package model

import (
	"math/rand"
	"time"
)

const (
	X_MIN = 2
	X_MAX = 640/16 - 3
	Y_MIN = 2
	Y_MAX = 400/16 - 3
)

type Player struct {
	X int
	Y int
}

type Game struct {
	Rng       *rand.Rand
	IsOver    bool
	Frame     int32
	Player    Player
	Score     int
	HighScore int
}

func Clamp(min int, val int, max int) int {
	if val < min {
		return min
	}
	if val > max {
		return max
	}
	return val
}

func NewGame() *Game {
	timestamp := time.Now().Unix()
	rng := rand.New(rand.NewSource(timestamp))

	g := &Game{}
	g.Rng = rng
	g.IsOver = false
	g.Frame = 0
	g.Player.X = 18
	g.Player.Y = Y_MAX
	return g
}

func (g *Game) Update(command string) {
	if g.IsOver {
		return
	}

	switch command {
	case "left":
		g.Player.X -= 1
		g.Player.X = Clamp(X_MIN, g.Player.X, X_MAX)
	case "right":
		g.Player.X += 1
		g.Player.X = Clamp(X_MIN, g.Player.X, X_MAX-2)
	case "shoot":
	}

	g.Frame += 1
}
