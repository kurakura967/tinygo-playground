//go:build !tinygo

// emu-demo デスクトップ版。wiodisplay エミュレータ（wio-emu）と通信して矩形アニメーションを表示する。
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

func main() {
	display := initdisplay.InitDisplay()

	// 背景を暗い青でクリア
	display.FillScreen(color.RGBA{R: 0x00, G: 0x00, B: 0x20, A: 0xFF})

	game := NewGame()

	// 初回: 矩形全体を描画
	display.FillRectangle(int16(game.x), int16(game.y), int16(rectSize), int16(rectSize), color.RGBA{R: 0xFF, G: 0x80, B: 0x00, A: 0xFF})

	for {
		// 入力処理（デスクトップ版: machine.Pin.Get() で RPC 経由取得）
		btnA := !machine.WIO_KEY_A.Get() // アクティブロー
		btnB := !machine.WIO_KEY_B.Get()
		btnC := !machine.WIO_KEY_C.Get()
		left := !machine.WIO_5S_LEFT.Get()
		right := !machine.WIO_5S_RIGHT.Get()
		up := !machine.WIO_5S_UP.Get()
		down := !machine.WIO_5S_DOWN.Get()

		// ボタン C で一時停止
		if !btnC {
			game.update(btnA, btnB, btnC, left, right, up, down)
		}

		// 矩形の色
		rectColor := color.RGBA{R: 0xFF, G: 0x80, B: 0x00, A: 0xFF} // オレンジ
		if btnA {
			rectColor = color.RGBA{R: 0xFF, G: 0x00, B: 0x00, A: 0xFF} // 赤
		} else if btnB {
			rectColor = color.RGBA{R: 0x00, G: 0xFF, B: 0x00, A: 0xFF} // 緑
		} else if btnC {
			rectColor = color.RGBA{R: 0x00, G: 0x00, B: 0xFF, A: 0xFF} // 青
		}

		// 新しい位置に矩形を描画（先に描く）
		display.FillRectangle(int16(game.x), int16(game.y), int16(rectSize), int16(rectSize), rectColor)
		// 旧位置の trailing edge を消去（後から消す）
		for _, r := range game.dirtyRects() {
			display.FillRectangle(int16(r[0]), int16(r[1]), int16(r[2]), int16(r[3]), color.RGBA{R: 0x00, G: 0x00, B: 0x20, A: 0xFF})
		}

		time.Sleep(32 * time.Millisecond)
	}
}
