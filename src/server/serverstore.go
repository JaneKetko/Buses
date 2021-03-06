package server

import (
	"time"

	"github.com/JaneKetko/Buses/src/domain"
)

//RouteServer - struct for storing info about route for decoding and encoding.
type routeServer struct {
	ID        int          `json:"id"`
	Points    PointsServer `json:"points"`
	Start     time.Time    `json:"start_time"`
	Cost      float32      `json:"cost"`
	FreeSeats int          `json:"freeseats"`
	AllSeats  int          `json:"allseats"`
}

//PointsServer - struct for showing points of route for decoding and encoding.
type PointsServer struct {
	StartPoint string `json:"startpoint"`
	EndPoint   string `json:"endpoint"`
}

//routeServerToRoute convert routeServer to Route
func routeServerToRoute(rServer routeServer) domain.Route {
	var route domain.Route
	cost := int(rServer.Cost * 100)
	route = domain.Route{
		ID: rServer.ID,
		Points: domain.Points{
			StartPoint: rServer.Points.StartPoint,
			EndPoint:   rServer.Points.EndPoint},
		Start:     rServer.Start,
		Cost:      cost,
		FreeSeats: rServer.FreeSeats,
		AllSeats:  rServer.AllSeats,
	}
	return route
}

//routeToRouteServer convert Route to routeServer
func routeToRouteServer(r domain.Route) routeServer {
	var route routeServer
	cost := float32(r.Cost) / 100
	route = routeServer{
		ID: r.ID,
		Points: PointsServer{
			StartPoint: r.Points.StartPoint,
			EndPoint:   r.Points.EndPoint},
		Start:     r.Start,
		Cost:      cost,
		FreeSeats: r.FreeSeats,
		AllSeats:  r.AllSeats,
	}
	return route
}
