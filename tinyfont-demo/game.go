// Package main は tinyfont を使ったテキスト表示デモ。
//
// このファイルはビルドタグなし。
// tinyfont はハードウェア依存がないため、デスクトップ Go でもそのままコンパイルできる。
// これが「TinyGo エコシステムのライブラリでも build tag 不要なものがある」という発表ポイント。
package main

import (
	"image/color"
	"strconv"

	"tinygo.org/x/tinyfont"
	"tinygo.org/x/tinyfont/freemono"
)

const (
	screenWidth  = 320
	screenHeight = 240
)

// TextDisplay は tinyfont が必要とする最小インターフェース。
// tinygo 版（koebiten.Displayer = *WioDisplay/ILI9341）・
// desktop 版（*ili9341.Device = RPC 経由）どちらも満たす。
type TextDisplay interface {
	SetPixel(x, y int16, c color.RGBA)
	Size() (x, y int16)
	Display() error
}

// page は1ページ分の表示コンテンツ。
type page struct {
	title string
	lines []string
}

var pages = []page{
	{
		title: "tinyfont",
		lines: []string{"text rendering", "for TinyGo"},
	},
	{
		title: "No build tag!",
		lines: []string{"Works on both", "desktop & device"},
	},
	{
		title: "game.go",
		lines: []string{"import tinyfont", "// no build tag :)"},
	},
	{
		title: "Go Conf 2026",
		lines: []string{"Sep. 11", "Short Session"},
	},
}

var (
	colorWhite = color.RGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF}
	colorCyan  = color.RGBA{R: 0x00, G: 0xFF, B: 0xFF, A: 0xFF}
	colorGray  = color.RGBA{R: 0x60, G: 0x60, B: 0x60, A: 0xFF}
)

// Game はデモのページ状態を管理する。
type Game struct {
	pageIdx int
	changed bool
}

func NewGame() *Game {
	return &Game{changed: true}
}

// update はボタン入力でページを切り替える。
// btnA は IsKeyJustPressed 相当（押した瞬間のみ true）で渡すこと。
func (g *Game) update(btnA bool) {
	if btnA {
		g.pageIdx = (g.pageIdx + 1) % len(pages)
		g.changed = true
	}
}

// drawText は現在のページを tinyfont でレンダリングする。
// TextDisplay を受け取るため、ビルドタグなしで動作する。
// tinygo / desktop どちらの display も TextDisplay を満たすので、
// このファイルから直接呼び出せる。
func drawText(d TextDisplay, g *Game) {
	p := pages[g.pageIdx]

	// タイトル (12pt・シアン)
	tinyfont.WriteLine(d, &freemono.Regular12pt7b, 10, 35, p.title, colorCyan)

	// 本文 (9pt・白)
	for i, line := range p.lines {
		tinyfont.WriteLine(d, &freemono.Regular9pt7b, 10, int16(70+i*22), line, colorWhite)
	}

	// 操作ヒントとページ番号 (9pt・グレー)
	tinyfont.WriteLine(d, &freemono.Regular9pt7b, 10, 225, "[ A ] next", colorGray)
	pageNum := strconv.Itoa(g.pageIdx+1) + "/" + strconv.Itoa(len(pages))
	tinyfont.WriteLine(d, &freemono.Regular9pt7b, int16(screenWidth-55), 225, pageNum, colorGray)
}
