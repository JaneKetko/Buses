package serverstore

import (
	"time"

	"github.com/JaneKetko/Buses/api/proto"
	"github.com/JaneKetko/Buses/src/stores/domain"

	"github.com/golang/protobuf/ptypes"
)

//RouteServer - struct for storing info about route for decoding and encoding.
type RouteServer struct {
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

//TicketServer - struct for ticket info.
type TicketServer struct {
	Points PointsServer `json:"route"`
	Start  time.Time    `json:"start_time"`
	Cost   float32      `json:"cost"`
	Place  int          `json:"place"`
}

//ConvertTicket converts domain.Ticket to TicketServer.
func ConvertTicket(t domain.Ticket) TicketServer {
	cost := float32(t.Cost) / 100
	ticket := TicketServer{
		Points: PointsServer{
			StartPoint: t.Points.StartPoint,
			EndPoint:   t.Points.EndPoint},
		Start: t.StartTime,
		Cost:  cost,
		Place: t.Place,
	}
	return ticket
}

//RouteServerToRoute convert routeServer to domain.Route.
func RouteServerToRoute(rServer RouteServer) domain.Route {
	cost := int(rServer.Cost * 100)
	route := domain.Route{
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

//RouteToRouteServer convert domain.Route to routeServer.
func RouteToRouteServer(r domain.Route) RouteServer {
	cost := float32(r.Cost) / 100
	route := RouteServer{
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

//PtypeBusToJSON converts Bus in proto to RouteServer.
func PtypeBusToJSON(in *proto.BusRoute) (*RouteServer, error) {
	cost := float32(in.Cost) / 100
	date, err := ptypes.Timestamp(in.Start)
	if err != nil {
		return nil, err
	}

	route := &RouteServer{
		ID: int(in.ID),
		Points: PointsServer{
			StartPoint: in.Points.StartPoint,
			EndPoint:   in.Points.EndPoint,
		},
		Start:     date,
		Cost:      cost,
		FreeSeats: int(in.FreeSeats),
		AllSeats:  int(in.AllSeats),
	}

	return route, nil
}

//ConvertTicketPType converts proto.Ticket to TicketServer.
func ConvertTicketPType(in *proto.Ticket) (*TicketServer, error) {
	cost := float32(in.Cost) / 100
	date, err := ptypes.Timestamp(in.Start)
	if err != nil {
		return nil, err
	}

	ticket := &TicketServer{
		Points: PointsServer{
			StartPoint: in.Points.StartPoint,
			EndPoint:   in.Points.EndPoint,
		},
		Start: date,
		Cost:  cost,
		Place: int(in.Place),
	}

	return ticket, nil
}
