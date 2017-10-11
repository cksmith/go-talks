package blelib

import (
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

// START1 OMIT

type MockServiceSuite struct {
	ClientSuite
	led    *MockLed
	sensor *MockSensor
}

func TestMockServiceSuite(t *testing.T) {
	suite.Run(t, new(MockServiceSuite))
}

func (s *MockServiceSuite) SetupTest() {
	s.led = NewMockLedService()
	s.sensor = NewMockTemperatureHumidityService()
	var err error
	s.d, err = NewSimDeviceClient(s.led, s.sensor)
	s.Require().NoError(err, "Could not initialize mock device client")
	s.ClientSuite.SetupTest()
}

// END1 OMIT

// START2 OMIT

func (s *MockServiceSuite) TestWriteLedValues() {
	for value := byte(0); value < 16; value += 1 {
		err := s.p.SetLEDs(value)
		s.NoError(err, "Failed to set LED values")
		s.Equal(value, s.led.LastWrittenWord, "Last written word did not match")
	}
}

func (s *MockServiceSuite) TestReadSensorValues() {
	s.sensor.NextTemperature = 25.1
	s.sensor.NextHumidity = 45.1
	temperature, humidity, err := s.p.GetTemperatureHumidity()
	s.NoError(err, "Failed to read temperature and humidity values")
	s.Equal(s.sensor.NextTemperature, temperature, "Temperature did not match")
	s.Equal(s.sensor.NextHumidity, humidity, "Humidity did not match")
	s.sensor.NextTemperature = -22.1
	temperature, _, err = s.p.GetTemperatureHumidity()
	s.NoError(err, "Failed to read temperature and humidity values")
	s.Equal(s.sensor.NextTemperature, temperature, "Negative temperature did not match")
}

// END2 OMIT

// START3 OMIT

type RealServiceSuite struct {
	ClientSuite
}

func TestRealServiceSuiteSuite(t *testing.T) {
	s := new(RealServiceSuite)
	s.runIfRealDevice(t, func(t *testing.T) {
		suite.Run(t, s)
	})
}

func (s *RealServiceSuite) SetupTest() {
	var err error
	s.d, err = NewDeviceClient()
	s.Require().NoError(err, "Could not initialize device client")
	s.ClientSuite.SetupTest()
}

// END3 OMIT

// START4 OMIT

func (s *RealServiceSuite) TestWriteLedValues() {
	for value := byte(0); value < 16; value += 1 {
		err := s.p.SetLEDs(value)
		s.NoError(err, "Failed to set LED values")
		time.Sleep(250 * time.Millisecond)
	}
	err := s.p.SetLEDs(0)
	s.NoError(err, "Failed to set LED values")
}

func (s *RealServiceSuite) TestReadSensorValues() {
	temperature, humidity, err := s.p.GetTemperatureHumidity()
	s.NoError(err, "Failed to read temperature and humidity values")
	log.WithFields(log.Fields{
		"temperature": temperature,
		"humidity":    humidity,
	}).Info("Read sensor values")
}

// END4 OMIT
