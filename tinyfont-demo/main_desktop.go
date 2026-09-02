//go:build !tinygo

// tinyfont-demo デスクトップ版。
// wio-emu サーバーと通信してエミュレータ画面にテキストを表示する。
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

	prevBtnA := false

	for {
		rawBtnA := !machine.WIO_KEY_A.Get()
		btnA := rawBtnA && !prevBtnA
		prevBtnA = rawBtnA

		game.update(btnA)

		if game.changed {
			// 背景クリア
			display.FillScreen(deskBg)
			// タイトルエリアとの区切り線
			display.FillRectangle(0, 48, screenWidth, 1, deskSep)
			// tinyfont でテキスト描画（*ili9341.Device は TextDisplay を満たす）
			drawText(display, game)
			game.changed = false
		}

		time.Sleep(16 * time.Millisecond)
	}
}
