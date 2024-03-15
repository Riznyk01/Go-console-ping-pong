package models

type Coordinates struct {
	X int
	Y int
}

type Rocket struct {
	Coord Coordinates
	Side  int
}

type Ball struct {
	Coord     Coordinates
	Direction Coordinates
}
