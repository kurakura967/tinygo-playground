package main

import (
	"image/color"
	"machine"
	"time"

	"tinygo.org/x/drivers/ili9341"
	"tinygo.org/x/tinyfont"
	"tinygo.org/x/tinyfont/freemono"
)

func main() {
	machine.SPI3.Configure(machine.SPIConfig{
		SCK:       machine.LCD_SCK_PIN,
		SDO:       machine.LCD_SDO_PIN,
		SDI:       machine.LCD_SDI_PIN,
		Frequency: 40_000_000,
	})

	backlight := machine.LCD_BACKLIGHT
	backlight.Configure(machine.PinConfig{Mode: machine.PinOutput})

	display := ili9341.NewSPI(
		machine.SPI3,
		machine.LCD_DC,
		machine.LCD_SS_PIN,
		machine.LCD_RESET,
	)

	display.Configure(ili9341.Config{})
	backlight.High()
	display.SetRotation(ili9341.Rotation270)

	display.FillScreen(color.RGBA{0, 0, 0, 255})

	const text = "Hello, World!"
	font := &freemono.Regular12pt7b
	_, w := tinyfont.LineWidth(font, text)
	x := int16(160) - int16(w)/2
	y := int16(120) + 8

	tinyfont.WriteLine(display, font, x, y, text, color.RGBA{255, 255, 255, 255})

	for {
		time.Sleep(1 * time.Second)
	}
}
