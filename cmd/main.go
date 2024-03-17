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
	// racket size
	racketsSide = 4
	// window title
	title = " Go PING-PONG "
	// chars for the filling pixels, empty pixels, header, footer filling
	filler       = '█'
	emptyFiller  = '░'
	headerFiller = '▬'
	footerFiller = '▬'
	// delays between screen updating and ball moving
	//renderDelay     = 150 * time.Millisecond
	ballMovingDelay = 15 * time.Millisecond
	// screen center for starting the game
	centerRow, centerCol = (windowHeight / 2) + 1, (windowWidth / 2) + 1
	failPause            = 2 * time.Millisecond
	sampleRate           = 44100
	bufferSize           = 1411
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

	lRacket := game.NewRacket(models.Coordinates{X: 0, Y: centerRow}, racketsSide)
	rRacket := game.NewRacket(models.Coordinates{X: windowWidth, Y: centerRow}, racketsSide)

	pingpong := NewGame(ball, lRacket, rRacket, soundPlayer)

	keysChannel := make(chan keyboard.KeyEvent)

	go pressingReceiver(keysChannel)
	go pingpong.handleKeyEventsForRackets(keysChannel)

	screen := [windowHeight][windowWidth]rune{}
	for y := 0; y < windowHeight; y++ {
		for x := 0; x < windowWidth; x++ {
			screen[y][x] = emptyFiller
		}
	}
	for {
		pingpong.move()
		pingpong.updateScreen(screen)
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
		case char := <-ch:
			if char.Rune == 'W' || char.Rune == 'w' {
				if g.leftRacket.Coord.Y-racketsSide > 1 {
					g.leftRacket.Coord.Y -= 1
				}
			} else if char.Rune == 'S' || char.Rune == 's' {
				if g.leftRacket.Coord.Y+racketsSide < windowHeight {
					g.leftRacket.Coord.Y += 1
				}
			}
		case key := <-ch:
			if key.Key == keyboard.KeyArrowUp {
				if g.rightRacket.Coord.Y-racketsSide > 1 {
					g.rightRacket.Coord.Y -= 1
				}
			} else if key.Key == keyboard.KeyArrowDown {
				if g.rightRacket.Coord.Y+racketsSide < windowHeight {
					g.rightRacket.Coord.Y += 1
				}
			}
		}
	}
}

func (g *Game) move() {
	g.ball.Coord.X += g.ball.Direction.X
	g.ball.Coord.Y += g.ball.Direction.Y

	if g.ball.Coord.X == windowWidth {
		g.ball.Direction.X *= -1
		g.audio.PlaySound("hit.wav")
		if g.ball.Coord.Y < (g.rightRacket.Coord.Y-racketsSide) || g.ball.Coord.Y > (g.rightRacket.Coord.Y+racketsSide) {
			g.ball.Coord.X = centerCol
			g.ball.Coord.Y = centerRow
			g.score[0]++
			g.audio.PlaySound("win.wav")
		}
	} else if g.ball.Coord.X == 1 {
		g.ball.Direction.X *= -1
		g.audio.PlaySound("hit.wav")
		if g.ball.Coord.Y < (g.leftRacket.Coord.Y-racketsSide) || g.ball.Coord.Y > (g.leftRacket.Coord.Y+racketsSide) {
			g.ball.Coord.X = centerCol
			g.ball.Coord.Y = centerRow
			g.score[1]++
			g.audio.PlaySound("win.wav")
		}
	}
	if g.ball.Coord.Y == windowHeight {
		g.audio.PlaySound("hit.wav")
		g.ball.Direction.Y *= -1
	} else if g.ball.Coord.Y == 1 {
		g.audio.PlaySound("hit.wav")
		g.ball.Direction.Y *= -1
	}
	<-time.After(ballMovingDelay)
}

// updateScreen update the screen with the given rune array
func (g *Game) updateScreen(screen [windowHeight][windowWidth]rune) {
	utils.ClearConsole()
	fmt.Print(headerBuilder())
	// Clear rackets' position lines
	for _, line := range screen {
		line[0] = emptyFiller
		line[len(line)-1] = emptyFiller
	}
	// Clear the ball's last position, fill the ball's current position,
	// and save the ball's position to the last coordinates
	log.Printf("a ball coordinates are: X %d, Y %d", g.ball.Coord.X, g.ball.Coord.Y)
	screen[g.ball.LastCoord.Y-1][g.ball.LastCoord.X-1] = emptyFiller
	g.ball.LastCoord.X = g.ball.Coord.X
	g.ball.LastCoord.Y = g.ball.Coord.Y
	screen[g.ball.Coord.Y-1][g.ball.Coord.X-1] = ballsFillers[rand.Intn(len(ballsFillers))]
	// Fill rackets' positions
	for l := g.leftRacket.Coord.Y - racketsSide; l <= g.leftRacket.Coord.Y+racketsSide; l++ {
		screen[l-1][0] = filler
	}
	for r := g.rightRacket.Coord.Y - racketsSide; r <= g.rightRacket.Coord.Y+racketsSide; r++ {
		screen[r-1][len(screen[r-1])-1] = filler
	}
	log.Printf("left racket coordinates are: X %d, Y %d", g.leftRacket.Coord.X, g.leftRacket.Coord.Y)
	log.Printf("right racket coordinates are: X %d, Y %d", g.leftRacket.Coord.X, g.leftRacket.Coord.Y)
	// Log the characters of the screens' array
	for _, line := range screen {
		log.Printf("%c\n", line)
	}
	// Print characters in each line of the screen
	for _, line := range screen {
		var pxRow string
		for _, px := range line {
			pxRow += string(px)
		}
		fmt.Printf("%s\n", pxRow)
	}
	fmt.Print(footerBuilder(g.score[0], g.score[1]))
	//<-time.After(renderDelay)
}

func headerBuilder() string {
	topFiller := strings.Repeat(string(headerFiller), (windowWidth-len(title))/2)
	return topFiller + title + topFiller
}
func footerBuilder(scoreLeft, scoreRight int) string {
	a := fmt.Sprintf("%c SCORE A: %v ", footerFiller, scoreLeft)
	b := fmt.Sprintf(" SCORE B: %v %c", scoreRight, footerFiller)
	bottomFiller := strings.Repeat(string(footerFiller), windowWidth-len(a)-len(b)+4)
	return a + bottomFiller + b
}
