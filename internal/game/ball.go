package game

import (
	"log"
	ping_pong "ping-pong"
	"ping-pong/internal/audio_engine"
	"ping-pong/internal/models"
	"time"
)

type Ball struct {
	coord       models.Coordinates
	lastCoord   models.Coordinates
	direction   models.Coordinates
	soundEngine *audio_engine.AudioPlayer
	config      *ping_pong.Config
}

func NewBall(coord models.Coordinates, lastcoord models.Coordinates, direction models.Coordinates, soundEngine *audio_engine.AudioPlayer, cfg *ping_pong.Config) *Ball {
	return &Ball{
		coord:       coord,
		lastCoord:   lastcoord,
		direction:   direction,
		soundEngine: soundEngine,
		config:      cfg,
	}
}

func (b *Ball) Move() {
	log.Printf("Move coord.X: %d coord.Y %d", b.coord.X, b.coord.Y)
	<-time.After(b.config.BallMovingDelay)
	b.coord.X += b.direction.X
	b.coord.Y += b.direction.Y

	if b.coord.X == b.config.WindowWidth {
		b.direction.X = -b.direction.X
		b.HitSound()
	} else if b.coord.X == 1 {
		b.direction.X = -b.direction.X
		b.HitSound()
	}
	if b.coord.Y == b.config.WindowHeight {
		b.HitSound()
		b.direction.Y = -b.direction.Y
	} else if b.coord.Y == 1 {
		b.HitSound()
		b.direction.Y = -b.direction.Y
	}
}
func (b *Ball) HitSound() {
	b.soundEngine.Play("hit.wav")
}
