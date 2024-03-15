package main

import (
	"fmt"
	"github.com/eiannone/keyboard"
	"os"
	"os/exec"
	"strings"
	"time"
)

const (
	cols       = 120
	rows       = 27
	rocketSide = 4
	filled     = '█'
	empty      = '░'
	headerChar = '▬'
	footerChar = '▬'
	failPause  = 2 * time.Millisecond
)

type coord struct {
	x int
	y int
}

var scoreA, scoreB = 0, 0
var delay = 25 * time.Millisecond
var ballDirection = &coord{x: 1, y: 1}
var ballCoord = &coord{
	x: centerCol, y: centerRow,
}
var RocketACoord = &coord{
	x: 0, y: centerRow,
}
var centerRow, centerCol = (rows / 2) + 1, (cols / 2) + 1

var RocketBCoord = &coord{
	x: cols,
	y: centerRow,
}

func main() {

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
				if RocketACoord.y-rocketSide > 0 {
					RocketACoord.y -= 1
				}
			} else if char == 'S' || char == 's' {
				if RocketACoord.y+rocketSide < rows {
					RocketACoord.y += 1
				}
			}

			if key == keyboard.KeyArrowUp {
				if RocketBCoord.y-rocketSide > 0 {
					RocketBCoord.y -= 1
				}
			} else if key == keyboard.KeyArrowDown {
				if RocketBCoord.y+rocketSide < rows {
					RocketBCoord.y += 1
				}
			}
		}
	}()

	for {
		ballMovement()
		printScreen(ballCoord)
	}
}

func ballMovement() {
	ballCoord.x += ballDirection.x
	ballCoord.y += ballDirection.y

	if ballCoord.x == cols {
		ballDirection.x *= -1
		if ballCoord.y < (RocketBCoord.y-rocketSide) || ballCoord.y > (RocketBCoord.y+rocketSide) {
			ballCoord.x = centerCol
			ballCoord.y = centerRow
			scoreA++
		}
	} else if ballCoord.x == 0 {
		ballDirection.x *= -1
		if ballCoord.y < (RocketACoord.y-rocketSide) || ballCoord.y > (RocketACoord.y+rocketSide) {
			ballCoord.x = centerCol
			ballCoord.y = centerRow
			scoreB++
		}
	}

	if ballCoord.y == rows {
		ballDirection.y *= -1
	} else if ballCoord.y == 0 {
		ballDirection.y *= -1
	}
}
func printScreen(ball *coord) {
	clearConsole()
	screen := make([][]rune, rows)
	for i := range screen {
		screen[i] = make([]rune, cols)
	}
	fmt.Print(headerBuilder())
	for y, r := range screen {
		for x := 0; x < cols; x++ {
			if y == ball.y && x == ball.x {
				r[x] = filled
			} else if x == 0 && (y >= (RocketACoord.y-rocketSide) && y <= (RocketACoord.y+rocketSide)) {
				r[x] = filled
			} else if x == cols-1 && (y >= (RocketBCoord.y-rocketSide) && y <= (RocketBCoord.y+rocketSide)) {
				r[x] = filled
			} else {
				r[x] = empty
			}
		}
	}
	for _, line := range screen {
		fmt.Printf("%s\n", string(line))
	}
	fmt.Print(footerBuilder())
	<-time.After(delay)
}
func clearConsole() {
	cmd := exec.Command("cmd", "/c", "cls")
	cmd.Stdout = os.Stdout
	cmd.Run()
}
func headerBuilder() string {
	title := " Go PING-PONG "
	filler := strings.Repeat(string(headerChar), (cols-len(title))/2)
	return filler + title + filler
}
func footerBuilder() string {
	a := fmt.Sprintf("%c SCORE A: %v ", footerChar, scoreA)
	b := fmt.Sprintf(" SCORE B: %v %c", scoreB, footerChar)
	filler := strings.Repeat("▬", cols-len(a)-len(b))
	return a + filler + b
}
