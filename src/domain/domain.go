package domain

import (
	"time"
)

//Route - struct for describing route of any bus.
type Route struct {
	ID        int
	Points    Points
	Start     time.Time
	Cost      int
	FreeSeats int
	AllSeats  int
}

//Points - struct for showing points of route.
type Points struct {
	StartPoint string
	EndPoint   string
}

//Ticket - struct for storing info about taked ticket
type Ticket struct {
	Points    Points
	StartTime time.Time
	Cost      int
	Place     int
}
