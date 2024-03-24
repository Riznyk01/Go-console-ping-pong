package game

import (
	"fmt"
	"github.com/eiannone/keyboard"
	"log"
	"math/rand"
	"ping-pong/internal/audio_engine"
	ping_pong "ping-pong/internal/config"
	"ping-pong/internal/utils"
	"strings"
	"time"
)

var ballsFillers = [3]rune{'▒', '▓', '█'}

type Game struct {
	Ball        *Ball
	LeftRacket  *Racket
	RightRacket *Racket
	Field       *Field
	Score       [2]int
	Audio       *audio_engine.AudioPlayer
	Cfg         *ping_pong.Config
	KeysChannel chan keyboard.KeyEvent
}

func NewGame(ball *Ball, aRacket *Racket, bRacket *Racket, field *Field, soundEngine *audio_engine.AudioPlayer, cfg *ping_pong.Config, keysChan chan keyboard.KeyEvent) *Game {
	return &Game{
		Ball:        ball,
		LeftRacket:  aRacket,
		RightRacket: bRacket,
		Field:       field,
		Score:       [2]int{0, 0},
		Audio:       soundEngine,
		Cfg:         cfg,
		KeysChannel: keysChan,
	}
}

func (g *Game) Play() {
	go g.GamePad()
	for {
		g.Ball.Move()
		g.CheckGoal()
		g.Output()
	}
}

func (g *Game) CheckGoal() {
	if g.Ball.coord.X == g.Cfg.WindowWidth {
		if g.Ball.coord.Y < (g.RightRacket.coord.Y-g.RightRacket.side) || g.Ball.coord.Y > (g.RightRacket.coord.Y+g.RightRacket.side) {
			g.Ball.coord.X, g.Ball.coord.Y = g.Cfg.CenterCol, g.Cfg.CenterRow
			g.Score[0]++
			g.Audio.Play("win.wav")
		}
	} else if g.Ball.coord.X == 1 {
		if g.Ball.coord.Y < (g.LeftRacket.coord.Y-g.LeftRacket.side) || g.Ball.coord.Y > (g.LeftRacket.coord.Y+g.LeftRacket.side) {
			g.Ball.coord.X, g.Ball.coord.Y = g.Cfg.CenterCol, g.Cfg.CenterRow
			g.Score[1]++
			g.Audio.Play("win.wav")
		}
	}
}

// Output update the screen with the given rune array
func (g *Game) Output() {
	utils.ClearConsole()
	fmt.Print(g.HeaderBuilder())
	// Clear rackets' position lines
	for _, line := range g.Field.Console {
		if line[0] != g.Cfg.RacketsLineFiller {
			line[0] = g.Cfg.RacketsLineFiller
		}
		if line[len(line)-1] != g.Cfg.RacketsLineFiller {
			line[len(line)-1] = g.Cfg.RacketsLineFiller
		}

	}
	// Clear the ball`s last position, fill the ball`s current position,
	// and save the ball`s position to the last coordinates
	log.Printf("a ball coordinates are: X %d, Y %d", g.Ball.coord.X, g.Ball.coord.Y)
	g.Field.Console[g.Ball.lastCoord.Y-1][g.Ball.lastCoord.X-1] = g.Cfg.EmptyFiller
	g.Ball.lastCoord.X, g.Ball.lastCoord.Y = g.Ball.coord.X, g.Ball.coord.Y
	g.Field.Console[g.Ball.coord.Y-1][g.Ball.coord.X-1] = ballsFillers[rand.Intn(len(ballsFillers))]

	// Fill rackets` positions
	for l := g.LeftRacket.coord.Y - g.LeftRacket.side; l <= g.LeftRacket.coord.Y+g.LeftRacket.side; l++ {
		g.Field.Console[l-1][0] = g.Cfg.Filler
	}
	for r := g.RightRacket.coord.Y - g.RightRacket.side; r <= g.RightRacket.coord.Y+g.RightRacket.side; r++ {
		g.Field.Console[r-1][len(g.Field.Console[r-1])-1] = g.Cfg.Filler
	}
	log.Printf("left racket coordinates are: X %d, Y %d", g.LeftRacket.coord.X, g.LeftRacket.coord.Y)
	log.Printf("right racket coordinates are: X %d, Y %d", g.LeftRacket.coord.X, g.LeftRacket.coord.Y)
	// Print characters in each line of the screen
	for _, line := range g.Field.Console {
		var pxRow string
		for _, px := range line {
			pxRow += string(px)
		}
		fmt.Printf("%s\n", pxRow)
	}
	fmt.Print(g.FooterBuilder(g.Score[0], g.Score[1]))
	<-time.After(g.Cfg.RenderDelay)
}

func (g *Game) HeaderBuilder() string {
	topFiller := strings.Repeat(string(g.Cfg.HeaderFiller), (g.Cfg.WindowWidth-len(g.Cfg.Title))/2)
	return topFiller + g.Cfg.Title + topFiller
}
func (g *Game) FooterBuilder(scoreLeft, scoreRight int) string {
	left, right := fmt.Sprintf("%c PLAYER 1: %v ", g.Cfg.FooterFiller, scoreLeft), fmt.Sprintf(" PLAYER 2: %v %c", scoreRight, g.Cfg.FooterFiller)
	bottomFiller := strings.Repeat(string(g.Cfg.FooterFiller), g.Cfg.WindowWidth-len(left)-len(right)+4)
	return left + bottomFiller + right
}

func (g *Game) GamePad() {
	for {
		select {
		case event := <-g.KeysChannel:
			if event.Rune == 'W' || event.Rune == 'w' {
				g.LeftRacket.Move(g.Cfg.RocketUpCommand)
			} else if event.Rune == 'S' || event.Rune == 's' {
				g.LeftRacket.Move(g.Cfg.RocketDownCommand)
			} else if event.Key == keyboard.KeyArrowUp {
				g.RightRacket.Move(g.Cfg.RocketUpCommand)
			} else if event.Key == keyboard.KeyArrowDown {
				g.RightRacket.Move(g.Cfg.RocketDownCommand)
			}
		}
	}
}
