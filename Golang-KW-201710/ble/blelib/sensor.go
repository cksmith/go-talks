package blelib

// START1 OMIT

type TemperatureHumiditySensor interface {
	Read() (temperature float32, humidity float32, retried int, err error)
}

type MockSensor struct {
	NextTemperature float32
	NextHumidity    float32
}

func NewMockTemperatureHumidityService() *MockSensor {
	return &MockSensor{}
}

func (s MockSensor) Read() (temperature float32, humidity float32, retried int, err error) {
	return s.NextTemperature, s.NextHumidity, 0, nil
}

// END1 OMIT
