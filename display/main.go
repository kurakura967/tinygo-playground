package main

import (
	"image/color"
	"time"

	"machine"

	"tinygo.org/x/drivers/examples/ili9341/initdisplay"
	"tinygo.org/x/tinyfont"
	"tinygo.org/x/tinyfont/freemono"
)

const (
	btnA = machine.WIO_KEY_A
	btnB = machine.WIO_KEY_B
	btnC = machine.WIO_KEY_C
)

var texts = map[machine.Pin]string{
	btnA: "Button A!!",
	btnB: "Button B!",
	btnC: "Button C!",
}

func main() {
	btnA.Configure(machine.PinConfig{Mode: machine.PinInputPullup})
	btnB.Configure(machine.PinConfig{Mode: machine.PinInputPullup})
	btnC.Configure(machine.PinConfig{Mode: machine.PinInputPullup})

	display := initdisplay.InitDisplay()
	font := &freemono.Regular12pt7b
	y := int16(120) + 8

	drawText := func(text string) {
		display.FillScreen(color.RGBA{0, 0, 0, 255})
		_, w := tinyfont.LineWidth(font, text)
		x := int16(160) - int16(w)/2
		tinyfont.WriteLine(display, font, x, y, text, color.RGBA{255, 255, 255, 255})
	}

	drawText("Hello, World!!")

	buttons := []machine.Pin{btnA, btnB, btnC}
	prev := map[machine.Pin]bool{}

	for {
		for _, btn := range buttons {
			curr := btn.Get()
			if !curr && prev[btn] {
				drawText(texts[btn])
			}
			prev[btn] = curr
		}
		time.Sleep(50 * time.Millisecond)
	}
}
