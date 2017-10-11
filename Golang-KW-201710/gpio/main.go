package main

import (
	"github.com/stianeikeland/go-rpio"
	"time"
)

// START1 OMIT
var PinNumbers []int = []int{17, 18, 27, 22}

// END1 OMIT

// START2 OMIT
func main() {
	err := rpio.Open() // HL
	if err != nil {
		panic("Cannot open GPIO")
	}
	defer rpio.Close() // HL

	for _, pinNumber := range PinNumbers {
		pin := rpio.Pin(pinNumber)
		pin.Output() // make the pin an output (as opposed to an input) // HL
		pin.Low()    // set the output low // HL
	}
	// END2 OMIT

	// START3 OMIT
	value := 0
	for {
		word := value
		for _, pinNumber := range PinNumbers {
			pin := rpio.Pin(pinNumber)
			if (word & 1) == 1 {
				pin.High()
			} else {
				pin.Low()
			}
			word >>= 1
		}

		value = (value + 1) % 16
		time.Sleep(250 * time.Millisecond)
	}
	// END3 OMIT
}
