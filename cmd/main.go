package main

import (
	"fmt"
	"github.com/eiannone/keyboard"
	"github.com/gopxl/beep/speaker"
	"log"
	"math/rand"
	"os"
	"ping-pong/internal/audio_engine"
	"ping-pong/internal/game"
	"ping-pong/internal/models"
	"ping-pong/internal/utils"
	"strings"
	"time"
)

const (
	// screen size
	windowWidth  = 120
	windowHeight = 27
	// window title
	title = " Go PING-PONG "
	// chars for the filling pixels, empty pixels, header, footer filling
	filler            = '▒'
	emptyFiller       = 'ₓ'
	headerFiller      = '▬'
	footerFiller      = '▬'
	racketsLineFiller = '⁞'
	// delays between screen updating and ball moving
	ballMovingDelay = 40 * time.Millisecond
	renderDelay     = 40 * time.Millisecond
	// screen center for starting the game
	centerRow, centerCol = (windowHeight / 2) + 1, (windowWidth / 2) + 1
	failPause            = 2 * time.Millisecond
	sampleRate           = 44100
	bufferSize           = 1411
	rocketUpCommand      = "up"
	rocketDownCommand    = "down"
)

var ballsFillers = [3]rune{'▒', '▓', '█'}

type Game struct {
	ball        *models.Ball
	leftRacket  *models.Racket
	rightRacket *models.Racket
	score       [2]int
	audio       *audio_engine.AudioPlayer
}

func NewGame(ball *models.Ball, aRacket *models.Racket, bRacket *models.Racket, soundEngine *audio_engine.AudioPlayer) *Game {
	return &Game{
		ball:        ball,
		leftRacket:  aRacket,
		rightRacket: bRacket,
		score:       [2]int{0, 0},
		audio:       soundEngine,
	}
}

func main() {
	logFile, err := os.OpenFile("logfile.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("error occurred while opening logfile:", err)
	}
	defer logFile.Close()
	log.SetOutput(logFile)

	speaker.Init(sampleRate, bufferSize)

	soundPlayer := audio_engine.NewAudioPlayer()
	soundPlayer.LoadSound("sound")

	ball := game.NewBall(
		models.Coordinates{X: centerCol, Y: centerRow}, models.Coordinates{X: centerCol, Y: centerRow}, models.Coordinates{X: 1, Y: 1})

	lRacket := game.NewRacket(models.Coordinates{X: 0, Y: centerRow}, 4)
	rRacket := game.NewRacket(models.Coordinates{X: windowWidth, Y: centerRow}, 4)

	pingpong := NewGame(ball, lRacket, rRacket, soundPlayer)

	keysChannel := make(chan keyboard.KeyEvent)

	go pressingReceiver(keysChannel)
	go pingpong.handleKeyEventsForRackets(keysChannel)

	screen := make([][]rune, windowHeight)
	for i := 0; i < windowHeight; i++ {
		lineRow := make([]rune, windowWidth)
		lineRow[0] = racketsLineFiller
		lineRow[len(lineRow)-1] = racketsLineFiller
		for linePx := 1; linePx < len(lineRow)-1; linePx++ {
			lineRow[linePx] = emptyFiller
		}
		screen[i] = lineRow
	}
	go pingpong.move()
	go pingpong.updateScreen(screen)
	for {
	}
}

func pressingReceiver(ch chan keyboard.KeyEvent) {
	keysEvents, err := keyboard.GetKeys(10)
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = keyboard.Close()
	}()

	for {
		event := <-keysEvents
		if event.Err != nil {
			panic(event.Err)
		}
		if event.Key == keyboard.KeyEsc {
			break
		}
		ch <- event
	}
}

func (g *Game) handleKeyEventsForRackets(ch chan keyboard.KeyEvent) {
	for {
		select {
		case event := <-ch:
			if event.Rune == 'W' || event.Rune == 'w' {
				g.racketMove(rocketUpCommand, g.leftRacket)
			} else if event.Rune == 'S' || event.Rune == 's' {
				g.racketMove(rocketDownCommand, g.leftRacket)
			} else if event.Key == keyboard.KeyArrowUp {
				g.racketMove(rocketUpCommand, g.rightRacket)
			} else if event.Key == keyboard.KeyArrowDown {
				g.racketMove(rocketDownCommand, g.rightRacket)
			}
		}
	}
}
func (g *Game) racketMove(event string, racket *models.Racket) {
	if event == rocketUpCommand {
		if racket.Coord.Y-racket.Side > 1 {
			racket.Coord.Y -= 1
			if racket.Coord.Y-racket.Side == 1 {
				g.audio.PlaySound("wood_hit.wav")
			}
		}
	} else if event == rocketDownCommand {
		if racket.Coord.Y+racket.Side < windowHeight {
			racket.Coord.Y += 1
			if racket.Coord.Y+racket.Side == windowHeight {
				g.audio.PlaySound("wood_hit.wav")
			}
		}
	}
}

