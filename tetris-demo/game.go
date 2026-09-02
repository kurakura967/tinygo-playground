// Package main implements a tetris game demo for Wio Terminal.
package main

const (
	screenWidth  = 320
	screenHeight = 240

	fieldCols = 10
	fieldRows = 20
	cellSize  = 12

	fieldX = 4  // フィールド左端 X
	fieldY = 0  // フィールド上端 Y

	nextX = 150 // ネクストピース表示左端 X
	nextY = 20  // ネクストピース表示上端 Y

	// スコア表示エリア（右パネル）
	scoreX      = 155
	scoreLabelY = 96  // "SCORE" ベースライン Y
	scoreValY   = 114 // スコア値ベースライン Y
	linesLabelY = 138 // "LINES" ベースライン Y
	linesValY   = 156 // ライン数値ベースライン Y
)

// pieces[type][row][col]
var pieces = [7][4][4]bool{
	// I
	{{false, false, false, false}, {true, true, true, true}, {false, false, false, false}, {false, false, false, false}},
	// O
	{{false, true, true, false}, {false, true, true, false}, {false, false, false, false}, {false, false, false, false}},
	// T
	{{false, true, false, false}, {true, true, true, false}, {false, false, false, false}, {false, false, false, false}},
	// S
	{{false, true, true, false}, {true, true, false, false}, {false, false, false, false}, {false, false, false, false}},
	// Z
	{{true, true, false, false}, {false, true, true, false}, {false, false, false, false}, {false, false, false, false}},
	// J
	{{true, false, false, false}, {true, true, true, false}, {false, false, false, false}, {false, false, false, false}},
	// L
	{{false, false, true, false}, {true, true, true, false}, {false, false, false, false}, {false, false, false, false}},
}

// GameState はゲームの状態を表す。
type GameState int

const (
	StateReady   GameState = iota
	StatePlaying
	StateGameOver
)

// Game はテトリスゲームの状態を管理する。
type Game struct {
	// フィールド: 0=空, 1-7=ピース種別（色管理用）
	field [fieldRows][fieldCols]int

	// 現在のピース
	pType  int        // 0-6
	pX, pY int        // フィールド座標
	pShape [4][4]bool // 現在の回転状態

	// 前フレームのピース位置（dirty rect 用）
	pXPrev, pYPrev int
	pShapePrev     [4][4]bool

	// ネクストピース
	nextType int

	// タイマー
	fallTimer int
	fallSpeed int // 落下間隔（フレーム数）。初期値30、最小5

	// DAS（Delayed Auto Shift）
	dasDir   int // -1=left, 0=none, 1=right
	dasTimer int

	// スコア
	score int
	lines int

	// 状態
	state       GameState
	initialized bool
	needRedraw  bool // ライン消去後の全体再描画フラグ
}

// seed は疑似乱数生成用のシード。
var seed uint32 = 12345

// randomPieceType は 0-6 のランダムなピース種別を返す。
func randomPieceType() int {
	seed = seed*1664525 + 1013904223
	return int(seed>>25) % 7
}

// NewGame は初期状態の Game を返す。
func NewGame() *Game {
	g := &Game{
		fallSpeed: 30,
		state:     StateReady,
	}
	g.nextType = randomPieceType()
	return g
}

// rotateCW は [4][4]bool を時計回りに90度回転する。
func rotateCW(s [4][4]bool) [4][4]bool {
	var r [4][4]bool
	for row := 0; row < 4; row++ {
		for col := 0; col < 4; col++ {
			r[col][3-row] = s[row][col]
		}
	}
	return r
}

// collides は shape をフィールドの (px, py) に置いたとき衝突するか判定する。
func (g *Game) collides(shape [4][4]bool, px, py int) bool {
	for row := 0; row < 4; row++ {
		for col := 0; col < 4; col++ {
			if !shape[row][col] {
				continue
			}
			fx := px + col
			fy := py + row
			if fx < 0 || fx >= fieldCols || fy >= fieldRows {
				return true
			}
			if fy >= 0 && g.field[fy][fx] != 0 {
				return true
			}
		}
	}
	return false
}

// spawnPiece は新しいピースを生成する。
func (g *Game) spawnPiece() {
	g.pType = g.nextType
	g.nextType = randomPieceType()
	g.pShape = pieces[g.pType]
	g.pX = fieldCols/2 - 2
	g.pY = -1 // 画面上端より上からスポーン
	g.pXPrev = g.pX
	g.pYPrev = g.pY
	g.pShapePrev = g.pShape
	if g.collides(g.pShape, g.pX, g.pY) {
		g.state = StateGameOver
	}
}

