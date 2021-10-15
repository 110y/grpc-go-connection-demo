package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/110y/run"
	"google.golang.org/grpc"

	"github.com/110y/grpc-go-connection-demo/callee/pb"
	pkggrpc "github.com/110y/grpc-go-connection-demo/grpc"
)

var _ pb.CalleeServiceServer = (*server)(nil)

func main() {
	run.Run(func(ctx context.Context) int {
		s, err := newServer()
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to create the server: %s", err)
			return 1
		}

		if err := s.Start(ctx); err != nil {
			fmt.Fprintf(os.Stderr, "failed to start or stop the server: %s", err)
			return 1
		}

		return 0
	})
}

func newServer() (*pkggrpc.Server, error) {
	s := &server{}

	port := os.Getenv("PORT")
	if port == "" {
		return nil, errors.New("must specify PORT")
	}

	portnum, err := strconv.Atoi(port)
	if err != nil {
		return nil, errors.New("must specify valid PORT as number")
	}

	return pkggrpc.NewServer(portnum, func(gs *grpc.Server) {
		pb.RegisterCalleeServiceServer(gs, s)
	}), nil
}

type server struct {
	pb.UnimplementedCalleeServiceServer
}

func (s *server) GetItem(ctx context.Context, _ *pb.GetItemRequest) (*pb.GetItemResponse, error) {
	return &pb.GetItemResponse{
		Name: "callee",
	}, nil
}
