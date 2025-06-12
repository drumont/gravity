package main

import (
	"context"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"gravity/gravity-agent/services"
	"net"
	"os"

	cb "gravity/proto/container/pb"
)

type gravityServer struct {
	cb.UnimplementedContainerServiceServer
}

func initLogger() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
}

func (s *gravityServer) RunContainer(ctx context.Context, req *cb.RunContainerRequest) (*cb.RunContainerResponse, error) {
	return services.Run(ctx, req)
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	cb.RegisterContainerServiceServer(grpcServer, &gravityServer{})

	log.Println("Gravity agent server is listening on :50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
