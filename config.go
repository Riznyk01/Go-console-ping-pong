package ping_pong

import (
	"runtime"
	"time"
)

type Config struct {
	WindowWidth       int
	WindowHeight      int
	Title             string
	Filler            rune
	EmptyFiller       rune
	HeaderFiller      rune
	FooterFiller      rune
	RacketsLineFiller rune
	BallMovingDelay   time.Duration
	RenderDelay       time.Duration
	CenterRow         int
	CenterCol         int
	FailPause         time.Time
	SampleRate        int
	BufferSize        int
	RocketUpCommand   string
	RocketDownCommand string
}

func NewConfig() (cfg *Config) {
	cfg = &Config{}

	if runtime.GOOS == "windows" {
		cfg.WindowWidth = 120
		cfg.WindowHeight = 27
	} else {
		cfg.WindowWidth = 80
		cfg.WindowHeight = 22
	}

	cfg.Title = " Go console PING-PONG "
	cfg.Filler = '▒'
	cfg.EmptyFiller = 'ₓ'
	cfg.HeaderFiller = '▬'
	cfg.FooterFiller = '▬'
	cfg.RacketsLineFiller = '⁞'
	cfg.BallMovingDelay = 40 * time.Millisecond
	cfg.RenderDelay = 40 * time.Millisecond
	cfg.CenterRow = (cfg.WindowHeight / 2) + 1
	cfg.CenterCol = (cfg.WindowWidth / 2) + 1
	cfg.SampleRate = 44100
	cfg.BufferSize = 1411
	cfg.RocketUpCommand = "up"
	cfg.RocketDownCommand = "down"

	return cfg
}
