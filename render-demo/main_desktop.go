//go:build !tinygo

// render-demo デスクトップ版。wiodisplay エミュレータ（wio-emu）と通信して矩形アニメーションを表示する。
// ボタン A でナイーブモード（毎フレーム全画面クリア）と最適化モード（dirty rect のみ更新）を切り替え。
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
	bgColor        = color.RGBA{R: 0x00, G: 0x00, B: 0x20, A: 0xFF} // 暗い青
	rectColor      = color.RGBA{R: 0xFF, G: 0x80, B: 0x00, A: 0xFF} // オレンジ
	naiveColor     = color.RGBA{R: 0xFF, G: 0x00, B: 0x00, A: 0xFF} // 赤 - ナイーブモードインジケーター
	optimizedColor = color.RGBA{R: 0x00, G: 0xFF, B: 0x00, A: 0xFF} // 緑 - 最適化モードインジケーター
)

func main() {
	display := initdisplay.InitDisplay()

	// 背景を暗い青でクリア
	display.FillScreen(bgColor)

	game := NewGame()

	// 初回: 矩形全体を描画
	display.FillRectangle(int16(game.x), int16(game.y), int16(rectSize), int16(rectSize), rectColor)
	// ナイーブモードインジケーター（左上 10x10）
	display.FillRectangle(0, 0, 10, 10, naiveColor)

	prevBtnA := false
	for {
		// 入力処理（デスクトップ版: machine.Pin.Get() で RPC 経由取得）
		rawBtnA := !machine.WIO_KEY_A.Get() // アクティブロー - モード切り替え
		toggleMode := rawBtnA && !prevBtnA  // 押した瞬間のみ true
		prevBtnA = rawBtnA

		left := !machine.WIO_5S_LEFT.Get()
		right := !machine.WIO_5S_RIGHT.Get()
		up := !machine.WIO_5S_UP.Get()
		down := !machine.WIO_5S_DOWN.Get()

		prevMode := game.mode
		game.update(toggleMode, left, right, up, down)

		switch game.mode {
		case ModeNaive:
			// ナイーブモード: 毎フレーム全画面クリアしてから矩形を描画
			display.FillScreen(bgColor)
			display.FillRectangle(int16(game.x), int16(game.y), int16(rectSize), int16(rectSize), rectColor)
			// モードインジケーター（左上 10x10）: 赤
			display.FillRectangle(0, 0, 10, 10, naiveColor)

		case ModeOptimized:
			// 最適化モード: dirty rect のみ更新
			if !game.initialized || prevMode != game.mode {
				// 初回またはモード切り替え直後: 全画面クリア
				display.FillScreen(bgColor)
				display.FillRectangle(int16(game.x), int16(game.y), int16(rectSize), int16(rectSize), rectColor)
				game.initialized = true
				// モードインジケーター（左上 10x10）: 緑
				display.FillRectangle(0, 0, 10, 10, optimizedColor)
			} else {
				// 新しい位置に矩形を描画（先に描く）
				display.FillRectangle(int16(game.x), int16(game.y), int16(rectSize), int16(rectSize), rectColor)
				// 旧位置の trailing edge を消去（後から消す）
				for _, r := range game.dirtyRects() {
					display.FillRectangle(int16(r[0]), int16(r[1]), int16(r[2]), int16(r[3]), bgColor)
				}
				// モードインジケーター（左上 10x10）: 緑
				display.FillRectangle(0, 0, 10, 10, optimizedColor)
			}
		}

		time.Sleep(32 * time.Millisecond)
	}
}
