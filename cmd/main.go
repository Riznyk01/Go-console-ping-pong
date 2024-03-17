package main

import (
	"fmt"
	"github.com/eiannone/keyboard"
	"github.com/gopxl/beep/speaker"
	"log"
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
	// rocket size
	rocketSide = 4
	// window title
	title = " Go PING-PONG "
	// chars for the filling pixels, empty pixels, header, footer filling
	filler       = '█'
	emptyFiller  = '░'
	headerFiller = '▬'
	footerFiller = '▬'
	// delay between screen updating
	delay = 16 * time.Millisecond
	// screen center for starting the game
	centerRow, centerCol = (windowHeight / 2) + 1, (windowWidth / 2) + 1
	failPause            = 2 * time.Millisecond
	sampleRate           = 44100
	bufferSize           = 1411
)

type Game struct {
	ball        *models.Ball
	leftRocket  *models.Rocket
	rightRocket *models.Rocket
	score       [2]int
	audio       *audio_engine.AudioPlayer
}

func NewGame(ball *models.Ball, aRocket *models.Rocket, bRocket *models.Rocket, soundEngine *audio_engine.AudioPlayer) *Game {
	return &Game{
		ball:        ball,
		leftRocket:  aRocket,
		rightRocket: bRocket,
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
		models.Coordinates{X: centerCol, Y: centerRow}, models.Coordinates{X: 1, Y: 1})

	lRocket := game.NewRocket(models.Coordinates{X: 0, Y: centerRow}, rocketSide)
	rRocket := game.NewRocket(models.Coordinates{X: windowWidth, Y: centerRow}, rocketSide)

	pingpong := NewGame(ball, lRocket, rRocket, soundPlayer)

	keysChannel := make(chan keyboard.KeyEvent)

	go pressingReceiver(keysChannel)
	go pingpong.handleKeyEventsForRockets(keysChannel)

	for {
		pingpong.move()
		pingpong.updateScreen()
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

func (g *Game) handleKeyEventsForRockets(ch chan keyboard.KeyEvent) {
	for {
		select {
		case char := <-ch:
			if char.Rune == 'W' || char.Rune == 'w' {
				if g.leftRocket.Coord.Y-rocketSide > 0 {
					g.leftRocket.Coord.Y -= 1
				}
			} else if char.Rune == 'S' || char.Rune == 's' {
				if g.leftRocket.Coord.Y+rocketSide < windowHeight {
					g.leftRocket.Coord.Y += 1
				}
			}
		case key := <-ch:
			if key.Key == keyboard.KeyArrowUp {
				if g.rightRocket.Coord.Y-rocketSide > 0 {
					g.rightRocket.Coord.Y -= 1
				}
			} else if key.Key == keyboard.KeyArrowDown {
				if g.rightRocket.Coord.Y+rocketSide < windowHeight {
					g.rightRocket.Coord.Y += 1
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
		if g.ball.Coord.Y < (g.rightRocket.Coord.Y-rocketSide) || g.ball.Coord.Y > (g.rightRocket.Coord.Y+rocketSide) {
			g.ball.Coord.X = centerCol
			g.ball.Coord.Y = centerRow
			g.score[0]++
			g.audio.PlaySound("win.wav")
		}
	} else if g.ball.Coord.X == 0 {
		g.ball.Direction.X *= -1
		g.audio.PlaySound("hit.wav")
		if g.ball.Coord.Y < (g.leftRocket.Coord.Y-rocketSide) || g.ball.Coord.Y > (g.leftRocket.Coord.Y+rocketSide) {
			g.ball.Coord.X = centerCol
			g.ball.Coord.Y = centerRow
			g.score[1]++
			g.audio.PlaySound("win.wav")
		}
	}
	if g.ball.Coord.Y == windowHeight {
		g.audio.PlaySound("hit.wav")
		g.ball.Direction.Y *= -1
	} else if g.ball.Coord.Y == 0 {
		g.audio.PlaySound("hit.wav")
		g.ball.Direction.Y *= -1
	}
}

func (g *Game) updateScreen() {
	utils.ClearConsole()
	screen := make([][]rune, windowHeight)
	for i := range screen {
		screen[i] = make([]rune, windowWidth)
	}
	fmt.Print(headerBuilder())
	for y, row := range screen {
		for x := 0; x < windowWidth; x++ {
			if y == g.ball.Coord.Y && x == g.ball.Coord.X {
				row[x] = filler
			} else if x == 0 && (y >= (g.leftRocket.Coord.Y-rocketSide) && y <= (g.leftRocket.Coord.Y+rocketSide)) {
				row[x] = filler
			} else if x == windowWidth-1 && (y >= (g.rightRocket.Coord.Y-rocketSide) && y <= (g.rightRocket.Coord.Y+rocketSide)) {
				row[x] = filler
			} else {
				row[x] = emptyFiller
			}
		}
	}
	for _, line := range screen {
		fmt.Printf("%s\n", string(line))
	}
	fmt.Print(footerBuilder(g.score[0], g.score[1]))
	<-time.After(delay)
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
