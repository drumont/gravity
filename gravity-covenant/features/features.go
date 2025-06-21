package features

import (
	"context"
	log "github.com/sirupsen/logrus"
	cb "gravity/proto/container/pb"
	"time"
)

func RunContainer(client cb.ContainerServiceClient, req *cb.RunContainerRequest) (*cb.RunContainerResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	resp, err := client.RunContainer(ctx, req)

	if err != nil {
		log.Fatalf("Failed to run container: %v", err)
	} else {
		log.Infof("Container registered successfully with ID: %s", resp.ContainerId)
	}

	return resp, err
}

func StreamContainerLogs(client cb.ContainerServiceClient, req *cb.StreamContainerLogsRequest) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	stream, err := client.StreamContainerLogs(ctx, req)
	if err != nil {
		log.Fatalf("Failed to stream container logs: %v", err)
		return err
	}

	for {
		line, err := stream.Recv()
		if err != nil {
			log.Errorf("Error receiving log line: %v", err)
			return err
		}
		log.Printf("%v: %v", line.ContainerId, line.Log)
	}
}
