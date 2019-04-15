package domain

import (
	"errors"
	"time"
)

var (
	//ErrNoFreeSeats - error for no free seats in bus.
	ErrNoFreeSeats = errors.New("no free seats in this bus")
	//ErrNoRoutes - error for no such route.
	ErrNoRoutes = errors.New("no such route")
	//ErrTypes - error with types.
	ErrTypes = errors.New("errors with types")
	//ErrNoRoutesByEndPoint - error for no such routes by this endpoint.
	ErrNoRoutesByEndPoint = errors.New("no such routes by this endpoint")
	//ErrInvalidDate - error for invalid date argument.
	ErrInvalidDate = errors.New("invalid date argument")
	//ErrInvalidArg - error with invallid arguments.
	ErrInvalidArg = errors.New("invalid arguments")
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

//Ticket - struct for storing info about taked ticket.
type Ticket struct {
	Points    Points
	StartTime time.Time
	Cost      int
	Place     int
}
