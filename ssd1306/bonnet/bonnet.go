package bonnet

import (
	"sync"
	"time"

	"github.com/mikegleasonjr/adafruit/ssd1306"
	"periph.io/x/periph/conn/gpio"
	"periph.io/x/periph/conn/gpio/gpioreg"
)

// Those are the available buttons on the bonnet
const (
	A     = "A"
	B     = "B"
	C     = "C"
	Up    = "UP"
	Down  = "DOWN"
	Left  = "LEFT"
	Right = "RIGHT"
)

var buttons = []struct {
	name string
	pin  string
}{
	{A, "GPIO5"},
	{B, "GPIO6"},
	{C, "GPIO4"},
	{Up, "GPIO17"},
	{Down, "GPIO22"},
	{Left, "GPIO27"},
	{Right, "GPIO23"},
}

// Event represent a button change event
type Event struct {
	When    time.Time
	Button  string
	Pressed bool
}

// Bonnet drives the Adafruit 128x64 OLED Bonnet for Raspberry Pi
// see https://www.adafruit.com/product/3531
type Bonnet struct {
	*ssd1306.I2C
	wg             *sync.WaitGroup
	stop           chan struct{}
	Events         chan Event
	internalEvents chan Event
}

// New creates a driver for the Adafruit 128x64 OLED Bonnet for Raspberry Pi.
// Name can be an I²C bus name, an alias or a number. Specify an empty name ""
// to get the first available bus.
func New(name string, rotated bool) (*Bonnet, error) {
	display, err := ssd1306.NewI2C(name, 128, 64, rotated)
	if err != nil {
		return nil, err
	}

	b := &Bonnet{
		I2C:            display,
		wg:             &sync.WaitGroup{},
		stop:           make(chan struct{}),
		Events:         make(chan Event),
		internalEvents: make(chan Event, 10),
	}

	b.listenToButtons()
	b.wg.Add(1)
	go b.drainEvents(time.Second)

	return b, nil
}

// Close closes the display
func (b *Bonnet) Close(turnoff bool) {
	b.I2C.Close(turnoff)
	close(b.stop)
	b.wg.Wait()
}

func (b *Bonnet) listenToButtons() {
	for _, btn := range buttons {
		b.wg.Add(1)

		go func(btn struct {
			name string
			pin  string
		}) {
			defer b.wg.Done()

			pin := gpioreg.ByName(btn.pin)
			pin.In(gpio.PullUp, gpio.BothEdges)
			for {
				changed := pin.WaitForEdge(1 * time.Second)
				if !changed {
					select {
					case <-b.stop:
						return
					default:
						continue
					}
				}
				select {
				case b.internalEvents <- Event{time.Now(), btn.name, pin.Read() == gpio.Low}:
				case <-b.stop:
					return
				default:
				}
			}
		}(btn)
	}
}

func (b *Bonnet) drainEvents(ttl time.Duration) {
	defer b.wg.Done()
	fifo := []Event{}

	for {
		select {
		case msg := <-b.internalEvents:
			fifo = append(fifo, msg)
			msg = fifo[0]
			since := time.Now().Sub(msg.When)
			if since < ttl {
				select {
				case b.Events <- msg:
				case <-time.After(ttl - since):
				case <-b.stop:
					return
				}
			}
			fifo = fifo[1:]
		case <-b.stop:
			return
		}
	}
}
