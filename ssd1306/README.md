# ssd1306

## Examples

### Adafruit 128x64 OLED Bonnet for Raspberry Pi

https://www.adafruit.com/product/3531

```golang
package main

import (
	"fmt"
	"image"
	"log"
	"time"

	"github.com/fogleman/gg"
	"github.com/mikegleasonjr/devices/ssd1306/adafruit/bonnet"
)

func main() {
	b, err := bonnet.New("", false)
	if err != nil {
		log.Fatal(err)
	}
	defer b.Close(true)

	t := time.NewTicker(time.Second)
	defer t.Stop()

	ctx := gg.NewContext(b.Bounds().Dx(), b.Bounds().Dy())

	for {
		select {
		case <-t.C:
			ctx.SetRGB(0, 0, 0)
			ctx.Clear()
			ctx.SetRGB(1, 1, 1)
			ctx.DrawStringAnchored(time.Now().Format("15:04:05"), 0, 0, 0, 1)
			b.Draw(b.Bounds(), ctx.Image(), image.ZP)
		case e := <-b.Events:
			fmt.Println("button", e.Button, "pressed:", e.Pressed)
			if e.Button == bonnet.C {
				return
			}
		}
	}
}
```

### Adafruit PiOLED - 128x32 Monochrome OLED Add-on for Raspberry Pi

https://www.adafruit.com/product/3527


```golang
package main

import (
	"image"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/fogleman/gg"
	"github.com/mikegleasonjr/devices/ssd1306/adafruit/pioled"
)

func main() {
	p, err := pioled.New("", false)
	if err != nil {
		log.Fatal(err)
	}
	defer p.Close(true)

	t := time.NewTicker(time.Second)
	defer t.Stop()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	ctx := gg.NewContext(p.Bounds().Dx(), p.Bounds().Dy())

	for {
		select {
		case <-t.C:
			ctx.SetRGB(0, 0, 0)
			ctx.Clear()
			ctx.SetRGB(1, 1, 1)
			ctx.DrawStringAnchored(time.Now().Format("15:04:05"), 0, 0, 0, 1)
			p.Draw(p.Bounds(), ctx.Image(), image.ZP)
		case <-c:
			return
		}
	}
}
```