func (g *Game) move() {
	for {
		<-time.After(ballMovingDelay)
		g.ball.Coord.X += g.ball.Direction.X
		g.ball.Coord.Y += g.ball.Direction.Y

		if g.ball.Coord.X == windowWidth {
			g.ball.Direction.X = -g.ball.Direction.X
			g.audio.PlaySound("hit.wav")
			if g.ball.Coord.Y < (g.rightRacket.Coord.Y-g.rightRacket.Side) || g.ball.Coord.Y > (g.rightRacket.Coord.Y+g.rightRacket.Side) {
				g.ball.Coord.X, g.ball.Coord.Y = centerCol, centerRow
				g.score[0]++
				g.audio.PlaySound("win.wav")
			}
		} else if g.ball.Coord.X == 1 {
			g.ball.Direction.X = -g.ball.Direction.X
			g.audio.PlaySound("hit.wav")
			if g.ball.Coord.Y < (g.leftRacket.Coord.Y-g.leftRacket.Side) || g.ball.Coord.Y > (g.leftRacket.Coord.Y+g.leftRacket.Side) {
				g.ball.Coord.X, g.ball.Coord.Y = centerCol, centerRow
				g.score[1]++
				g.audio.PlaySound("win.wav")
			}
		}
		if g.ball.Coord.Y == windowHeight {
			g.audio.PlaySound("hit.wav")
			g.ball.Direction.Y = -g.ball.Direction.Y
		} else if g.ball.Coord.Y == 1 {
			g.audio.PlaySound("hit.wav")
			g.ball.Direction.Y = -g.ball.Direction.Y
		}
	}
}

// updateScreen update the screen with the given rune array
func (g *Game) updateScreen(screen [][]rune) {
	for {
		utils.ClearConsole()
		fmt.Print(headerBuilder())
		// Clear rackets' position lines
		for _, line := range screen {
			if line[0] == filler {
				line[0] = racketsLineFiller
			}
			if line[len(line)-1] == filler {
				line[len(line)-1] = racketsLineFiller
			}

		}
		// Clear the balls' last position, fill the balls' current position,
		// and save the balls' position to the last coordinates
		log.Printf("a ball coordinates are: X %d, Y %d", g.ball.Coord.X, g.ball.Coord.Y)
		screen[g.ball.LastCoord.Y-1][g.ball.LastCoord.X-1] = emptyFiller
		g.ball.LastCoord.X, g.ball.LastCoord.Y = g.ball.Coord.X, g.ball.Coord.Y
		screen[g.ball.Coord.Y-1][g.ball.Coord.X-1] = ballsFillers[rand.Intn(len(ballsFillers))]
		// Fill rackets' positions
		for l := g.leftRacket.Coord.Y - g.leftRacket.Side; l <= g.leftRacket.Coord.Y+g.leftRacket.Side; l++ {
			screen[l-1][0] = filler
		}
		for r := g.rightRacket.Coord.Y - g.rightRacket.Side; r <= g.rightRacket.Coord.Y+g.rightRacket.Side; r++ {
			screen[r-1][len(screen[r-1])-1] = filler
		}
		log.Printf("left racket coordinates are: X %d, Y %d", g.leftRacket.Coord.X, g.leftRacket.Coord.Y)
		log.Printf("right racket coordinates are: X %d, Y %d", g.leftRacket.Coord.X, g.leftRacket.Coord.Y)
		// Print characters in each line of the screen
		for _, line := range screen {
			var pxRow string
			for _, px := range line {
				pxRow += string(px)
			}
			fmt.Printf("%s\n", pxRow)
		}
		fmt.Print(footerBuilder(g.score[0], g.score[1]))
		<-time.After(renderDelay)
	}
}

func headerBuilder() string {
	topFiller := strings.Repeat(string(headerFiller), (windowWidth-len(title))/2)
	return topFiller + title + topFiller
}
func footerBuilder(scoreLeft, scoreRight int) string {
	left, right := fmt.Sprintf("%c SCORE: %v ", footerFiller, scoreLeft), fmt.Sprintf(" SCORE: %v %c", scoreRight, footerFiller)
	bottomFiller := strings.Repeat(string(footerFiller), windowWidth-len(left)-len(right)+4)
	return left + bottomFiller + right
}
