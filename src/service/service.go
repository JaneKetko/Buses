package service

import (
	"github.com/JaneKetko/Buses/src/grpcserver"
	"github.com/JaneKetko/Buses/src/server"
)

//Service - struct for bus servers.
type Service struct {
	grpc *grpcserver.GRPSServer
	rest *server.RESTServer
}

//NewService - constructor for Service.
func NewService(g *grpcserver.GRPSServer, r *server.RESTServer) *Service {
	return &Service{g, r}
}

//RunService - start service.
func (s *Service) RunService() {
	go s.grpc.RunServer()
	s.rest.RunServer()
}
