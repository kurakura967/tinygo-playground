//go:build tinygo

// emu-demo は koebiten の RunGame パターンを使った矩形アニメーションのサンプル。
// Wio Terminal 実機向けビルド:
//
//	tinygo build -target wioterminal -o demo.uf2 .
package main

import (
	"log"

	"github.com/kurakura967/wiodisplay/hardware"
	"github.com/sago35/koebiten"
	"tinygo.org/x/drivers/pixel"
)

var (
	colorBg   = pixel.NewColor[pixel.RGB565BE](0x00, 0x00, 0x20) // 暗い青
	colorRect = pixel.NewColor[pixel.RGB565BE](0xFF, 0x80, 0x00) // オレンジ
	colorA    = pixel.NewColor[pixel.RGB565BE](0xFF, 0x00, 0x00) // 赤
	colorB    = pixel.NewColor[pixel.RGB565BE](0x00, 0xFF, 0x00) // 緑
	colorC    = pixel.NewColor[pixel.RGB565BE](0x00, 0x00, 0xFF) // 青
)

// Update は koebiten.Game インターフェースの実装。
func (g *Game) Update() error {
	btnA := koebiten.IsKeyPressed(koebiten.Key0)  // WIO_KEY_A
	btnB := koebiten.IsKeyPressed(koebiten.Key1)  // WIO_KEY_B
	btnC := koebiten.IsKeyPressed(koebiten.Key2)  // WIO_KEY_C
	left := koebiten.IsKeyPressed(koebiten.KeyLeft)
	right := koebiten.IsKeyPressed(koebiten.KeyRight)
	up := koebiten.IsKeyPressed(koebiten.KeyUp)
	down := koebiten.IsKeyPressed(koebiten.KeyDown)

	// ボタン C で一時停止
	if btnC {
		return nil
	}

	g.update(btnA, btnB, btnC, left, right, up, down)
	return nil
}

// Draw は koebiten.Game インターフェースの実装。
func (g *Game) Draw(screen *koebiten.Image) {
	// ボタン状態に応じて矩形の色を変える
	rectColor := colorRect
	if koebiten.IsKeyPressed(koebiten.Key0) {
		rectColor = colorA
	} else if koebiten.IsKeyPressed(koebiten.Key1) {
		rectColor = colorB
	} else if koebiten.IsKeyPressed(koebiten.Key2) {
		rectColor = colorC
	}

	if !g.initialized {
		// 初回: 背景全体と矩形全体を描画
		koebiten.DrawFilledRect(nil, 0, 0, screenWidth, screenHeight, colorBg)
		koebiten.DrawFilledRect(nil, g.x, g.y, rectSize, rectSize, rectColor)
		g.initialized = true
		return
	}

	// 新しい位置に矩形を描画（先に描く）
	koebiten.DrawFilledRect(nil, g.x, g.y, rectSize, rectSize, rectColor)
	// 旧位置の trailing edge を消去（後から消す）
	for _, r := range g.dirtyRects() {
		koebiten.DrawFilledRect(nil, r[0], r[1], r[2], r[3], colorBg)
	}
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
	koebiten.SetWindowTitle("Wio Terminal Demo")

	game := NewGame()
	if err := koebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
