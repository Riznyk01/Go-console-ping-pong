package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

const (
	cols = 119
	rows = 27
)

var filled = '█'
var empty = ' '

type coord struct {
	x int
	y int
}

var scoreA, scoreB = 0, 0
var delay = 100 * time.Millisecond

var ballDirection = &coord{x: 1, y: 1}

var ballCoord = &coord{
	x: centerCol,
	y: centerRow,
}

var RocketACoord = &coord{
	x: 0,
	y: centerRow,
}

var RocketSide = 4
var centerRow = (rows / 2) + 1
var centerCol = (cols / 2) + 1

var RocketBCoord = &coord{
	x: cols,
	y: centerRow,
}

func main() {
	// ball movement
	for {
		ballCoord.x += ballDirection.x
		ballCoord.y += ballDirection.y
		if ballCoord.x == cols || ballCoord.x == 0 {
			ballDirection.x *= -1
		}
		if ballCoord.y == rows || ballCoord.y == 0 {
			ballDirection.y *= -1
		}

		printScreen(ballCoord)
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
			} else if x == 0 && (y >= (RocketACoord.y-RocketSide) && y <= (RocketACoord.y+RocketSide)) {
				r[x] = filled
			} else if x == cols-1 && (y >= (RocketBCoord.y-RocketSide) && y <= (RocketBCoord.y+RocketSide)) {
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
	title := " PING-PONG "
	filler := strings.Repeat("▒", (cols-len(title))/2)
	return filler + title + filler
}
func footerBuilder() string {
	footerText := fmt.Sprintf("▒ SCORE A:%d ▒ SCORE B:%d ", scoreA, scoreB)
	filler := strings.Repeat("▒", cols-len(footerText))
	return footerText + filler
}
