package main

import (
	"context"
	"fmt"
	"os"

	"github.com/110y/run"
	"google.golang.org/grpc"

	"github.com/110y/grpc-go-connection-demo/callee/pb"
	pkggrpc "github.com/110y/grpc-go-connection-demo/grpc"
)

var _ pb.CalleeServiceServer = (*server)(nil)

func main() {
	run.Run(func(ctx context.Context) int {
		s := newServer()

		if err := s.Start(ctx); err != nil {
			fmt.Fprintf(os.Stderr, "failed to start or stop the server: %s", err)
			return 1
		}

		return 0
	})
}

func newServer() *pkggrpc.Server {
	s := &server{}
	return pkggrpc.NewServer(5000, func(gs *grpc.Server) {
		pb.RegisterCalleeServiceServer(gs, s)
	})
}

type server struct {
	pb.UnimplementedCalleeServiceServer
}

func (s *server) GetItem(ctx context.Context, _ *pb.GetItemRequest) (*pb.GetItemResponse, error) {
	return &pb.GetItemResponse{
		Name: "callee",
	}, nil
}
