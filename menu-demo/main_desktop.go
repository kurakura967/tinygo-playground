//go:build !tinygo

// menu-demo デスクトップ版。
// wio-emu サーバーと通信してエミュレータ画面にメニューを表示する。
//
// 使い方:
//  1. wio-emu サーバーを起動する:
//     cd ../../wiodisplay && go run ./cmd/wio-emu
//  2. このアプリを起動する:
//     go run .
package main

import (
	"image/color"
	"time"

	"github.com/kurakura967/wiodisplay/initdisplay"
	"github.com/kurakura967/wiodisplay/machine"
)

var (
	deskBg  = color.RGBA{R: 0x00, G: 0x00, B: 0x00, A: 0xFF}
	deskSep = color.RGBA{R: 0x20, G: 0x20, B: 0x40, A: 0xFF}
)

func main() {
	display := initdisplay.InitDisplay()
	game := NewGame()

	prevUp   := false
	prevDown := false
	prevBtnA := false
	prevBtnB := false

	for {
		rawUp   := !machine.WIO_5S_UP.Get()
		rawDown := !machine.WIO_5S_DOWN.Get()
		rawBtnA := !machine.WIO_KEY_A.Get()
		rawBtnB := !machine.WIO_KEY_B.Get()

		up   := rawUp   && !prevUp
		down := rawDown && !prevDown
		btnA := rawBtnA && !prevBtnA
		btnB := rawBtnB && !prevBtnB

		prevUp   = rawUp
		prevDown = rawDown
		prevBtnA = rawBtnA
		prevBtnB = rawBtnB

		game.update(up, down, btnA, btnB)

		if game.changed {
			display.FillScreen(deskBg)
			display.FillRectangle(0, 40, screenWidth, 1, deskSep)

			switch game.state {
			case StateMenu:
				drawMenu(display, game)
			case StateDetail:
				drawDetail(display, game)
			}
			game.changed = false
		}

		time.Sleep(16 * time.Millisecond)
	}
}
