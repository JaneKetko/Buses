package serverstore

import (
	"time"

	"github.com/JaneKetko/Buses/api/proto"
	"github.com/JaneKetko/Buses/src/stores/domain"

	"github.com/golang/protobuf/ptypes"
)

//RouteForServer - struct for storing info about route for decoding and encoding.
type RouteForServer struct {
	ID        int             `json:"id"`
	Points    PointsForServer `json:"points"`
	Start     time.Time       `json:"start_time"`
	Cost      float32         `json:"cost"`
	FreeSeats int             `json:"freeseats"`
	AllSeats  int             `json:"allseats"`
}

//PointsForServer - struct for showing points of route for decoding and encoding.
type PointsForServer struct {
	StartPoint string `json:"startpoint"`
	EndPoint   string `json:"endpoint"`
}

//TicketForServer - struct for ticket info.
type TicketForServer struct {
	Points PointsForServer `json:"route"`
	Start  time.Time       `json:"start_time"`
	Cost   float32         `json:"cost"`
	Place  int             `json:"place"`
}

//TicketToJSON converts domain.Ticket to TicketServer.
func TicketToJSON(t domain.Ticket) TicketForServer {
	cost := float32(t.Cost) / 100
	ticket := TicketForServer{
		Points: PointsForServer{
			StartPoint: t.Points.StartPoint,
			EndPoint:   t.Points.EndPoint},
		Start: t.StartTime,
		Cost:  cost,
		Place: t.Place,
	}
	return ticket
}

//RouteServerToRoute convert routeServer to domain.Route.
func RouteServerToRoute(rServer RouteForServer) domain.Route {
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

//RouteToRouteServer convert domain.Route to RouteForServer.
func RouteToRouteServer(r domain.Route) RouteForServer {
	cost := float32(r.Cost) / 100
	route := RouteForServer{
		ID: r.ID,
		Points: PointsForServer{
			StartPoint: r.Points.StartPoint,
			EndPoint:   r.Points.EndPoint},
		Start:     r.Start,
		Cost:      cost,
		FreeSeats: r.FreeSeats,
		AllSeats:  r.AllSeats,
	}
	return route
}

//PtypeBusToJSON converts Bus in proto to RouteForServer.
func PtypeBusToJSON(in *proto.BusRoute) (*RouteForServer, error) {
	cost := float32(in.Cost) / 100
	date, err := ptypes.Timestamp(in.Start)
	if err != nil {
		return nil, err
	}

	route := &RouteForServer{
		ID: int(in.ID),
		Points: PointsForServer{
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

//RouteToPType converts domain.Route to Bus in proto.
func RouteToPType(route domain.Route) (*proto.BusRoute, error) {
	date, err := ptypes.TimestampProto(route.Start)
	if err != nil {
		return nil, err
	}
	busroute := &proto.BusRoute{
		ID: int64(route.ID),
		Points: &proto.RoutePoints{
			StartPoint: route.Points.StartPoint,
			EndPoint:   route.Points.EndPoint,
		},
		Start:     date,
		Cost:      int64(route.Cost),
		FreeSeats: int64(route.FreeSeats),
		AllSeats:  int64(route.AllSeats),
	}
	return busroute, nil
}

//TicketPTypeToJSON converts proto.Ticket to TicketServer.
func TicketPTypeToJSON(in *proto.Ticket) (*TicketForServer, error) {
	cost := float32(in.Cost) / 100
	date, err := ptypes.Timestamp(in.Start)
	if err != nil {
		return nil, err
	}

	ticket := &TicketForServer{
		Points: PointsForServer{
			StartPoint: in.Points.StartPoint,
			EndPoint:   in.Points.EndPoint,
		},
		Start: date,
		Cost:  cost,
		Place: int(in.Place),
	}

	return ticket, nil
}
