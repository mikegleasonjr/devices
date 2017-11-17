package pioled

import (
	"github.com/mikegleasonjr/devices/ssd1306"
)

// PiOLED drives the Adafruit PiOLED - 128x32 Monochrome OLED Add-on for Raspberry Pi
// https://www.adafruit.com/product/3527
type PiOLED struct {
	*ssd1306.I2C
}

// New creates a driver for the Adafruit PiOLED - 128x32 Monochrome OLED Add-on for Raspberry Pi.
// Name can be an I²C bus name, an alias or a number. Specify an empty name ("")
// to get the first available bus.
func New(name string, rotated bool) (*PiOLED, error) {
	display, err := ssd1306.NewI2C(name, 128, 32, rotated)
	if err != nil {
		return nil, err
	}

	if err = display.Tx([]byte{
		0x00, // I²C transaction has stream of command bytes
		0xDA, // set comm pins
		0x02, // 2
	}); err != nil {
		return nil, err
	}

	return &PiOLED{display}, nil
}
