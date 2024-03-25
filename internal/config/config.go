package config

import (
	"github.com/pelletier/go-toml"
	"log"
	"runtime"
	"time"
)

type Config struct {
	WindowWidth       int           `toml:"window_width"`
	WindowHeight      int           `toml:"window_height"`
	Title             string        `toml:"title"`
	Filler            rune          `toml:"filler"`
	EmptyFiller       rune          `toml:"empty_filler"`
	HeaderFiller      rune          `toml:"header_filler"`
	FooterFiller      rune          `toml:"footer_filler"`
	RacketsLineFiller rune          `toml:"rackets_line_filler"`
	BallMovingDelay   time.Duration `toml:"ball_moving_delay"`
	RenderDelay       time.Duration `toml:"render_delay"`
	CenterRow         int           `toml:"center_row"`
	CenterCol         int           `toml:"center_col"`
	SampleRate        int           `toml:"sample_rate"`
	BufferSize        int           `toml:"buffer_size"`
	RocketUpCommand   string        `toml:"rocket_up_command"`
	RocketDownCommand string        `toml:"rocket_down_command"`
	SoundPath         string        `toml:"sound_path"`
}

func MustLoad() *Config {

	configPath := "config/config.toml"
	// check if file exists
	cfg, err := toml.LoadFile(configPath)
	if err != nil {
		log.Fatalf("error loading config file: %s", err)
	}

	var config Config

	if err := cfg.Unmarshal(&config); err != nil {
		log.Fatalf("error decoding config: %s", err)
	}
	SetFieldResolution(&config)
	return &config
}
func SetFieldResolution(cfg *Config) {
	if runtime.GOOS == "windows" {
		cfg.WindowWidth = 120
		cfg.WindowHeight = 27
		cfg.CenterCol = cfg.WindowWidth / 2
		cfg.CenterRow = cfg.WindowHeight / 2
	} else {
		cfg.WindowWidth = 80
		cfg.WindowHeight = 24
		cfg.CenterCol = cfg.WindowWidth / 2
		cfg.CenterRow = cfg.WindowHeight / 2
	}
}
