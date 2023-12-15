package model

import (
	"math/rand"
	"time"
)

type Game struct {
	Rng    *rand.Rand
	IsOver bool
	Frame  int32
}

func NewGame() *Game {
	timestamp := time.Now().Unix()
	rng := rand.New(rand.NewSource(timestamp))

	g := &Game{}
	g.Rng = rng
	g.IsOver = false
	g.Frame = 0
	return g
}

func (g *Game) Update(command string) {
	if g.IsOver {
		return
	}

	switch command {
	case "left":
	case "right":
	case "shoot":
	}

	g.Frame += 1
}
