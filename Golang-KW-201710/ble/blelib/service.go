package blelib

// START1 OMIT
import (
	"encoding/binary"
	"github.com/cksmith/gatt"
	log "github.com/sirupsen/logrus"
)

func NewService(led Led, sensor TemperatureHumiditySensor) *gatt.Service {
	// Globally unique IDs required for service and characteristics
	s := gatt.NewService(ServiceUUID)

	c := s.AddCharacteristic(LEDCharacteristicUUID)
	c.AddDescriptor(gatt.UUID16(0x2901)).SetStringValue("LED Values")
	c.HandleWriteFunc(func(r gatt.Request, data []byte) byte {
		log.WithField("data", data[0]).Info("Write LEDs")
		led.Write(data[0])
		return gatt.StatusSuccess
	})

	// END1 OMIT

	// START2 OMIT
	c = s.AddCharacteristic(TemperatureHumidityCharacteristicUUID)
	c.AddDescriptor(gatt.UUID16(0x2901)).SetStringValue("Temperature and Humidity")
	c.HandleReadFunc(
		func(rsp gatt.ResponseWriter, req *gatt.ReadRequest) {
			temperature, humidity, retries, err := sensor.Read()
			if err == nil {
				log.WithFields(log.Fields{
					"temperature": temperature,
					"humidity":    humidity,
					"retries":     retries,
				}).Info("Sensor values")
				temperatureInt := int32(temperature * 100.0)
				humidityInt := uint32(humidity * 100.0)
				binary.Write(rsp, binary.LittleEndian, temperatureInt)
				binary.Write(rsp, binary.LittleEndian, humidityInt)
			} else {
				log.WithError(err).Error("Error reading sensor")
				rsp.Write([]byte{0})
			}
		})
	return s
}

// END2 OMIT
