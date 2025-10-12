package app

import (
	"net"

	"github.com/cart-overflow/cart-api/pkg/pb"
	"github.com/cart-overflow/cart/internal/cart"
	"google.golang.org/grpc"
)

type Server struct {
	lis net.Listener
	rpc *grpc.Server
}

type Handlers struct {
	Cart *cart.Handler
}

func NewServer(addr string, h Handlers) (*Server, error) {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}

	rpc := grpc.NewServer()
	if h.Cart != nil {
		pb.RegisterCartServiceServer(rpc, h.Cart)
	}

	s := Server{lis, rpc}
	return &s, nil
}

func (s *Server) Run() error {
	return s.rpc.Serve(s.lis)
}

func (s *Server) Stop() {
	s.rpc.GracefulStop()
}
