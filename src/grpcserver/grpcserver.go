package grpcserver

import (
	"context"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	pb "github.com/JaneKetko/Buses/api/proto"
	"github.com/JaneKetko/Buses/src/config"
	"github.com/JaneKetko/Buses/src/routemanager"
	"github.com/JaneKetko/Buses/src/stores/domain"

	"github.com/golang/protobuf/ptypes"
)

//GRPSServer - struct for server.
type GRPSServer struct {
	manager *routemanager.RouteManager
	config  *config.Config
}

//NewGRPSServer - init server.
func NewGRPSServer(r *routemanager.RouteManager, c *config.Config) *GRPSServer {
	return &GRPSServer{
		manager: r,
		config:  c,
	}
}

//RunServer - start serving.
func (s *GRPSServer) RunServer() {

	lis, err := net.Listen("tcp", s.config.PortGRPCServer)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	srv := grpc.NewServer()
	pb.RegisterBusesManagerServer(srv, s)
	reflection.Register(srv)
	log.Println("Started server...")
	if err := srv.Serve(lis); err != nil {
		log.Fatalf("errors with serving: %v", err)
	}

}

func convertTypes(route domain.Route) (*pb.BusRoute, error) {
	date, err := ptypes.TimestampProto(route.Start)
	if err != nil {
		return nil, err
	}
	busroute := &pb.BusRoute{
		ID: int64(route.ID),
		Points: &pb.RoutePoints{
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

//GetRoutes - get only cuurent routes.
func (s *GRPSServer) GetRoutes(ctx context.Context, in *pb.Nothing) (*pb.ListRoutes, error) {
	routes, err := s.manager.GetCurrentRoutes(ctx)
	if err != nil {
		return nil, err
	}

	busrt := make([]*pb.BusRoute, 0)
	for _, route := range routes {
		busroute, err := convertTypes(route)
		if err != nil {
			return nil, err
		}
		busrt = append(busrt, busroute)
	}

	return &pb.ListRoutes{BusRoutes: busrt}, nil
}

//GetRoute - get route by id.
func (s *GRPSServer) GetRoute(ctx context.Context, in *pb.IDRequest) (*pb.SingleRoute, error) {
	route, err := s.manager.GetRouteByID(ctx, int(in.ID))
	if err != nil {
		return nil, err
	}

	busroute, err := convertTypes(*route)
	if err != nil {
		return nil, err
	}

	return &pb.SingleRoute{Route: busroute}, nil
}

//CreateRoute - create route.
func (s *GRPSServer) CreateRoute(ctx context.Context, in *pb.SingleRoute) (*pb.Nothing, error) {
	date, err := ptypes.Timestamp(in.Route.Start)
	if err != nil {
		return nil, err
	}

	route := &domain.Route{
		ID: int(in.Route.ID),
		Points: domain.Points{
			StartPoint: in.Route.Points.StartPoint,
			EndPoint:   in.Route.Points.EndPoint,
		},
		Start:     date,
		Cost:      int(in.Route.Cost),
		FreeSeats: int(in.Route.FreeSeats),
		AllSeats:  int(in.Route.AllSeats),
	}

	err = s.manager.CreateRoute(ctx, route)
	if err != nil {
		return nil, err
	}

	return &pb.Nothing{}, nil
}

//DeleteRoute - delete route by id.
func (s *GRPSServer) DeleteRoute(ctx context.Context, in *pb.IDRequest) (*pb.Nothing, error) {
	err := s.manager.DeleteRoute(ctx, int(in.ID))
	if err != nil {
		return nil, err
	}

	return &pb.Nothing{}, nil
}

//BuyTicket - buy ticket for one route.
func (s *GRPSServer) BuyTicket(ctx context.Context, in *pb.IDRequest) (*pb.TicketResponse, error) {
	tick, err := s.manager.BuyTicket(ctx, int(in.ID))
	if err != nil {
		return nil, err
	}

	date, err := ptypes.TimestampProto(tick.StartTime)
	if err != nil {
		return nil, err
	}

	ticket := &pb.Ticket{
		Points: &pb.RoutePoints{
			StartPoint: tick.Points.StartPoint,
			EndPoint:   tick.Points.EndPoint,
		},
		Start: date,
		Cost:  int64(tick.Cost),
		Place: int64(tick.Place),
	}

	return &pb.TicketResponse{Ticket: ticket}, nil
}

//SearchRoutes - search routes with datetime and endpoint.
func (s *GRPSServer) SearchRoutes(ctx context.Context, in *pb.Search) (*pb.ListRoutes, error) {
	date, err := time.Parse("2006-01-02", in.StartTime)
	if err != nil {
		return nil, domain.ErrInvalidDate
	}
	routes, err := s.manager.SearchRoutes(ctx, date, in.EndPoint)
	if err != nil {
		return nil, err
	}

	busrt := make([]*pb.BusRoute, 0)
	for _, route := range routes {
		busroute, err := convertTypes(route)
		if err != nil {
			return nil, err
		}
		busrt = append(busrt, busroute)
	}

	return &pb.ListRoutes{BusRoutes: busrt}, nil
}
