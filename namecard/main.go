// namecard は Go Conference 2026 用のスピーカー名刺を Wio Terminal に表示する。
//
// Desktop（エミュレータ）:
//
//	go run .
//
// Wio Terminal:
//
//	tinygo build -target wioterminal -o namecard.uf2 .
package main

import (
	"image/color"
	"time"

	"github.com/kurakura967/wiodisplay/driver/ili9341"
	"github.com/kurakura967/wiodisplay/initdisplay"
	"github.com/kurakura967/wiodisplay/machine"
	"tinygo.org/x/tinyfont"
	"tinygo.org/x/tinyfont/freemono"
	"tinygo.org/x/tinyfont/shnm"
)

const (
	screenW = int16(320)
	screenH = int16(240)

	logoY = int16(0)

	// 下部エリア: ロゴの下 100px
	bottomY = int16(logoH)
	bottomH = screenH - bottomY

	// アバター: 下部エリア左端に縦中央揃え
	avatarX = int16(4)
	avatarY = bottomY + (bottomH-int16(avatarH))/2

	// テキスト: アバターの右
	textX = avatarX + int16(avatarW) + 8

	// QR コード: 画面中央
	qrX = (screenW - int16(qrW)) / 2
	qrY = (screenH - int16(qrH)) / 2
)

var (
	colorBg   = color.RGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF} // 白
	colorText = color.RGBA{R: 0x1B, G: 0x34, B: 0x6B, A: 0xFF} // ネイビー
)

// drawScreen1 は名刺画面を描画する（ロゴ + アバター + 名前 + ハンドル）。
func drawScreen1(d *ili9341.Device) {
	d.FillScreen(colorBg)
	d.DrawRGBBitmap8(0, logoY, logo[:], int16(logoW), int16(logoH))
	d.DrawRGBBitmap8(avatarX, avatarY, avatar[:], int16(avatarW), int16(avatarH))

	nameFont := &freemono.Bold12pt7b
	nameY := bottomY + 38
	tinyfont.WriteLine(d, nameFont, textX, nameY, "Hiroki.Kurasawa", colorText)

	handleY := nameY + 34
	tinyfont.WriteLine(d, nameFont, textX, handleY, "@kurasawah", colorText)
}

// drawScreen2 はセッション情報画面を描画する（セッション名のみ）。
func drawScreen2(d *ili9341.Device) {
	d.FillScreen(colorBg)

	font := &shnm.Shnmk12
	y := int16(105)
	tinyfont.WriteLine(d, font, 4, y, "TinyGo 開発サイクルを高速化する：", colorText)
	y += 22
	tinyfont.WriteLine(d, font, 4, y, "Go で作るエミュレータ入門", colorText)
}

// drawScreen3 は QR コード画面を描画する。
func drawScreen3(d *ili9341.Device) {
	d.FillScreen(colorBg)
	d.DrawRGBBitmap8(qrX, qrY, qr[:], int16(qrW), int16(qrH))
}

func main() {
	display := initdisplay.InitDisplay()

	btnA := machine.WIO_KEY_A
	btnB := machine.WIO_KEY_B
	btnC := machine.WIO_KEY_C
	btnA.Configure(machine.PinConfig{Mode: machine.PinInputPullup})
	btnB.Configure(machine.PinConfig{Mode: machine.PinInputPullup})
	btnC.Configure(machine.PinConfig{Mode: machine.PinInputPullup})

	drawScreen1(display)

	for {
		switch {
		case !btnA.Get():
			drawScreen1(display)
			time.Sleep(200 * time.Millisecond)
		case !btnB.Get():
			drawScreen2(display)
			time.Sleep(200 * time.Millisecond)
		case !btnC.Get():
			drawScreen3(display)
			time.Sleep(200 * time.Millisecond)
		}
		time.Sleep(10 * time.Millisecond)
	}
}
