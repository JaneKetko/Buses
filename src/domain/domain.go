package domain

import (
	"time"
)

//Route - struct for describing route of any bus.
type Route struct {
	ID        int       `json:"id"`
	Points    Points    `json:"points"`
	Start     time.Time `json:"start_time"`
	Cost      float32   `json:"cost"`
	FreeSeats int       `json:"freeseats"`
	AllSeats  int       `json:"allseats"`
}

//Points - struct for showing points of route.
type Points struct {
	StartPoint string `json:"startpoint"`
	EndPoint   string `json:"endpoint"`
}
