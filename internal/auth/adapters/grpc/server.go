package grpc

import (
	"net"

	"google.golang.org/grpc"

	"github.com/cardio-analyst/backend/internal/auth/port/service"
	"github.com/cardio-analyst/backend/pkg/api/proto/auth"
)

type Server struct {
	server   *grpc.Server
	listener net.Listener
	services service.Services
	auth.UnimplementedAuthServiceServer
}

func NewServer(server *grpc.Server, listener net.Listener, services service.Services) *Server {
	var s Server
	s.server = server
	s.listener = listener
	s.services = services
	return &s
}

func (s *Server) Serve() error {
	return s.server.Serve(s.listener)
}

func (s *Server) Close() error {
	s.server.GracefulStop()
	return s.listener.Close()
}
