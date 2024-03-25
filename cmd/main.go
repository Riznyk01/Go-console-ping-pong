package main

import (
	"github.com/eiannone/keyboard"
	"github.com/gopxl/beep"
	"github.com/gopxl/beep/speaker"
	"log"
	"os"
	"ping-pong/internal/audio_engine"
	"ping-pong/internal/config"
	"ping-pong/internal/controller"
	"ping-pong/internal/game"
	"ping-pong/internal/models"
)

func main() {
	logFile, err := os.OpenFile("logfile.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("error occurred while opening logfile:", err)
	}
	defer logFile.Close()
	log.SetOutput(logFile)
	cfg := config.MustLoad()
	log.Println(cfg)

	speaker.Init(beep.SampleRate(cfg.SampleRate), cfg.BufferSize)
	soundPlayer := audio_engine.NewAudioPlayer()
	soundPlayer.LoadSound(cfg.SoundPath)

	ball := game.NewBall(
		models.Coordinates{X: cfg.CenterCol, Y: cfg.CenterRow},
		models.Coordinates{X: cfg.CenterCol, Y: cfg.CenterRow},
		models.Coordinates{X: 1, Y: 1},
		soundPlayer,
		cfg)
	lRacket := game.NewRacket(models.Coordinates{X: 0, Y: cfg.CenterRow}, 4, soundPlayer, cfg)
	rRacket := game.NewRacket(models.Coordinates{X: cfg.WindowWidth, Y: cfg.CenterRow}, 4, soundPlayer, cfg)
	consoleScreen := game.NewScreen(cfg)
	keysChan := make(chan keyboard.KeyEvent)
	pong := game.NewGame(ball, lRacket, rRacket, consoleScreen, soundPlayer, cfg, keysChan)

	go controller.Controller(keysChan)
	pong.Play()
}
