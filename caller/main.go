package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/110y/run"
	"google.golang.org/grpc"

	calleepb "github.com/110y/grpc-go-connection-demo/callee/pb"
	"github.com/110y/grpc-go-connection-demo/caller/pb"
	pkggrpc "github.com/110y/grpc-go-connection-demo/grpc"
)

var _ pb.CallerServiceServer = (*server)(nil)

func main() {
	run.Run(func(ctx context.Context) int {
		s, err := newServer(ctx)
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

func newServer(ctx context.Context) (*pkggrpc.Server, error) {
	opts := []grpc.DialOption{
		grpc.WithInsecure(),
		grpc.WithBlock(),
		grpc.WithDefaultCallOptions(grpc.WaitForReady(true)),
	}

	calleeHost := os.Getenv("CALLEE_HOST")
	if calleeHost == "" {
		return nil, errors.New("must specify CALLEE_HOST")
	}

	calleePort := os.Getenv("CALLEE_PORT")
	if calleePort == "" {
		return nil, errors.New("must specify CALLEE_PORT")
	}

	conn, err := grpc.DialContext(ctx, fmt.Sprintf("%s:%s", calleeHost, calleePort), opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to dial to the callee server: %w", err)
	}

	calleeClient := calleepb.NewCalleeServiceClient(conn)

	s := &server{calleeClient: calleeClient}

	port := os.Getenv("PORT")
	if port == "" {
		return nil, errors.New("must specify PORT")
	}

	portnum, err := strconv.Atoi(port)
	if err != nil {
		return nil, errors.New("must specify valid PORT as number")
	}

	return pkggrpc.NewServer(portnum, func(gs *grpc.Server) {
		pb.RegisterCallerServiceServer(gs, s)
	}), nil
}

type server struct {
	calleeClient calleepb.CalleeServiceClient
	pb.UnimplementedCallerServiceServer
}

func (s *server) GetItem(ctx context.Context, _ *pb.GetItemRequest) (*pb.GetItemResponse, error) {
	res, err := s.calleeClient.GetItem(ctx, &calleepb.GetItemRequest{Id: "f8ba61b5-3332-404f-9b64-e594f7335fa6"})
	if err != nil {
		return nil, fmt.Errorf("failed to call callee: %w", err)
	}

	return &pb.GetItemResponse{
		Name: fmt.Sprintf("caller gets item from callee: %s", res.Name),
	}, nil
}
