package blelib

import (
	"encoding/binary"
	"errors"
	"github.com/cksmith/gatt"
	"github.com/cksmith/gatt/examples/option"
	log "github.com/sirupsen/logrus"
	"sync"
	"time"
)

const (
	DefaultInitTimeout = 1 * time.Second
	DefaultScanTimeout = 10 * time.Second
)

type MyDevice struct {
	gatt.Device
	PoweredDownCh <-chan struct{}

	poweredOnCh      chan error
	poweredDownCh    chan struct{}
	deviceFound      bool
	deviceFoundMutex sync.Mutex
	foundCh          chan gatt.Peripheral
	connectedCh      chan *MyPeripheral
	disconnectedCh   chan gatt.Peripheral
	p                *MyPeripheral
}

func newDevice(d gatt.Device) *MyDevice {
	poweredDownCh := make(chan struct{})
	return &MyDevice{
		Device:         d,
		PoweredDownCh:  poweredDownCh,
		poweredOnCh:    make(chan error),
		poweredDownCh:  poweredDownCh,
		foundCh:        make(chan gatt.Peripheral),
		connectedCh:    make(chan *MyPeripheral),
		disconnectedCh: make(chan gatt.Peripheral),
	}
}

type MyPeripheral struct {
	gatt.Peripheral
	disconnectedCh       chan struct{}
	ledCharacteristic    *gatt.Characteristic
	sensorCharacteristic *gatt.Characteristic
}

func newPeripheral(p gatt.Peripheral) *MyPeripheral {
	return &MyPeripheral{
		Peripheral:     p,
		disconnectedCh: make(chan struct{}),
	}
}

func onStateChanged(d *MyDevice) func(gattD gatt.Device, s gatt.State) {
	return func(gattD gatt.Device, s gatt.State) {
		log.WithField("state", s).Info("New state")
		switch s {
		case gatt.StatePoweredOn:
			d.poweredOnCh <- nil
		case gatt.StatePoweredOff:
			close(d.poweredDownCh)
		}
	}
}

func onPeripheralDiscovered(d *MyDevice) func(p gatt.Peripheral, a *gatt.Advertisement, rssi int) {
	return func(p gatt.Peripheral, a *gatt.Advertisement, rssi int) {
		d.deviceFoundMutex.Lock()
		defer d.deviceFoundMutex.Unlock()
		if !d.deviceFound {
			d.StopScanning()
			d.foundCh <- p
			d.deviceFound = true
		}
	}
}

func onPeripheralConnected(d *MyDevice) func(gattP gatt.Peripheral, err error) {
	return func(gattP gatt.Peripheral, err error) {
		log := log.WithFields(log.Fields{
			"id":   gattP.ID(),
			"name": gattP.Name(),
		})
		if d.p == nil {
			ss, err := gattP.DiscoverServices([]gatt.UUID{ServiceUUID})
			if err != nil {
				log.WithError(err).Error("Failed to discover services")
				return
			}
			for _, s := range ss {
				if ServiceUUID.Equal(s.UUID()) {
					cs, err := gattP.DiscoverCharacteristics(CharacteristicUUIDs, s)
					if err != nil {
						log.WithError(err).Error("Failed to discover characteristics")
						return
					}
					p := newPeripheral(gattP)
					cmap := make(map[string]*gatt.Characteristic)
					for _, c := range cs {
						cmap[c.UUID().String()] = c
					}
					p.ledCharacteristic = cmap[LEDCharacteristicUUID.String()]
					p.sensorCharacteristic = cmap[TemperatureHumidityCharacteristicUUID.String()]
					log.Info("Connected")
					d.p = p
					d.connectedCh <- p
					return
				}
			}
			log.Error("Failed to find service")
		} else {
			log.Error("Already connected")
		}
	}
}

func onPeripheralDisconnected(d *MyDevice) func(p gatt.Peripheral, err error) {
	return func(gattP gatt.Peripheral, err error) {
		log := log.WithFields(log.Fields{
			"id":   gattP.ID(),
			"name": gattP.Name(),
		})
		log.Info("Disconnected")
		if d.p != nil {
			close(d.p.disconnectedCh)
			d.p = nil
		} else {
			log.Info("Disconnecting peripheral that had not yet connected")
			d.disconnectedCh <- gattP
		}
	}
}

func NewDeviceClient() (*MyDevice, error) {
	gattD, err := gatt.NewDevice(option.DefaultClientOptions...)
	var d *MyDevice
	if err == nil {
		d, err = initDeviceClient(gattD)
	}
	return d, err
}

func NewSimDeviceClient(led Led, sensor TemperatureHumiditySensor) (*MyDevice, error) {
	gattD := gatt.NewSimDeviceClient(NewService(led, sensor), DeviceName)
	return initDeviceClient(gattD)
}

func initDeviceClient(gattD gatt.Device) (*MyDevice, error) {
	d := newDevice(gattD)
	err := gattD.Init(onStateChanged(d))
	if err == nil {
		select {
		case err = <-d.poweredOnCh:
		case <-time.After(DefaultInitTimeout):
			err = errors.New("Power on timeout")
		}
		if err == nil {
			gattD.Handle(
				gatt.PeripheralDiscovered(onPeripheralDiscovered(d)),
				gatt.PeripheralConnected(onPeripheralConnected(d)),
				gatt.PeripheralDisconnected(onPeripheralDisconnected(d)),
			)
		}
	}

	return d, err
}

func (d *MyDevice) Connect() (*MyPeripheral, error) {
	d.deviceFoundMutex.Lock()
	d.deviceFound = false
	d.deviceFoundMutex.Unlock()
	d.Scan([]gatt.UUID{ServiceUUID}, false)
	defer d.StopScanning()
	select {
	case p := <-d.foundCh:
		log.Info("Connecting...")
		p.Device().Connect(p)
		select {
		case connectedP := <-d.connectedCh:
			return connectedP, nil
		case <-d.disconnectedCh:
			return nil, errors.New("Disconnect callback received during connect")
		}
	case <-time.After(DefaultScanTimeout):
		log.Error("Scan timeout")
		return nil, errors.New("Scan timeout")
	}
}

func (d *MyDevice) CancelConnection() {
	if d.p != nil {
		d.Device.CancelConnection(d.p.Peripheral)
		<-d.p.disconnectedCh
	}
	d.p = nil
}

func (d *MyDevice) Stop() error {
	return d.Device.Stop()
}

// START1 OMIT

func (p MyPeripheral) SetLEDs(word byte) error {
	return p.WriteCharacteristic(p.ledCharacteristic, []byte{word}, false)
}

func (p MyPeripheral) GetTemperatureHumidity() (temperature float32, humidity float32, err error) {
	var data []byte
	data, err = p.ReadCharacteristic(p.sensorCharacteristic)
	if err == nil && len(data) == 8 {
		temperatureUint := binary.LittleEndian.Uint32(data[0:4])
		temperature = float32(int32(temperatureUint)) / 100.0
		humidityUint := binary.LittleEndian.Uint32(data[4:8])
		humidity = float32(humidityUint) / 100.0
	} else if len(data) < 8 {
		err = errors.New("Temperature/humidity characteristic failed to read sensor")
	}
	return
}

// END1 OMIT
