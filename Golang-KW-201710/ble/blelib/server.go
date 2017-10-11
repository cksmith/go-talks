package blelib

import (
	"context"
	"fmt"
	"github.com/cksmith/gatt"
	"github.com/cksmith/gatt/examples/option"
	"github.com/cksmith/gatt/examples/service"
)

// START1 OMIT
func NewServer(ctx context.Context, led Led, sensor TemperatureHumiditySensor) {
	go func() {
		d, err := gatt.NewDevice(option.DefaultServerOptions...)
		if err != nil {
			panic("Failed to open device")
		}

		// Register optional handlers. Not called on Mac OS X (handled by OS).
		d.Handle(
			gatt.CentralConnected(func(c gatt.Central) { fmt.Println("Connect:", c.ID()) }),
			gatt.CentralDisconnected(func(c gatt.Central) { fmt.Println("Disconnect:", c.ID()) }),
		)
		// END1 OMIT
		// START2 OMIT
		// A mandatory handler for monitoring device state.
		onStateChanged := func(d gatt.Device, s gatt.State) {
			switch s {
			case gatt.StatePoweredOn:
				// Setup GAP and GATT services for Linux implementation
				// Mac OS X provides them
				d.AddService(service.NewGapService(DeviceName))
				d.AddService(service.NewGattService())

				s := NewService(led, sensor)
				d.AddService(s)
				d.AdvertiseNameAndServices(DeviceName, []gatt.UUID{s.UUID()})
			default:
			}
		}

		d.Init(onStateChanged)

		<-ctx.Done()

		d.StopAdvertising()
		d.RemoveAllServices()
		d.Stop()
	}()
}

// END2 OMIT
