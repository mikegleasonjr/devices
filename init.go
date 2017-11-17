package devices

import (
	"sync"

	"periph.io/x/periph/host"
)

var once = sync.Once{}
var err error

// Init initializes your device.
func Init() error {
	once.Do(func() {
		_, err = host.Init()
	})
	return err
}
