// Package main はボタン操作で動くメニュー UI デモ。
//
// このファイルはビルドタグなし。
// tinyfont を使ったメニュー描画ロジックを tinygo / desktop 共通で定義している。
package main

import (
	"image/color"

	"tinygo.org/x/tinyfont"
	"tinygo.org/x/tinyfont/freemono"
)

const (
	screenWidth  = 320
	screenHeight = 240
)

// TextDisplay は tinyfont が必要とする最小インターフェース。
// tinygo 版（koebiten.Displayer）・desktop 版（*ili9341.Device）どちらも満たす。
type TextDisplay interface {
	SetPixel(x, y int16, c color.RGBA)
	Size() (x, y int16)
	Display() error
}

// MenuState はメニューの表示状態を表す。
type MenuState int

const (
	StateMenu   MenuState = iota // メニュー一覧表示
	StateDetail                  // 選択項目の詳細表示
)

// MenuItem はメニューの1項目を表す。
type MenuItem struct {
	label  string   // メニュー一覧に表示するラベル
	detail []string // 詳細画面に表示するテキスト行
}

var menuItems = []MenuItem{
	{
		label:  "Hello, World!",
		detail: []string{"Hello,", "World!"},
	},
	{
		label:  "About TinyGo",
		detail: []string{"TinyGo v0.35", "for Wio Terminal"},
	},
	{
		label:  "About wiodisplay",
		detail: []string{"Wio Terminal", "Emulator", "by kurakura967"},
	},
	{
		label:  "Go Conf 2026",
		detail: []string{"Sep. 11, 2026", "Short Session", "20 min"},
	},
}

var (
	colorWhite = color.RGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF}
	colorCyan  = color.RGBA{R: 0x00, G: 0xFF, B: 0xFF, A: 0xFF}
	colorGray  = color.RGBA{R: 0x60, G: 0x60, B: 0x60, A: 0xFF}
)

// Game はメニューのカーソル位置と状態を管理する。
type Game struct {
	cursor  int
	state   MenuState
	changed bool
}

func NewGame() *Game {
	return &Game{changed: true}
}

// update はボタン入力に基づいてメニュー状態を更新する。
// up / down / btnA / btnB は IsKeyJustPressed 相当（押した瞬間のみ true）で渡すこと。
func (g *Game) update(up, down, btnA, btnB bool) {
	switch g.state {
	case StateMenu:
		if up && g.cursor > 0 {
			g.cursor--
			g.changed = true
		}
		if down && g.cursor < len(menuItems)-1 {
			g.cursor++
			g.changed = true
		}
		if btnA {
			g.state = StateDetail
			g.changed = true
		}
	case StateDetail:
		if btnB {
			g.state = StateMenu
			g.changed = true
		}
	}
}

// drawMenu はメニュー一覧画面を描画する。
func drawMenu(d TextDisplay, g *Game) {
	tinyfont.WriteLine(d, &freemono.Regular12pt7b, 10, 28, "Menu", colorCyan)

	for i, item := range menuItems {
		y := int16(62 + i*22)
		if i == g.cursor {
			tinyfont.WriteLine(d, &freemono.Regular9pt7b, 10, y, "> "+item.label, colorCyan)
		} else {
			tinyfont.WriteLine(d, &freemono.Regular9pt7b, 10, y, "  "+item.label, colorWhite)
		}
	}

	tinyfont.WriteLine(d, &freemono.Regular9pt7b, 10, 225, "[A]select [U/D]move", colorGray)
}

// drawDetail は選択した項目の詳細画面を描画する。
func drawDetail(d TextDisplay, g *Game) {
	item := menuItems[g.cursor]

	tinyfont.WriteLine(d, &freemono.Regular12pt7b, 10, 28, item.label, colorCyan)

	for i, line := range item.detail {
		tinyfont.WriteLine(d, &freemono.Regular9pt7b, 10, int16(68+i*22), line, colorWhite)
	}

	tinyfont.WriteLine(d, &freemono.Regular9pt7b, 10, 225, "[ B ] back", colorGray)
}
