package ssd1306

import (
	"image"
	"sync"

	"periph.io/x/periph/conn/i2c"
	"periph.io/x/periph/conn/i2c/i2creg"
	"periph.io/x/periph/devices/ssd1306"
	"periph.io/x/periph/host"
)

var hostInit = func() func() error {
	var once = sync.Once{}
	var err error

	return func() error {
		once.Do(func() {
			_, err = host.Init()
		})
		return err
	}
}()

// I2C is an AdaFruit OLED display with hardware I2C that don't have a reset pin.
type I2C struct {
	c i2c.BusCloser
	s *ssd1306.Dev
}

// NewI2C creates an OLED display
func NewI2C(name string, width, height int, rotated bool) (*I2C, error) {
	if err := hostInit(); err != nil {
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
	return &I2C{c, s}, nil
}

// Close closes the display
func (s *I2C) Close(turnoff bool) {
	if turnoff {
		s.s.Halt()
	}
	s.c.Close()
}

// Draw writes an image to the display
func (s *I2C) Draw(r image.Rectangle, src image.Image, sp image.Point) {
	s.s.Draw(r, src, sp)
}

// Bounds returns the display size
func (s *I2C) Bounds() image.Rectangle {
	return s.s.Bounds()
}
