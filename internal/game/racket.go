package game

import (
	"ping-pong/internal/audio_engine"
	ping_pong "ping-pong/internal/config"
	"ping-pong/internal/models"
)

type Racket struct {
	coord  models.Coordinates
	side   int
	sound  *audio_engine.AudioPlayer
	config *ping_pong.Config
}

func NewRacket(coord models.Coordinates, side int, soundEngine *audio_engine.AudioPlayer, cfg *ping_pong.Config) *Racket {
	return &Racket{
		coord:  coord,
		side:   side,
		sound:  soundEngine,
		config: cfg,
	}
}

func (r *Racket) Move(command string) {
	if command == r.config.RocketUpCommand {
		if r.coord.Y-r.side > 1 {
			r.coord.Y -= 1
			r.MoveSound()
			if r.coord.Y-r.side == 1 {
				r.HitSound()
			}
		}
	} else if command == r.config.RocketDownCommand {
		if r.coord.Y+r.side < r.config.WindowHeight {
			r.MoveSound()
			r.coord.Y += 1
			if r.coord.Y+r.side == r.config.WindowHeight {
				r.HitSound()
			}
		}
	}
}
func (r *Racket) HitSound() {
	r.sound.Play("wood_hit.wav")
}
func (r *Racket) MoveSound() {
	r.sound.Play("rackets_moving.wav")
}
