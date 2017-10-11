// +build !linux

package blelib

import (
	"errors"
)

type sensor struct {
}

func NewTemperatureHumidityService() TemperatureHumiditySensor {
	return &sensor{}
}

func (s sensor) Read() (temperature float32, humidity float32, retried int, err error) {
	return 0, 0, 0, errors.New("Not supported on non-Linux platforms")
}
