package blelib

// START1 OMIT

import (
	"github.com/d2r2/go-dht"
)

const (
	SensorModel      = dht.DHT22
	SensorPin        = 4
	MaxRetries       = 10
	BoostPerformance = false // shouldn't be necessary with RPi 3, requires root
)

type sensor struct {
}

func NewTemperatureHumidityService() TemperatureHumiditySensor {
	return &sensor{}
}

func (s sensor) Read() (temperature float32, humidity float32, retried int, err error) {
	return dht.ReadDHTxxWithRetry(SensorModel, SensorPin, BoostPerformance, MaxRetries)
}

// END1 OMIT
