package config

import (
	"github.com/pelletier/go-toml"
	"log"
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

	if err := cfg.Unmarshal(&config); err != nil {
		log.Fatalf("unable to decode config: %s", err)
	}

	return &config
}
