package game

import "ping-pong/internal/models"

func NewRocket(coord models.Coordinates, side int) *models.Rocket {
	return &models.Rocket{
		Coord: coord,
		Side:  side,
	}
}
