//go:build tinygo

// tetris-demo は koebiten の RunGame パターンを使ったテトリスゲームのサンプル。
// Wio Terminal 実機向けビルド:
//
//	tinygo build -target wioterminal -o tetris.uf2 .
package main

import (
	"image/color"
	"log"
	"strconv"

	"github.com/kurakura967/wiodisplay/hardware"
	"github.com/sago35/koebiten"
	"tinygo.org/x/drivers/pixel"
	"tinygo.org/x/tinyfont"
	"tinygo.org/x/tinyfont/freemono"
)

var (
	colorBg     = pixel.NewColor[pixel.RGB565BE](0x00, 0x00, 0x00) // 黒
	colorBorder = pixel.NewColor[pixel.RGB565BE](0x40, 0x40, 0x40) // グレー（フィールド枠）
	// ピース種別ごとの色（0=空, 1-7=I/O/T/S/Z/J/L）
	colorPieces = [8]pixel.BaseColor{
		pixel.NewColor[pixel.RGB565BE](0x00, 0x00, 0x00), // 0: 空（黒）
		pixel.NewColor[pixel.RGB565BE](0x00, 0xFF, 0xFF), // 1: I シアン
		pixel.NewColor[pixel.RGB565BE](0xFF, 0xFF, 0x00), // 2: O 黄
		pixel.NewColor[pixel.RGB565BE](0x80, 0x00, 0xFF), // 3: T 紫
		pixel.NewColor[pixel.RGB565BE](0x00, 0xFF, 0x00), // 4: S 緑
		pixel.NewColor[pixel.RGB565BE](0xFF, 0x00, 0x00), // 5: Z 赤
		pixel.NewColor[pixel.RGB565BE](0x00, 0x00, 0xFF), // 6: J 青
		pixel.NewColor[pixel.RGB565BE](0xFF, 0x80, 0x00), // 7: L オレンジ
	}
)

// drawCell はフィールド座標 (fx, fy) にセルを描画する。
func drawCell(fx, fy int, c pixel.BaseColor) {
	koebiten.DrawFilledRect(nil, fieldX+fx*cellSize, fieldY+fy*cellSize, cellSize-1, cellSize-1, c)
}

// eraseCell はフィールド座標 (fx, fy) のセルを背景色で消去する。
func eraseCell(fx, fy int) {
	koebiten.DrawFilledRect(nil, fieldX+fx*cellSize, fieldY+fy*cellSize, cellSize, cellSize, colorBg)
}

// drawPiece は shape をフィールド座標 (px, py) に色 c で描画する。
func drawPiece(shape [4][4]bool, px, py int, c pixel.BaseColor) {
	for row := 0; row < 4; row++ {
		for col := 0; col < 4; col++ {
			if shape[row][col] && py+row >= 0 {
				drawCell(px+col, py+row, c)
			}
		}
	}
}

// erasePiece は shape をフィールド座標 (px, py) から消去する。
// フィールドに固定ブロックがあればその色、なければ背景色で塗る。
func erasePiece(shape [4][4]bool, px, py int, g *Game) {
	for row := 0; row < 4; row++ {
		for col := 0; col < 4; col++ {
			if shape[row][col] && py+row >= 0 {
				fx := px + col
				fy := py + row
				// フィールドに固定ブロックがあればその色、なければ背景色
				if g.field[fy][fx] != 0 {
					drawCell(fx, fy, colorPieces[g.field[fy][fx]])
				} else {
					eraseCell(fx, fy)
				}
			}
		}
	}
}

// drawScore はスコアとライン数を画面右パネルに描画する。
// tinyfont を使って文字列を直接ハードウェア display に書き込む。
func drawScore(g *Game) {
	d := hardware.Device.GetDisplay()
	white := color.RGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF}
	// SCORE ラベル
	koebiten.DrawFilledRect(nil, scoreX, scoreLabelY-12, 70, 14, colorBg)
	tinyfont.WriteLine(d, &freemono.Regular9pt7b, int16(scoreX), int16(scoreLabelY), "SCORE", white)
	// スコア値
	koebiten.DrawFilledRect(nil, scoreX, scoreValY-12, 70, 14, colorBg)
	tinyfont.WriteLine(d, &freemono.Regular9pt7b, int16(scoreX), int16(scoreValY), strconv.Itoa(g.score), white)
	// LINES ラベル
	koebiten.DrawFilledRect(nil, scoreX, linesLabelY-12, 70, 14, colorBg)
	tinyfont.WriteLine(d, &freemono.Regular9pt7b, int16(scoreX), int16(linesLabelY), "LINES", white)
	// ライン数値
	koebiten.DrawFilledRect(nil, scoreX, linesValY-12, 70, 14, colorBg)
	tinyfont.WriteLine(d, &freemono.Regular9pt7b, int16(scoreX), int16(linesValY), strconv.Itoa(g.lines), white)
}

