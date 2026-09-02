//go:build tinygo

// menu-demo は tinyfont を使ったメニュー UI デモ（実機向け）。
// Wio Terminal 実機向けビルド:
//
//	tinygo build -target wioterminal -o menu-demo.uf2 .
package main

import (
	"log"

	"github.com/kurakura967/wiodisplay/hardware"
	"github.com/sago35/koebiten"
	"tinygo.org/x/drivers/pixel"
)

var (
	bgColor  = pixel.NewColor[pixel.RGB565BE](0x00, 0x00, 0x00)
	sepColor = pixel.NewColor[pixel.RGB565BE](0x20, 0x20, 0x40)
)

// Update は koebiten.Game インターフェースの実装。
func (g *Game) Update() error {
	up   := koebiten.IsKeyJustPressed(koebiten.KeyUp)
	down := koebiten.IsKeyJustPressed(koebiten.KeyDown)
	btnA := koebiten.IsKeyJustPressed(koebiten.Key0)
	btnB := koebiten.IsKeyJustPressed(koebiten.Key1)
	g.update(up, down, btnA, btnB)
	return nil
}

// Draw は koebiten.Game インターフェースの実装。
// changed フラグが立っているときのみ再描画する。
func (g *Game) Draw(screen *koebiten.Image) {
	if !g.changed {
		return
	}
	koebiten.DrawFilledRect(nil, 0, 0, screenWidth, screenHeight, bgColor)
	koebiten.DrawFilledRect(nil, 0, 40, screenWidth, 1, sepColor)

	d := hardware.Device.GetDisplay()
	switch g.state {
	case StateMenu:
		drawMenu(d, g)
	case StateDetail:
		drawDetail(d, g)
	}
	g.changed = false
}

// Layout は koebiten.Game インターフェースの実装。
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	if err := koebiten.SetHardware(hardware.Device); err != nil {
		log.Fatal(err)
	}
	koebiten.SetWindowSize(screenWidth, screenHeight)
	koebiten.SetWindowTitle("Menu Demo")
	if err := koebiten.RunGame(NewGame()); err != nil {
		log.Fatal(err)
	}
}