// lockPiece は現在のピースをフィールドに固定する。
func (g *Game) lockPiece() {
	for row := 0; row < 4; row++ {
		for col := 0; col < 4; col++ {
			if g.pShape[row][col] {
				fy := g.pY + row
				fx := g.pX + col
				if fy >= 0 && fy < fieldRows {
					g.field[fy][fx] = g.pType + 1
				}
			}
		}
	}
	g.clearLines()
	g.spawnPiece()
}

// clearLines は揃った行を消去する。
func (g *Game) clearLines() {
	cleared := 0
	for row := fieldRows - 1; row >= 0; row-- {
		full := true
		for col := 0; col < fieldCols; col++ {
			if g.field[row][col] == 0 {
				full = false
				break
			}
		}
		if full {
			// この行を消して上をずらす
			for r := row; r > 0; r-- {
				g.field[r] = g.field[r-1]
			}
			g.field[0] = [fieldCols]int{}
			cleared++
			row++ // 同じ行を再チェック
		}
	}
	if cleared > 0 {
		// スコア加算（1=100, 2=300, 3=500, 4+=800）
		switch cleared {
		case 1:
			g.score += 100
		case 2:
			g.score += 300
		case 3:
			g.score += 500
		default:
			g.score += 800
		}
		g.lines += cleared
		g.needRedraw = true
		// 落下速度を上げる
		if g.fallSpeed > 5 {
			g.fallSpeed -= cleared
		}
	}
}

// update はボタン入力に基づいてゲーム状態を更新する。
// rotate は IsKeyJustPressed で取得すること。
func (g *Game) update(left, right, down, rotate bool) {
	switch g.state {
	case StateReady:
		if rotate || left || right || down {
			// 任意のキーでゲーム開始
			g.spawnPiece()
			g.state = StatePlaying
			g.initialized = false
		}
	case StatePlaying:
		// DAS 処理（左右の押しっぱなし）
		if left && !right {
			if g.dasDir != -1 {
				g.dasDir = -1
				g.dasTimer = 0
				// 最初の1回は即移動
				if !g.collides(g.pShape, g.pX-1, g.pY) {
					g.pXPrev = g.pX
					g.pX--
				}
			} else {
				g.dasTimer++
				if g.dasTimer >= 8 { // 8フレームごとに移動
					if !g.collides(g.pShape, g.pX-1, g.pY) {
						g.pXPrev = g.pX
						g.pX--
					}
					g.dasTimer = 6 // 連続移動は少し速め
				}
			}
		} else if right && !left {
			if g.dasDir != 1 {
				g.dasDir = 1
				g.dasTimer = 0
				if !g.collides(g.pShape, g.pX+1, g.pY) {
					g.pXPrev = g.pX
					g.pX++
				}
			} else {
				g.dasTimer++
				if g.dasTimer >= 8 {
					if !g.collides(g.pShape, g.pX+1, g.pY) {
						g.pXPrev = g.pX
						g.pX++
					}
					g.dasTimer = 6
				}
			}
		} else {
			g.dasDir = 0
			g.dasTimer = 0
		}

		// 回転
		if rotate {
			rotated := rotateCW(g.pShape)
			if !g.collides(rotated, g.pX, g.pY) {
				g.pShapePrev = g.pShape
				g.pShape = rotated
			} else {
				// 壁キック: 左右1マスずらして試みる
				if !g.collides(rotated, g.pX+1, g.pY) {
					g.pShapePrev = g.pShape
					g.pShape = rotated
					g.pXPrev = g.pX
					g.pX++
				} else if !g.collides(rotated, g.pX-1, g.pY) {
					g.pShapePrev = g.pShape
					g.pShape = rotated
					g.pXPrev = g.pX
					g.pX--
				}
			}
		}

		// 落下
		g.fallTimer++
		interval := g.fallSpeed
		if down {
			interval = 2 // ソフトドロップ
		}
		if g.fallTimer >= interval {
			g.fallTimer = 0
			g.pYPrev = g.pY
			if !g.collides(g.pShape, g.pX, g.pY+1) {
				g.pY++
			} else {
				g.lockPiece()
			}
		}

	case StateGameOver:
		if rotate || left || right || down {
			// リスタート
			*g = Game{fallSpeed: 30, state: StateReady}
			g.nextType = randomPieceType()
			g.initialized = false
			g.needRedraw = false
		}
	}
}
