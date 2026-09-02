// Package main implements a bouncing rectangle animation demo.
package main

const (
	screenWidth  = 320
	screenHeight = 240
	rectSize     = 40
	speed        = 3
)

// Game はバウンシング矩形のゲーム状態を管理する。
type Game struct {
	x, y        int
	prevX       int
	prevY       int
	vx          int
	vy          int
	initialized bool
}

// NewGame は初期状態の Game を返す。
func NewGame() *Game {
	x := screenWidth/2 - rectSize/2
	y := screenHeight/2 - rectSize/2
	return &Game{
		x:     x,
		y:     y,
		prevX: x,
		prevY: y,
		vx:    speed,
		vy:    speed,
	}
}

// dirtyRects は旧矩形から新矩形に移動したときに消去が必要な領域を返す。
// 最大2つの矩形（水平方向の帯・垂直方向の帯）を返す。
// 戻り値: [][4]int{{x, y, w, h}, ...}
func (g *Game) dirtyRects() [][4]int {
	var rects [][4]int

	// X方向の差分
	if g.x > g.prevX {
		// 右移動: 旧左端から移動量の幅の帯を消す
		rects = append(rects, [4]int{g.prevX, g.prevY, g.x - g.prevX, rectSize})
	} else if g.x < g.prevX {
		// 左移動: 新矩形右端から移動量の幅の帯を消す
		rects = append(rects, [4]int{g.x + rectSize, g.prevY, g.prevX - g.x, rectSize})
	}

	// Y方向の差分
	if g.y > g.prevY {
		// 下移動: 旧上端から移動量の高さの帯を消す
		rects = append(rects, [4]int{g.prevX, g.prevY, rectSize, g.y - g.prevY})
	} else if g.y < g.prevY {
		// 上移動: 新矩形下端から移動量の高さの帯を消す
		rects = append(rects, [4]int{g.prevX, g.y + rectSize, rectSize, g.prevY - g.y})
	}

	return rects
}

// leadingRects は新矩形にあって旧矩形になかった領域（描画が必要な新しい部分）を返す。
// 戻り値: [][4]int{{x, y, w, h}, ...}
func (g *Game) leadingRects() [][4]int {
	var rects [][4]int
	dx := g.x - g.prevX
	dy := g.y - g.prevY
	if dx > 0 {
		// 右移動: 新rect右端の帯
		rects = append(rects, [4]int{g.prevX + rectSize, g.y, dx, rectSize})
	} else if dx < 0 {
		// 左移動: 新rect左端の帯
		rects = append(rects, [4]int{g.x, g.y, -dx, rectSize})
	}
	if dy > 0 {
		// 下移動: 新rect下端の帯
		rects = append(rects, [4]int{g.x, g.prevY + rectSize, rectSize, dy})
	} else if dy < 0 {
		// 上移動: 新rect上端の帯
		rects = append(rects, [4]int{g.x, g.y, rectSize, -dy})
	}
	return rects
}

// update はボタン入力に基づいて矩形の位置を更新する。
// tinygo 版の Update() と desktop 版のループどちらからも呼ばれる。
func (g *Game) update(btnA, btnB, btnC, left, right, up, down bool) {
	// 現在位置を前フレーム位置として保存
	g.prevX = g.x
	g.prevY = g.y

	// ボタン A でオートバウンスをリセット
	if btnA {
		g.vx = speed
		g.vy = speed
	}

	// ボタン B で速度を上げる
	if btnB {
		g.x += g.vx * 2
		g.y += g.vy * 2
	} else {
		g.x += g.vx
		g.y += g.vy
	}

	// ジョイスティックで方向制御（オートバウンスに加算）
	if left {
		g.x -= speed
	}
	if right {
		g.x += speed
	}
	if up {
		g.y -= speed
	}
	if down {
		g.y += speed
	}

	// 壁での反射
	if g.x < 0 {
		g.x = 0
		g.vx = -g.vx
	}
	if g.x+rectSize > screenWidth {
		g.x = screenWidth - rectSize
		g.vx = -g.vx
	}
	if g.y < 0 {
		g.y = 0
		g.vy = -g.vy
	}
	if g.y+rectSize > screenHeight {
		g.y = screenHeight - rectSize
		g.vy = -g.vy
	}
}
