package main

import (
	"fmt"
	"github.com/eiannone/keyboard"
	"ping-pong/internal/game"
	"ping-pong/internal/models"
	"ping-pong/internal/utils"
	"strings"
	"time"
)

const (
	//screen size
	cols = 120
	rows = 27
	//rocket size
	rocketSide = 4
	//a chars for the filling pixels, empty pixels, header, footer filling
	filled     = '█'
	empty      = '░'
	headerChar = '▬'
	footerChar = '▬'
	//delay between screen updating
	delay = 25 * time.Millisecond
	//screen center for starting the game
	centerRow, centerCol = (rows / 2) + 1, (cols / 2) + 1
	failPause            = 2 * time.Millisecond
)

type Game struct {
	ball    *models.Ball
	rocketA *models.Rocket
	rocketB *models.Rocket
	scoreA  int
	scoreB  int
}

func NewGame(ball *models.Ball, rocketA *models.Rocket, rocketB *models.Rocket) *Game {
	return &Game{
		ball:    ball,
		rocketA: rocketA,
		rocketB: rocketB,
		scoreA:  0,
		scoreB:  0,
	}
}

func main() {

	ball := game.NewBall(
		models.Coordinates{X: centerCol, Y: centerRow}, models.Coordinates{X: 1, Y: 1})

	rocketA := game.NewRocket(models.Coordinates{X: 0, Y: centerRow}, rocketSide)
	rocketB := game.NewRocket(models.Coordinates{X: cols, Y: centerRow}, rocketSide)

	pingpong := NewGame(ball, rocketA, rocketB)

	go func() {
		if err := keyboard.Open(); err != nil {
			panic(err)
		}
		defer func() {
			_ = keyboard.Close()
		}()
		for {
			char, key, err := keyboard.GetKey()
			if err != nil {
				panic(err)
			}
			if key == keyboard.KeyEsc {
				break
			}
			if char == 'W' || char == 'w' {
				if rocketA.Coord.Y-rocketSide > 0 {
					rocketA.Coord.Y -= 1
				}
			} else if char == 'S' || char == 's' {
				if rocketA.Coord.Y+rocketSide < rows {
					rocketA.Coord.Y += 1
				}
			}

			if key == keyboard.KeyArrowUp {
				if rocketB.Coord.Y-rocketSide > 0 {
					rocketB.Coord.Y -= 1
				}
			} else if key == keyboard.KeyArrowDown {
				if rocketB.Coord.Y+rocketSide < rows {
					rocketB.Coord.Y += 1
				}
			}
		}
	}()

	for {
		pingpong.move()
		pingpong.updateScreen()
	}
}

func (g *Game) move() {
	g.ball.Coord.X += g.ball.Direction.X
	g.ball.Coord.Y += g.ball.Direction.Y

	if g.ball.Coord.X == cols {
		g.ball.Direction.X *= -1
		if g.ball.Coord.Y < (g.rocketB.Coord.Y-rocketSide) || g.ball.Coord.Y > (g.rocketB.Coord.Y+rocketSide) {
			g.ball.Coord.X = centerCol
			g.ball.Coord.Y = centerRow
			g.scoreA++
		}
	} else if g.ball.Coord.X == 0 {
		g.ball.Direction.X *= -1
		if g.ball.Coord.Y < (g.rocketA.Coord.Y-rocketSide) || g.ball.Coord.Y > (g.rocketA.Coord.Y+rocketSide) {
			g.ball.Coord.X = centerCol
			g.ball.Coord.Y = centerRow
			g.scoreB++
		}
	}
	if g.ball.Coord.Y == rows {
		g.ball.Direction.Y *= -1
	} else if g.ball.Coord.Y == 0 {
		g.ball.Direction.Y *= -1
	}
}

func (g *Game) updateScreen() {
	utils.ClearConsole()
	screen := make([][]rune, rows)
	for i := range screen {
		screen[i] = make([]rune, cols)
	}
	fmt.Print(headerBuilder())
	for y, r := range screen {
		for x := 0; x < cols; x++ {
			if y == g.ball.Coord.Y && x == g.ball.Coord.X {
				r[x] = filled
			} else if x == 0 && (y >= (g.rocketA.Coord.Y-rocketSide) && y <= (g.rocketA.Coord.Y+rocketSide)) {
				r[x] = filled
			} else if x == cols-1 && (y >= (g.rocketB.Coord.Y-rocketSide) && y <= (g.rocketB.Coord.Y+rocketSide)) {
				r[x] = filled
			} else {
				r[x] = empty
			}
		}
	}
	for _, line := range screen {
		fmt.Printf("%s\n", string(line))
	}
	fmt.Print(footerBuilder(g.scoreA, g.scoreB))
	<-time.After(delay)
}

func headerBuilder() string {
	title := " Go PING-PONG "
	filler := strings.Repeat(string(headerChar), (cols-len(title))/2)
	return filler + title + filler
}
func footerBuilder(scoreLeft, scoreRight int) string {
	a := fmt.Sprintf("%c SCORE A: %v ", footerChar, scoreLeft)
	b := fmt.Sprintf(" SCORE B: %v %c", scoreRight, footerChar)
	filler := strings.Repeat("▬", cols-len(a)-len(b))
	return a + filler + b
}
