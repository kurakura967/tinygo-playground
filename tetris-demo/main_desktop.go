//go:build !tinygo

// tetris-demo デスクトップ版。wiodisplay エミュレータ（wio-emu）と通信してテトリスゲームを表示する。
//
// 使い方:
//  1. wio-emu サーバーを起動する:
//     cd cmd/wio-emu && go run .
//  2. このアプリを起動する:
//     go run .
package main

import (
	"image/color"
	"strconv"
	"time"

	"github.com/kurakura967/wiodisplay/driver/ili9341"
	"github.com/kurakura967/wiodisplay/initdisplay"
	"github.com/kurakura967/wiodisplay/machine"
	"tinygo.org/x/tinyfont"
	"tinygo.org/x/tinyfont/freemono"
)

var (
	deskBg      = color.RGBA{R: 0x00, G: 0x00, B: 0x00, A: 0xFF} // 黒
	deskBorder  = color.RGBA{R: 0x40, G: 0x40, B: 0x40, A: 0xFF} // グレー
	// ピース種別ごとの色（0=空, 1-7=I/O/T/S/Z/J/L）
	deskPieces = [8]color.RGBA{
		{R: 0x00, G: 0x00, B: 0x00, A: 0xFF}, // 0: 空（黒）
		{R: 0x00, G: 0xFF, B: 0xFF, A: 0xFF}, // 1: I シアン
		{R: 0xFF, G: 0xFF, B: 0x00, A: 0xFF}, // 2: O 黄
		{R: 0x80, G: 0x00, B: 0xFF, A: 0xFF}, // 3: T 紫
		{R: 0x00, G: 0xFF, B: 0x00, A: 0xFF}, // 4: S 緑
		{R: 0xFF, G: 0x00, B: 0x00, A: 0xFF}, // 5: Z 赤
		{R: 0x00, G: 0x00, B: 0xFF, A: 0xFF}, // 6: J 青
		{R: 0xFF, G: 0x80, B: 0x00, A: 0xFF}, // 7: L オレンジ
	}
)

// deskDrawCell はフィールド座標 (fx, fy) にセルを描画する（デスクトップ版）。
func deskDrawCell(display interface {
	FillRectangle(x, y, width, height int16, c color.RGBA)
}, fx, fy int, c color.RGBA) {
	display.FillRectangle(int16(fieldX+fx*cellSize), int16(fieldY+fy*cellSize), cellSize-1, cellSize-1, c)
}

// deskEraseCell はフィールド座標 (fx, fy) のセルを背景色で消去する（デスクトップ版）。
func deskEraseCell(display interface {
	FillRectangle(x, y, width, height int16, c color.RGBA)
}, fx, fy int) {
	display.FillRectangle(int16(fieldX+fx*cellSize), int16(fieldY+fy*cellSize), cellSize, cellSize, deskBg)
}

// drawScore はスコアとライン数を画面右パネルに描画する（デスクトップ版）。
func drawScore(display *ili9341.Device, g *Game) {
	white := color.RGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF}
	// SCORE ラベル
	display.FillRectangle(int16(scoreX), int16(scoreLabelY-12), 70, 14, deskBg)
	tinyfont.WriteLine(display, &freemono.Regular9pt7b, int16(scoreX), int16(scoreLabelY), "SCORE", white)
	// スコア値
	display.FillRectangle(int16(scoreX), int16(scoreValY-12), 70, 14, deskBg)
	tinyfont.WriteLine(display, &freemono.Regular9pt7b, int16(scoreX), int16(scoreValY), strconv.Itoa(g.score), white)
	// LINES ラベル
	display.FillRectangle(int16(scoreX), int16(linesLabelY-12), 70, 14, deskBg)
	tinyfont.WriteLine(display, &freemono.Regular9pt7b, int16(scoreX), int16(linesLabelY), "LINES", white)
	// ライン数値
	display.FillRectangle(int16(scoreX), int16(linesValY-12), 70, 14, deskBg)
	tinyfont.WriteLine(display, &freemono.Regular9pt7b, int16(scoreX), int16(linesValY), strconv.Itoa(g.lines), white)
}

