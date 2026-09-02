//go:build tinygo

// breakout-demo は koebiten の RunGame パターンを使ったブレイクアウトゲームのサンプル。
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
	colorBg     = pixel.NewColor[pixel.RGB565BE](0x00, 0x00, 0x00) // 黒
	colorPaddle = pixel.NewColor[pixel.RGB565BE](0xFF, 0xFF, 0xFF) // 白
	colorBall   = pixel.NewColor[pixel.RGB565BE](0xFF, 0xFF, 0x00) // 黄
	// ブロック行ごとの色（0行目が上）
	colorBlocks = [blockRows]pixel.BaseColor{
		pixel.NewColor[pixel.RGB565BE](0xFF, 0x00, 0x00), // 赤
		pixel.NewColor[pixel.RGB565BE](0xFF, 0x80, 0x00), // オレンジ
		pixel.NewColor[pixel.RGB565BE](0xFF, 0xFF, 0x00), // 黄
		pixel.NewColor[pixel.RGB565BE](0x00, 0xFF, 0x00), // 緑
		pixel.NewColor[pixel.RGB565BE](0x00, 0x80, 0xFF), // 青
	}
)

// Update は koebiten.Game インターフェースの実装。
func (g *Game) Update() error {
	left   := koebiten.IsKeyPressed(koebiten.KeyLeft)
	right  := koebiten.IsKeyPressed(koebiten.KeyRight)
	launch := koebiten.IsKeyJustPressed(koebiten.Key0) // ボタン A
	g.update(left, right, launch)
	return nil
}

var (
	colorGameOver = pixel.NewColor[pixel.RGB565BE](0xFF, 0x00, 0x00) // 赤
	colorClear    = pixel.NewColor[pixel.RGB565BE](0x00, 0xFF, 0x00) // 緑
)

// Draw は koebiten.Game インターフェースの実装。
// ClearBuffer は no-op のため、dirty rect で必要な箇所だけ更新する。
func (g *Game) Draw(screen *koebiten.Image) {
	if !g.initialized {
		// 初回またはリスタート直後: 全画面クリアして全要素を描画
		koebiten.DrawFilledRect(nil, 0, 0, screenWidth, screenHeight, colorBg)
		for r := 0; r < blockRows; r++ {
			for c := 0; c < blockCols; c++ {
				if g.blocks[r][c] {
					bx, by, bw, bh := blockRect(r, c)
					koebiten.DrawFilledRect(nil, bx, by, bw, bh, colorBlocks[r])
				}
			}
		}
		koebiten.DrawFilledRect(nil, g.paddleX, paddleY, paddleW, paddleH, colorPaddle)
		koebiten.DrawFilledRect(nil, g.ballX, g.ballY, ballSize, ballSize, colorBall)
		g.initialized = true
		return
	}

	// 破壊されたブロックを消去
	for r := 0; r < blockRows; r++ {
		for c := 0; c < blockCols; c++ {
			if g.destroyed[r][c] {
				bx, by, bw, bh := blockRect(r, c)
				koebiten.DrawFilledRect(nil, bx, by, bw, bh, colorBg)
				g.destroyed[r][c] = false
			}
		}
	}

	// パドルの旧位置を消去して新位置に描画
	if g.paddleX != g.paddleXPrev {
		koebiten.DrawFilledRect(nil, g.paddleXPrev, paddleY, paddleW, paddleH, colorBg)
		koebiten.DrawFilledRect(nil, g.paddleX, paddleY, paddleW, paddleH, colorPaddle)
	}

	// ボールの旧位置を消去して新位置に描画
	koebiten.DrawFilledRect(nil, g.ballXPrev, g.ballYPrev, ballSize, ballSize, colorBg)
	koebiten.DrawFilledRect(nil, g.ballX, g.ballY, ballSize, ballSize, colorBall)

	// GameOver / Clear の状態表示
	switch g.state {
	case StateGameOver:
		koebiten.DrawFilledRect(nil, 60, 108, 200, 4, colorGameOver)
		koebiten.DrawFilledRect(nil, 60, 116, 200, 4, colorGameOver)
	case StateClear:
		koebiten.DrawFilledRect(nil, 60, 108, 200, 4, colorClear)
		koebiten.DrawFilledRect(nil, 60, 116, 200, 4, colorClear)
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
	koebiten.SetWindowTitle("Breakout Demo")
	if err := koebiten.RunGame(NewGame()); err != nil {
		log.Fatal(err)
	}
}
