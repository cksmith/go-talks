package blelib

import (
	"github.com/stianeikeland/go-rpio"
)

var PinNumbers []int = []int{17, 18, 27, 22}

// START1 OMIT

type Led interface {
	Open() error
	Close() error
	Write(word byte)
}

type led struct {
	// Don't need to store any state
}

func NewLedService() Led {
	return &led{}
}

// Method implementations for Open(), Close(), Write()
// ...
// END1 OMIT

func (l led) Open() error {
	err := rpio.Open()
	if err == nil {
		for _, pinNumber := range PinNumbers {
			pin := rpio.Pin(pinNumber)
			pin.Output()
			pin.Low()
		}
	}
	return err
}

func (l led) Close() error {
	return rpio.Close()
}

func (l led) Write(word byte) {
	for _, pinNumber := range PinNumbers {
		pin := rpio.Pin(pinNumber)
		if (word & 1) == 1 {
			pin.High()
		} else {
			pin.Low()
		}
		word >>= 1
	}
}

// START2 OMIT

type MockLed struct {
	LastWrittenWord byte // Publicly accessible
}

func NewMockLedService() *MockLed {
	return &MockLed{}
}

func (l MockLed) Open() error {
	return nil
}

func (l MockLed) Close() error {
	return nil
}

func (l *MockLed) Write(word byte) {
	l.LastWrittenWord = word
}

// END2 OMIT
