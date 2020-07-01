package network

import (
	"net"

	"google.golang.org/grpc"
)

type ServiceProvider func(server *grpc.Server)

type Server struct {
	Addr string

	serviceProviders []ServiceProvider
}

func NewServer(addr string) *Server {
	return &Server{Addr: addr}
}

func (s *Server) ListenAndServe() error {
	addr := s.Addr
	if addr == "" {
		panic("unable to start server without an address")
	}

	l, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	return s.Serve(l)
}

func (s *Server) Serve(l net.Listener) error {
	grpcServer := grpc.NewServer()

	for _, sp := range s.serviceProviders {
		sp(grpcServer)
	}

	return grpcServer.Serve(l)
}

func (s *Server) RegisterServices(sp ServiceProvider) {
	s.serviceProviders = append(s.serviceProviders, sp)
}
