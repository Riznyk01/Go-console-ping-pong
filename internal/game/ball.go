package game

import "ping-pong/internal/models"

func NewBall(coord models.Coordinates, lastCoord models.Coordinates, direction models.Coordinates) *models.Ball {
	return &models.Ball{
		Coord:     coord,
		LastCoord: lastCoord,
		Direction: direction,
	}
}
