//go:build tinygo

// render-demo は koebiten の RunGame パターンを使ったレンダリングモード比較デモ。
// ボタン A でナイーブモード（毎フレーム全画面クリア）と最適化モード（dirty rect のみ更新）を切り替え。
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
	colorBg        = pixel.NewColor[pixel.RGB565BE](0x00, 0x00, 0x20) // 暗い青
	colorRect      = pixel.NewColor[pixel.RGB565BE](0xFF, 0x80, 0x00) // オレンジ
	colorNaive     = pixel.NewColor[pixel.RGB565BE](0xFF, 0x00, 0x00) // 赤 - ナイーブモードインジケーター
	colorOptimized = pixel.NewColor[pixel.RGB565BE](0x00, 0xFF, 0x00) // 緑 - 最適化モードインジケーター
)

// Update は koebiten.Game インターフェースの実装。
func (g *Game) Update() error {
	btnA := koebiten.IsKeyJustPressed(koebiten.Key0) // WIO_KEY_A - モード切り替え
	left := koebiten.IsKeyPressed(koebiten.KeyLeft)
	right := koebiten.IsKeyPressed(koebiten.KeyRight)
	up := koebiten.IsKeyPressed(koebiten.KeyUp)
	down := koebiten.IsKeyPressed(koebiten.KeyDown)

	g.update(btnA, left, right, up, down)
	return nil
}

// Draw は koebiten.Game インターフェースの実装。
func (g *Game) Draw(screen *koebiten.Image) {
	switch g.mode {
	case ModeNaive:
		// ナイーブモード: 毎フレーム全画面クリアしてから矩形を描画
		koebiten.DrawFilledRect(nil, 0, 0, screenWidth, screenHeight, colorBg)
		koebiten.DrawFilledRect(nil, g.x, g.y, rectSize, rectSize, colorRect)
		// モードインジケーター（左上 10x10）: 赤
		koebiten.DrawFilledRect(nil, 0, 0, 10, 10, colorNaive)

	case ModeOptimized:
		// 最適化モード: dirty rect のみ更新
		if !g.initialized {
			// 初回またはモード切り替え直後: 全画面クリア
			koebiten.DrawFilledRect(nil, 0, 0, screenWidth, screenHeight, colorBg)
			koebiten.DrawFilledRect(nil, g.x, g.y, rectSize, rectSize, colorRect)
			g.initialized = true
			// モードインジケーター（左上 10x10）: 緑
			koebiten.DrawFilledRect(nil, 0, 0, 10, 10, colorOptimized)
			return
		}
		// 新しい位置に矩形を描画（先に描く）
		koebiten.DrawFilledRect(nil, g.x, g.y, rectSize, rectSize, colorRect)
		// 旧位置の trailing edge を消去（後から消す）
		for _, r := range g.dirtyRects() {
			koebiten.DrawFilledRect(nil, r[0], r[1], r[2], r[3], colorBg)
		}
		// モードインジケーター（左上 10x10）: 緑
		koebiten.DrawFilledRect(nil, 0, 0, 10, 10, colorOptimized)
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
	koebiten.SetWindowTitle("Render Mode Demo")

	game := NewGame()
	if err := koebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