func main() {
	display := initdisplay.InitDisplay()

	game := NewGame()

	// エッジ検出用
	prevBtnA := false
	prevUp := false

	for {
		// 入力処理（デスクトップ版: machine.Pin.Get() で RPC 経由取得）
		left     := !machine.WIO_5S_LEFT.Get()  // アクティブロー
		right    := !machine.WIO_5S_RIGHT.Get()
		down     := !machine.WIO_5S_DOWN.Get()
		rawBtnA  := !machine.WIO_KEY_A.Get()
		rawUp    := !machine.WIO_5S_UP.Get()
		// IsKeyJustPressed 相当: false→true の遷移時のみ true
		rotate := (rawBtnA && !prevBtnA) || (rawUp && !prevUp)
		prevBtnA = rawBtnA
		prevUp = rawUp

		game.update(left, right, down, rotate)

		if !game.initialized {
			// 全画面クリア + フィールド全体再描画
			display.FillScreen(deskBg)
			// フィールド枠
			display.FillRectangle(int16(fieldX-1), int16(fieldY), 1, int16(fieldRows*cellSize), deskBorder)
			display.FillRectangle(int16(fieldX+fieldCols*cellSize), int16(fieldY), 1, int16(fieldRows*cellSize), deskBorder)
			// フィールド内の固定ブロック
			for row := 0; row < fieldRows; row++ {
				for col := 0; col < fieldCols; col++ {
					if game.field[row][col] != 0 {
						display.FillRectangle(
							int16(fieldX+col*cellSize), int16(fieldY+row*cellSize),
							cellSize-1, cellSize-1,
							deskPieces[game.field[row][col]])
					}
				}
			}
			// ネクストピース
			display.FillRectangle(int16(nextX), int16(nextY), 4*cellSize, 4*cellSize, deskBg)
			nShape := pieces[game.nextType]
			for row := 0; row < 4; row++ {
				for col := 0; col < 4; col++ {
					if nShape[row][col] {
						display.FillRectangle(
							int16(nextX+col*cellSize), int16(nextY+row*cellSize),
							cellSize-1, cellSize-1,
							deskPieces[game.nextType+1])
					}
				}
			}
			// スコア
			drawScore(display, game)
			game.initialized = true
			game.needRedraw = false
			// 現在のピースを描画
			pShape := game.pShape
			for row := 0; row < 4; row++ {
				for col := 0; col < 4; col++ {
					if pShape[row][col] && game.pY+row >= 0 {
						display.FillRectangle(
							int16(fieldX+(game.pX+col)*cellSize),
							int16(fieldY+(game.pY+row)*cellSize),
							cellSize-1, cellSize-1,
							deskPieces[game.pType+1])
					}
				}
			}
			game.pXPrev = game.pX
			game.pYPrev = game.pY
			game.pShapePrev = game.pShape

			time.Sleep(16 * time.Millisecond)
			continue
		}

		if game.needRedraw {
			// ライン消去後: 再初期化フラグを立てて次のフレームで全体再描画
			game.initialized = false
			time.Sleep(16 * time.Millisecond)
			continue
		}

		// 前フレームのピースを消去
		prevShape := game.pShapePrev
		for row := 0; row < 4; row++ {
			for col := 0; col < 4; col++ {
				if prevShape[row][col] && game.pYPrev+row >= 0 {
					fx := game.pXPrev + col
					fy := game.pYPrev + row
					if game.field[fy][fx] != 0 {
						display.FillRectangle(
							int16(fieldX+fx*cellSize), int16(fieldY+fy*cellSize),
							cellSize-1, cellSize-1,
							deskPieces[game.field[fy][fx]])
					} else {
						display.FillRectangle(
							int16(fieldX+fx*cellSize), int16(fieldY+fy*cellSize),
							cellSize, cellSize,
							deskBg)
					}
				}
			}
		}

		// 新しいピースを描画
		curShape := game.pShape
		for row := 0; row < 4; row++ {
			for col := 0; col < 4; col++ {
				if curShape[row][col] && game.pY+row >= 0 {
					display.FillRectangle(
						int16(fieldX+(game.pX+col)*cellSize),
						int16(fieldY+(game.pY+row)*cellSize),
						cellSize-1, cellSize-1,
						deskPieces[game.pType+1])
				}
			}
		}

		// pPrev を更新
		game.pXPrev = game.pX
		game.pYPrev = game.pY
		game.pShapePrev = game.pShape
		// ネクストピース表示を更新
		display.FillRectangle(int16(nextX), int16(nextY), 4*cellSize, 4*cellSize, deskBg)
		nShape := pieces[game.nextType]
		for row := 0; row < 4; row++ {
			for col := 0; col < 4; col++ {
				if nShape[row][col] {
					display.FillRectangle(
						int16(nextX+col*cellSize), int16(nextY+row*cellSize),
						cellSize-1, cellSize-1,
						deskPieces[game.nextType+1])
				}
			}
		}

		// GameOver 表示: 赤い横棒
		if game.state == StateGameOver {
			display.FillRectangle(int16(fieldX), int16(screenHeight/2-2), int16(fieldCols*cellSize), 4,
				color.RGBA{R: 0xFF, G: 0x00, B: 0x00, A: 0xFF})
		}

		time.Sleep(16 * time.Millisecond)
	}
}
