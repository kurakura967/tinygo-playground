// Package main implements a breakout game demo for Wio Terminal.
package main

const (
	screenWidth  = 320
	screenHeight = 240

	paddleW     = 60
	paddleH     = 8
	paddleY     = 220
	paddleSpeed = 4

	ballSize = 8

	blockCols    = 10
	blockRows    = 5
	blockW       = 28
	blockH       = 12
	blockGapX    = 2
	blockGapY    = 2
	blockOffsetX = 10
	blockOffsetY = 20
)

// GameState はゲームの状態を表す。
type GameState int

const (
	StateReady    GameState = iota // ボタン A 待ち
	StatePlaying                   // ゲーム中
	StateGameOver                  // ボールが落ちた
	StateClear                     // 全ブロック破壊
)

// Game はブレイクアウトゲームの状態を管理する。
type Game struct {
	// パドル
	paddleX     int
	paddleXPrev int

	// ボール
	ballX, ballY         int
	ballXPrev, ballYPrev int
	ballVX, ballVY       int

	// ブロック（true = 生存）
	blocks [blockRows][blockCols]bool
	// 破壊されたブロックを Draw で消去するためのフラグ
	destroyed [blockRows][blockCols]bool

	// 状態
	state       GameState
	initialized bool
}

// NewGame は初期状態の Game を返す。
func NewGame() *Game {
	g := &Game{}
	g.paddleX = screenWidth/2 - paddleW/2
	g.paddleXPrev = g.paddleX
	g.ballX = g.paddleX + paddleW/2 - ballSize/2
	g.ballY = paddleY - ballSize
	g.ballXPrev = g.ballX
	g.ballYPrev = g.ballY
	for r := 0; r < blockRows; r++ {
		for c := 0; c < blockCols; c++ {
			g.blocks[r][c] = true
		}
	}
	g.state = StateReady
	return g
}

// blockRect はブロック (r, c) の矩形座標を返す。
func blockRect(r, c int) (x, y, w, h int) {
	x = blockOffsetX + c*(blockW+blockGapX)
	y = blockOffsetY + r*(blockH+blockGapY)
	return x, y, blockW, blockH
}

// rectsOverlap は2つの矩形が重なるか判定する。
func rectsOverlap(ax, ay, aw, ah, bx, by, bw, bh int) bool {
	return ax < bx+bw && ax+aw > bx && ay < by+bh && ay+ah > by
}

// abs は整数の絶対値を返す。
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// update はボタン入力に基づいてゲーム状態を更新する。
// tinygo 版の Update() と desktop 版のループどちらからも呼ばれる。
func (g *Game) update(left, right, launch bool) {
	// パドル移動（StateReady と StatePlaying で有効）
	if g.state == StateReady || g.state == StatePlaying {
		g.paddleXPrev = g.paddleX
		if left {
			g.paddleX -= paddleSpeed
		}
		if right {
			g.paddleX += paddleSpeed
		}
		// 画面端でクランプ
		if g.paddleX < 0 {
			g.paddleX = 0
		}
		if g.paddleX > screenWidth-paddleW {
			g.paddleX = screenWidth - paddleW
		}
	}

	switch g.state {
	case StateReady:
		// ボールをパドルの中央上に追従
		g.ballXPrev = g.ballX
		g.ballYPrev = g.ballY
		g.ballX = g.paddleX + paddleW/2 - ballSize/2
		g.ballY = paddleY - ballSize
		// launch でゲーム開始
		if launch {
			g.ballVX = 2
			g.ballVY = -3
			g.state = StatePlaying
		}

	case StatePlaying:
		// 前フレーム位置を保存してボールを移動
		g.ballXPrev = g.ballX
		g.ballYPrev = g.ballY
		g.ballX += g.ballVX
		g.ballY += g.ballVY

		// 壁（左右）で反射
		if g.ballX < 0 {
			g.ballX = 0
			g.ballVX = -g.ballVX
		}
		if g.ballX+ballSize > screenWidth {
			g.ballX = screenWidth - ballSize
			g.ballVX = -g.ballVX
		}
		// 壁（上）で反射
		if g.ballY < 0 {
			g.ballY = 0
			g.ballVY = -g.ballVY
		}

		// パドルとの衝突判定
		if g.ballY+ballSize >= paddleY &&
			g.ballX+ballSize > g.paddleX &&
			g.ballX < g.paddleX+paddleW &&
			g.ballVY > 0 {
			g.ballVY = -g.ballVY
		}

		// ボールが画面下に落ちたら GameOver
		if g.ballY > screenHeight {
			g.state = StateGameOver
		}

		// ブロックとの衝突判定
		allDestroyed := true
		for r := 0; r < blockRows; r++ {
			for c := 0; c < blockCols; c++ {
				if !g.blocks[r][c] {
					continue
				}
				allDestroyed = false
				bx, by, bw, bh := blockRect(r, c)
				if rectsOverlap(g.ballX, g.ballY, ballSize, ballSize, bx, by, bw, bh) {
					g.blocks[r][c] = false
					g.destroyed[r][c] = true
					// どちらの面から当たったかで反射方向を決める
					ballCX := g.ballX + ballSize/2
					ballCY := g.ballY + ballSize/2
					blkCX := bx + blockW/2
					blkCY := by + blockH/2
					dx := ballCX - blkCX
					dy := ballCY - blkCY
					if abs(dx)*blockH > abs(dy)*blockW {
						g.ballVX = -g.ballVX
					} else {
						g.ballVY = -g.ballVY
					}
				}
			}
		}

		// 全ブロックが破壊されたか確認
		if allDestroyed {
			g.state = StateClear
		}

	case StateGameOver, StateClear:
		// ボタン A でリスタート
		if launch {
			g.paddleX = screenWidth/2 - paddleW/2
			g.paddleXPrev = g.paddleX
			g.ballX = g.paddleX + paddleW/2 - ballSize/2
			g.ballY = paddleY - ballSize
			g.ballXPrev = g.ballX
			g.ballYPrev = g.ballY
			g.ballVX = 0
			g.ballVY = 0
			for r := 0; r < blockRows; r++ {
				for c := 0; c < blockCols; c++ {
					g.blocks[r][c] = true
					g.destroyed[r][c] = false
				}
			}
			g.state = StateReady
			g.initialized = false
		}
	}
}
