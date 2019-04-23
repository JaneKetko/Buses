package grpcserver

import (
	"context"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	pb "github.com/JaneKetko/Buses/api/proto"
	"github.com/JaneKetko/Buses/src/routemanager"
	"github.com/JaneKetko/Buses/src/stores/domain"
	sst "github.com/JaneKetko/Buses/src/stores/serverstore"

	"github.com/golang/protobuf/ptypes"
)

//GRPCServer - struct for server.
type GRPCServer struct {
	manager *routemanager.RouteManager
	address string
	Srv     *grpc.Server
}

//NewGRPSServer - init server.
func NewGRPCServer(r *routemanager.RouteManager, addr string) *GRPCServer {
	return &GRPCServer{
		manager: r,
		address: addr,
		Srv:     grpc.NewServer(),
	}
}

//RunServer - start serving.
func (s *GRPCServer) RunServer() {

	lis, err := net.Listen("tcp", s.address)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	pb.RegisterBusesManagerServer(s.Srv, s)
	reflection.Register(s.Srv)
	log.Println("Started server...")
	if err := s.Srv.Serve(lis); err != nil {
		log.Fatalf("errors with serving: %v", err)
	}

}

//GetRoutes - get only cuurent routes.
func (s *GRPCServer) GetRoutes(ctx context.Context, in *pb.Nothing) (*pb.ListRoutes, error) {
	routes, err := s.manager.GetCurrentRoutes(ctx)
	if err != nil {
		return nil, err
	}

	busrt := make([]*pb.BusRoute, 0)
	for _, route := range routes {
		busroute, err := sst.RouteToPType(route)
		if err != nil {
			return nil, err
		}
		busrt = append(busrt, busroute)
	}

	return &pb.ListRoutes{BusRoutes: busrt}, nil
}

//GetRoute - get route by id.
func (s *GRPCServer) GetRoute(ctx context.Context, in *pb.IDRequest) (*pb.SingleRoute, error) {
	route, err := s.manager.GetRouteByID(ctx, int(in.ID))
	if err != nil {
		return nil, err
	}

	busroute, err := sst.RouteToPType(*route)
	if err != nil {
		return nil, err
	}

	return &pb.SingleRoute{Route: busroute}, nil
}

//CreateRoute - create route.
func (s *GRPCServer) CreateRoute(ctx context.Context, in *pb.SingleRoute) (*pb.Nothing, error) {
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
func (s *GRPCServer) DeleteRoute(ctx context.Context, in *pb.IDRequest) (*pb.Nothing, error) {
	err := s.manager.DeleteRoute(ctx, int(in.ID))
	if err != nil {
		return nil, err
	}

	return &pb.Nothing{}, nil
}

//BuyTicket - buy ticket for one route.
func (s *GRPCServer) BuyTicket(ctx context.Context, in *pb.IDRequest) (*pb.TicketResponse, error) {
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
func (s *GRPCServer) SearchRoutes(ctx context.Context, in *pb.Search) (*pb.ListRoutes, error) {
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
		busroute, err := sst.RouteToPType(route)
		if err != nil {
			return nil, err
		}
		busrt = append(busrt, busroute)
	}

	return &pb.ListRoutes{BusRoutes: busrt}, nil
}
