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

const (
	DIRECTION_LEFT  = 0
	DIRECTION_RIGHT = 1
	DIRECTION_UP    = 2
	DIRECTION_DOWN  = 3
)

type Player struct {
	X int
	Y int
}

type Bullet struct {
	X            int
	Y            int
	Direction    int
	ShouldRemove bool
}

type Target struct {
	X int
	Y int
}

type Game struct {
	Rng       *rand.Rand
	IsOver    bool
	Frame     int32
	Player    *Player
	Bullets   []*Bullet
	Targets   []*Target
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

func RandRange(rng *rand.Rand, min int, maxInclusive int) int {
	return rng.Intn(maxInclusive-min) + min
}

func NewGame() *Game {
	timestamp := time.Now().Unix()
	rng := rand.New(rand.NewSource(timestamp))

	g := &Game{}
	g.Rng = rng
	g.IsOver = false
	g.Frame = 0
	g.Player = &Player{
		X: 18,
		Y: Y_MAX,
	}
	g.Bullets = make([]*Bullet, 0, 10)
	g.Targets = make([]*Target, 0, 100)
	return g
}

func (g *Game) Update(commands []string) {
	if g.IsOver {
		return
	}

	g.handleCommands(commands)

	for _, bullet := range g.Bullets {
		switch bullet.Direction {
		case DIRECTION_LEFT:
			if bullet.X == X_MIN {
				bullet.Direction = DIRECTION_RIGHT
			} else {
				bullet.X -= 1
			}
		case DIRECTION_RIGHT:
			if bullet.X == X_MAX {
				bullet.Direction = DIRECTION_LEFT
			} else {
				bullet.X += 1
			}
		case DIRECTION_UP:
			if bullet.Y == Y_MIN {
				bullet.Direction = DIRECTION_DOWN
			} else {
				bullet.Y -= 1
			}
		case DIRECTION_DOWN:
			if bullet.Y == Y_MAX {
				bullet.ShouldRemove = true
			} else {
				bullet.Y += 1
			}
		}
	}

	for _, bullet := range g.Bullets {
		if g.Player.X <= bullet.X && bullet.X <= g.Player.X+2 && g.Player.Y == bullet.Y {
			g.IsOver = true
		}
	}

	removedBullets := make([]*Bullet, 0, cap(g.Bullets))
	for _, bullet := range g.Bullets {
		if !bullet.ShouldRemove {
			removedBullets = append(removedBullets, bullet)
		}
	}
	g.Bullets = removedBullets

	if g.Frame%60 == 0 {
		g.spawnTarget()
	}

	g.Frame += 1
}

func (g *Game) handleCommands(commands []string) {
	for _, command := range commands {
		switch command {
		case "left":
			g.Player.X -= 1
			g.Player.X = Clamp(X_MIN, g.Player.X, X_MAX)
		case "right":
			g.Player.X += 1
			g.Player.X = Clamp(X_MIN, g.Player.X, X_MAX-2)
		case "shoot":
			g.shoot()
		}
	}

}

func (g *Game) spawnTarget() {
	target := &Target{
		X: RandRange(g.Rng, X_MIN+1, X_MAX-1),
		Y: RandRange(g.Rng, Y_MIN, 15),
	}
	g.Targets = append(g.Targets, target)
}

func (g *Game) shoot() {
	bullet := &Bullet{
		X:         g.Player.X + 1,
		Y:         g.Player.Y - 1,
		Direction: DIRECTION_UP,
	}
	g.Bullets = append(g.Bullets, bullet)
}
