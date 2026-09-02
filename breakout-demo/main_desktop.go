//go:build !tinygo

// breakout-demo デスクトップ版。wiodisplay エミュレータ（wio-emu）と通信してブレイクアウトゲームを表示する。
//
// 使い方:
//  1. wio-emu サーバーを起動する:
//     cd cmd/wio-emu && go run .
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
	bgColor     = color.RGBA{R: 0x00, G: 0x00, B: 0x00, A: 0xFF} // 黒
	paddleColor = color.RGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF} // 白
	ballColor   = color.RGBA{R: 0xFF, G: 0xFF, B: 0x00, A: 0xFF} // 黄
	blockColors = [blockRows]color.RGBA{
		{R: 0xFF, G: 0x00, B: 0x00, A: 0xFF}, // 赤
		{R: 0xFF, G: 0x80, B: 0x00, A: 0xFF}, // オレンジ
		{R: 0xFF, G: 0xFF, B: 0x00, A: 0xFF}, // 黄
		{R: 0x00, G: 0xFF, B: 0x00, A: 0xFF}, // 緑
		{R: 0x00, G: 0x80, B: 0xFF, A: 0xFF}, // 青
	}
)

func main() {
	display := initdisplay.InitDisplay()

	game := NewGame()

	// エッジ検出用
	prevBtnA := false

	for {
		// 入力処理（デスクトップ版: machine.Pin.Get() で RPC 経由取得）
		left  := !machine.WIO_5S_LEFT.Get()  // アクティブロー
		right := !machine.WIO_5S_RIGHT.Get()
		rawBtnA := !machine.WIO_KEY_A.Get()
		// IsKeyJustPressed 相当: false→true の遷移時のみ true
		launch := rawBtnA && !prevBtnA
		prevBtnA = rawBtnA

		game.update(left, right, launch)

		if !game.initialized {
			// 初回: 全画面を黒でクリアして全要素を描画
			display.FillScreen(bgColor)
			// 全ブロックを描画
			for r := 0; r < blockRows; r++ {
				for c := 0; c < blockCols; c++ {
					if game.blocks[r][c] {
						bx, by, bw, bh := blockRect(r, c)
						display.FillRectangle(int16(bx), int16(by), int16(bw), int16(bh), blockColors[r])
					}
				}
			}
			// パドルを描画
			display.FillRectangle(int16(game.paddleX), int16(paddleY), int16(paddleW), int16(paddleH), paddleColor)
			// ボールを描画
			display.FillRectangle(int16(game.ballX), int16(game.ballY), int16(ballSize), int16(ballSize), ballColor)
			game.initialized = true

			time.Sleep(16 * time.Millisecond)
			continue
		}

		// 破壊されたブロックを消去
		for r := 0; r < blockRows; r++ {
			for c := 0; c < blockCols; c++ {
				if game.destroyed[r][c] {
					bx, by, bw, bh := blockRect(r, c)
					display.FillRectangle(int16(bx), int16(by), int16(bw), int16(bh), bgColor)
					game.destroyed[r][c] = false
				}
			}
		}

		// パドルの旧位置を消去して新位置に描画
		if game.paddleX != game.paddleXPrev {
			display.FillRectangle(int16(game.paddleXPrev), int16(paddleY), int16(paddleW), int16(paddleH), bgColor)
			display.FillRectangle(int16(game.paddleX), int16(paddleY), int16(paddleW), int16(paddleH), paddleColor)
		}

		// ボールの旧位置を消去して新位置に描画
		display.FillRectangle(int16(game.ballXPrev), int16(game.ballYPrev), int16(ballSize), int16(ballSize), bgColor)
		display.FillRectangle(int16(game.ballX), int16(game.ballY), int16(ballSize), int16(ballSize), ballColor)

		// GameOver / Clear の状態表示（横棒で示す）
		switch game.state {
		case StateGameOver:
			display.FillRectangle(60, 108, 200, 4, color.RGBA{R: 0xFF, G: 0x00, B: 0x00, A: 0xFF})
			display.FillRectangle(60, 116, 200, 4, color.RGBA{R: 0xFF, G: 0x00, B: 0x00, A: 0xFF})
		case StateClear:
			display.FillRectangle(60, 108, 200, 4, color.RGBA{R: 0x00, G: 0xFF, B: 0x00, A: 0xFF})
			display.FillRectangle(60, 116, 200, 4, color.RGBA{R: 0x00, G: 0xFF, B: 0x00, A: 0xFF})
		}

		time.Sleep(16 * time.Millisecond)
	}
}
