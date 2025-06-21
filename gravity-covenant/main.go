package main

import (
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"gravity/gravity-covenant/features"
	cb "gravity/proto/container/pb"
	"os"
)

func initLogger() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
}

func main() {
	// The main function is intentionally left empty.
	// This is a placeholder for the worker implementation.
	// You can add your worker logic here, such as connecting to the master,
	// sending heartbeats, and handling responses.

	// Example:
	// - Connect to the master server
	// - Send periodic heartbeat messages
	// - Handle responses from the master

	// Read from configuration or command line arguments
	// - Initialize logging
	// - Set up error handling
	// - Implement the worker's main loop
	// - Clean up resources on exit
	// - Handle graceful shutdown
	// - Implement retry logic for network failures
	// - Use context for cancellation and timeouts
	// - Implement worker-specific logic, such as task processing or status reporting
	// - Optionally, implement a health check mechanism
	// - Optionally, implement a shutdown signal handler
	// - Optionally, implement a metrics reporting mechanism
	// - Optionally, implement a configuration reload mechanism
	// - Optionally, implement a logging framework
	// - Optionally, implement a monitoring system
	// initLogger()
	var id = uuid.New().String()
	log.Infof("Starting covenant on local with id: %s", id)
	var agentAddress string = "192.168.1.17:50051"

	log.Infof("Connecting with agent at %s", agentAddress)
	connexion, err := grpc.NewClient(agentAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to agent: %v", err)
	}
	defer connexion.Close()
	log.Infof("Connected to agent at %s", agentAddress)

	containerClient := cb.NewContainerServiceClient(connexion)

	req := &cb.RunContainerRequest{
		RequestId: id,
		Memory:    1024, // 1 GB
		Vcpu:      1,
		Image:     "nginx:latest",
		Env:       map[string]string{"ENV_VAR": "value"},
		Ports:     make([]int32, 0),
	}

	resp, err := features.RunContainer(containerClient, req)
	if err != nil {
		log.Fatalf("Failed to run container: %v", err)
	}
	log.Infof("Container registered successfully with ID: %s", resp.ContainerId)

	err = features.StreamContainerLogs(containerClient, &cb.StreamContainerLogsRequest{
		ContainerId: resp.ContainerId,
	})
	if err != nil {
		log.Fatalf("Failed to stream container logs: %v", err)
	}
}
