package blelib

import (
	"flag"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
	"testing"
)

// START1 OMIT

var (
	testWithDevice bool
)

func init() {
	log.SetLevel(log.DebugLevel)

	flag.BoolVar(&testWithDevice, "device", false, "Test using a real device")
	flag.Parse()
}

// END1 OMIT
// START2 OMIT

type ClientSuite struct {
	suite.Suite
	d *MyDevice
	p *MyPeripheral
}

func (s *ClientSuite) SetupTest() {
	var err error
	// Scan and connect to the first device that implements our service
	s.p, err = s.d.Connect()
	s.Require().NoError(err, "Failed to connect")
}

func (s *ClientSuite) TearDownTest() {
	if s.d != nil {
		s.d.CancelConnection()
	}
}

func (s *ClientSuite) runIfRealDevice(t *testing.T, f func(t *testing.T)) {
	if !testWithDevice {
		t.Skip("Skipping real device tests. -device parameter not specified.")
	} else {
		f(t)
	}
}

// END2 OMIT
