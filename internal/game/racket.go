package game

import "ping-pong/internal/models"

func NewRacket(coord models.Coordinates, side int) *models.Racket {
	return &models.Racket{
		Coord: coord,
		Side:  side,
	}
}
