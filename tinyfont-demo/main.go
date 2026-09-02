//go:build tinygo

// tinyfont-demo は tinyfont によるテキスト表示デモ（実機向け）。
// Wio Terminal 実機向けビルド:
//
//	tinygo build -target wioterminal -o tinyfont-demo.uf2 .
package main

import (
	"log"

	"github.com/kurakura967/wiodisplay/hardware"
	"github.com/sago35/koebiten"
	"tinygo.org/x/drivers/pixel"
)

var bgColor = pixel.NewColor[pixel.RGB565BE](0x00, 0x00, 0x00)
var sepColor = pixel.NewColor[pixel.RGB565BE](0x20, 0x20, 0x40)

// Update は koebiten.Game インターフェースの実装。
func (g *Game) Update() error {
	btnA := koebiten.IsKeyJustPressed(koebiten.Key0)
	g.update(btnA)
	return nil
}

// Draw は koebiten.Game インターフェースの実装。
// changed フラグが立っているときのみ全体再描画する。
func (g *Game) Draw(screen *koebiten.Image) {
	if !g.changed {
		return
	}
	// 背景クリア
	koebiten.DrawFilledRect(nil, 0, 0, screenWidth, screenHeight, bgColor)
	// タイトルエリアとの区切り線
	koebiten.DrawFilledRect(nil, 0, 48, screenWidth, 1, sepColor)
	// tinyfont でテキスト描画（hardware.Device.GetDisplay() は TextDisplay を満たす）
	drawText(hardware.Device.GetDisplay(), g)
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
	koebiten.SetWindowTitle("tinyfont Demo")
	if err := koebiten.RunGame(NewGame()); err != nil {
		log.Fatal(err)
	}
}
