package ssd1306

import (
	"image"

	"github.com/mikegleasonjr/devices"
	"periph.io/x/periph/conn/i2c"
	"periph.io/x/periph/conn/i2c/i2creg"
	"periph.io/x/periph/devices/ssd1306"
)

// I2C is an OLED display with hardware I2C.
type I2C struct {
	*ssd1306.Dev
	c i2c.BusCloser
}

// NewI2C creates an OLED display
func NewI2C(name string, width, height int, rotated bool) (*I2C, error) {
	if err := devices.Init(); err != nil {
		return nil, err
	}

	c, err := i2creg.Open(name)
	if err != nil {
		return nil, err
	}

	s, err := ssd1306.NewI2C(c, width, height, rotated)
	if err != nil {
		c.Close()
		return nil, err
	}

	s.Draw(s.Bounds(), image.Black, image.ZP)
	return &I2C{c: c, Dev: s}, nil
}

// Close closes the display
func (s *I2C) Close(turnoff bool) {
	if turnoff {
		s.Halt()
	}
	s.c.Close()
}

// Tx sends raw commands to the displays
func (s *I2C) Tx(w []byte) error {
	return s.c.Tx(0x3C, w, nil)
}
