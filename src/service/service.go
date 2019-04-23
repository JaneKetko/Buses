package service

import (
	"github.com/JaneKetko/Buses/src/grpcserver"
	"github.com/JaneKetko/Buses/src/server"
)

//Service - struct for bus servers.
type Service struct {
	Grpc *grpcserver.GRPCServer
	Rest *server.RESTServer
}

//NewService - constructor for Service.
func NewService(g *grpcserver.GRPCServer, r *server.RESTServer) *Service {
	return &Service{g, r}
}

//RunService - start service.
func (s *Service) RunService() {
	go s.Grpc.RunServer()
	s.Rest.RunServer()
}

//StopService - stop service.
func (s *Service) StopService() {
	s.Grpc.Srv.GracefulStop()
}