// drawNextPiece はネクストピースエリアを更新する。
func drawNextPiece(g *Game) {
	// ネクストエリアをクリア
	koebiten.DrawFilledRect(nil, nextX, nextY, 4*cellSize, 4*cellSize, colorBg)
	shape := pieces[g.nextType]
	for row := 0; row < 4; row++ {
		for col := 0; col < 4; col++ {
			if shape[row][col] {
				koebiten.DrawFilledRect(nil,
					nextX+col*cellSize, nextY+row*cellSize,
					cellSize-1, cellSize-1,
					colorPieces[g.nextType+1])
			}
		}
	}
}

// Update は koebiten.Game インターフェースの実装。
func (g *Game) Update() error {
	left   := koebiten.IsKeyPressed(koebiten.KeyLeft)
	right  := koebiten.IsKeyPressed(koebiten.KeyRight)
	down   := koebiten.IsKeyPressed(koebiten.KeyDown)
	rotate := koebiten.IsKeyJustPressed(koebiten.KeyUp) || koebiten.IsKeyJustPressed(koebiten.Key0)
	g.update(left, right, down, rotate)
	return nil
}

// Draw は koebiten.Game インターフェースの実装。
// ClearBuffer は no-op のため、dirty rect で必要な箇所だけ更新する。
func (g *Game) Draw(screen *koebiten.Image) {
	if !g.initialized {
		// 全画面クリア + フィールド全体再描画
		koebiten.DrawFilledRect(nil, 0, 0, screenWidth, screenHeight, colorBg)
		// フィールド枠
		koebiten.DrawFilledRect(nil, fieldX-1, fieldY, 1, fieldRows*cellSize, colorBorder)
		koebiten.DrawFilledRect(nil, fieldX+fieldCols*cellSize, fieldY, 1, fieldRows*cellSize, colorBorder)
		// フィールド内の固定ブロック
		for row := 0; row < fieldRows; row++ {
			for col := 0; col < fieldCols; col++ {
				if g.field[row][col] != 0 {
					drawCell(col, row, colorPieces[g.field[row][col]])
				}
			}
		}
		// ネクストピース
		drawNextPiece(g)
		// スコア
		drawScore(g)
		g.initialized = true
		g.needRedraw = false
		// 現在のピースを描画
		drawPiece(g.pShape, g.pX, g.pY, colorPieces[g.pType+1])
		g.pXPrev = g.pX
		g.pYPrev = g.pY
		g.pShapePrev = g.pShape
		return
	}

	if g.needRedraw {
		// ライン消去後: フィールド全体を再描画
		g.initialized = false
		return
	}

	// 前フレームのピースを消去
	erasePiece(g.pShapePrev, g.pXPrev, g.pYPrev, g)
	// 新しいピースを描画
	drawPiece(g.pShape, g.pX, g.pY, colorPieces[g.pType+1])
	// pPrev を更新
	g.pXPrev = g.pX
	g.pYPrev = g.pY
	g.pShapePrev = g.pShape
	// ネクストピース表示を更新
	drawNextPiece(g)

	// GameOver 表示: 赤い横棒
	if g.state == StateGameOver {
		koebiten.DrawFilledRect(nil, fieldX, screenHeight/2-2, fieldCols*cellSize, 4,
			pixel.NewColor[pixel.RGB565BE](0xFF, 0x00, 0x00))
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
	koebiten.SetWindowTitle("Tetris Demo")
	if err := koebiten.RunGame(NewGame()); err != nil {
		log.Fatal(err)
	}
}
