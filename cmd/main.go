package main

import (
	"fmt"
	"github.com/eiannone/keyboard"
	"github.com/gopxl/beep/speaker"
	"ping-pong/internal/audio_engine"
	"ping-pong/internal/game"
	"ping-pong/internal/models"
	"ping-pong/internal/utils"
	"strings"
	"time"
)

const (
	// screen size
	cols = 120
	rows = 27
	// rocket size
	rocketSide = 4
	// window title
	title = " Go PING-PONG "
	// chars for the filling pixels, empty pixels, header, footer filling
	filled     = '█'
	empty      = '░'
	headerChar = '▬'
	footerChar = '▬'
	// delay between screen updating
	delay = 16 * time.Millisecond
	// screen center for starting the game
	centerRow, centerCol = (rows / 2) + 1, (cols / 2) + 1
	failPause            = 2 * time.Millisecond
	sampleRate           = 44100
	bufferSize           = 1411
)

type Game struct {
	ball    *models.Ball
	aRocket *models.Rocket
	bRocket *models.Rocket
	aScore  int
	bScore  int
	sound   *audio_engine.AudioPlayer
}

func NewGame(ball *models.Ball, aRocket *models.Rocket, bRocket *models.Rocket, sound *audio_engine.AudioPlayer) *Game {
	return &Game{
		ball:    ball,
		aRocket: aRocket,
		bRocket: bRocket,
		aScore:  0,
		bScore:  0,
		sound:   sound,
	}
}

func main() {
	speaker.Init(sampleRate, bufferSize)

	player := audio_engine.NewAudioPlayer()
	player.LoadSound("sound")

	ball := game.NewBall(
		models.Coordinates{X: centerCol, Y: centerRow}, models.Coordinates{X: 1, Y: 1})

	rocketA := game.NewRocket(models.Coordinates{X: 0, Y: centerRow}, rocketSide)
	rocketB := game.NewRocket(models.Coordinates{X: cols, Y: centerRow}, rocketSide)

	pingpong := NewGame(ball, rocketA, rocketB, player)

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
				if g.aRocket.Coord.Y-rocketSide > 0 {
					g.aRocket.Coord.Y -= 1
				}
			} else if char.Rune == 'S' || char.Rune == 's' {
				if g.aRocket.Coord.Y+rocketSide < rows {
					g.aRocket.Coord.Y += 1
				}
			}
		case key := <-ch:
			if key.Key == keyboard.KeyArrowUp {
				if g.bRocket.Coord.Y-rocketSide > 0 {
					g.bRocket.Coord.Y -= 1
				}
			} else if key.Key == keyboard.KeyArrowDown {
				if g.bRocket.Coord.Y+rocketSide < rows {
					g.bRocket.Coord.Y += 1
				}
			}
		}
	}
}

func (g *Game) move() {
	g.ball.Coord.X += g.ball.Direction.X
	g.ball.Coord.Y += g.ball.Direction.Y

	if g.ball.Coord.X == cols {
		g.ball.Direction.X *= -1
		g.sound.PlaySound("hit.wav")
		if g.ball.Coord.Y < (g.bRocket.Coord.Y-rocketSide) || g.ball.Coord.Y > (g.bRocket.Coord.Y+rocketSide) {
			g.ball.Coord.X = centerCol
			g.ball.Coord.Y = centerRow
			g.aScore++
			g.sound.PlaySound("win.wav")
		}
	} else if g.ball.Coord.X == 0 {
		g.ball.Direction.X *= -1
		g.sound.PlaySound("hit.wav")
		if g.ball.Coord.Y < (g.aRocket.Coord.Y-rocketSide) || g.ball.Coord.Y > (g.aRocket.Coord.Y+rocketSide) {
			g.ball.Coord.X = centerCol
			g.ball.Coord.Y = centerRow
			g.bScore++
			g.sound.PlaySound("win.wav")
		}
	}
	if g.ball.Coord.Y == rows {
		g.sound.PlaySound("hit.wav")
		g.ball.Direction.Y *= -1
	} else if g.ball.Coord.Y == 0 {
		g.sound.PlaySound("hit.wav")
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
	for y, row := range screen {
		for x := 0; x < cols; x++ {
			if y == g.ball.Coord.Y && x == g.ball.Coord.X {
				row[x] = filled
			} else if x == 0 && (y >= (g.aRocket.Coord.Y-rocketSide) && y <= (g.aRocket.Coord.Y+rocketSide)) {
				row[x] = filled
			} else if x == cols-1 && (y >= (g.bRocket.Coord.Y-rocketSide) && y <= (g.bRocket.Coord.Y+rocketSide)) {
				row[x] = filled
			} else {
				row[x] = empty
			}
		}
	}
	for _, line := range screen {
		fmt.Printf("%s\n", string(line))
	}
	fmt.Print(footerBuilder(g.aScore, g.bScore))
	<-time.After(delay)
}

func headerBuilder() string {
	filler := strings.Repeat(string(headerChar), (cols-len(title))/2)
	return filler + title + filler
}
func footerBuilder(scoreLeft, scoreRight int) string {
	a := fmt.Sprintf("%c SCORE A: %v ", footerChar, scoreLeft)
	b := fmt.Sprintf(" SCORE B: %v %c", scoreRight, footerChar)
	filler := strings.Repeat("▬", cols-len(a)-len(b))
	return a + filler + b
}
