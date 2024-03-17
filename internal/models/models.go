package models

type Coordinates struct {
	X int
	Y int
}

type Racket struct {
	Coord Coordinates
	Side  int
}

type Ball struct {
	Coord     Coordinates
	LastCoord Coordinates
	Direction Coordinates
}
