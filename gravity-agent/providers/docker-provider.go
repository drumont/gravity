package providers

import (
	"context"
	"github.com/docker/docker/api/types/container"
	dockerclient "github.com/docker/docker/client"
	log "github.com/sirupsen/logrus"
	"io"
)

type DockerProvider struct {
	config map[string]string
	client *dockerclient.Client
}

func NewDockerProvider(config map[string]string) *DockerProvider {
	log.Debug("Creating new DockerProvider instance")
	cli, err := initClient()
	if err != nil {
		log.Fatalf("Failed to initialize Docker client: %v", err)
		return nil
	}
	log.Debug("Docker client initialized successfully")

	return &DockerProvider{
		config: config,
		client: cli,
	}
}

func (dp *DockerProvider) LaunchContainer(ctx context.Context, request map[string]string) (container.CreateResponse, error) {
	log.Debugf("Launching container with ID")
	// Here you would typically use dp.client to launch the container
	// For now, we just log the action
	resp, err := dp.client.ContainerCreate(ctx, &container.Config{
		Image: request["image"],
		Tty:   false,
	}, nil, nil, nil, "")

	if err != nil {
		log.Errorf("Failed to create container: %v", err)
		return container.CreateResponse{}, err
	}

	if err := dp.client.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		log.Errorf("Failed to start container: %v", err)
		return container.CreateResponse{}, err
	}

	log.Infof("Container launched successfully")
	return resp, nil
}

func (dp *DockerProvider) RetrieveLogs(containerId string) (io.ReadCloser, error) {
	log.Infof("Retrieving logs for container ID: %s", containerId)
	// Here you would typically use dp.client to retrieve logs
	// For now, we just log the action
	// For now, we just log the action
	opts := container.LogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Follow:     false,
		Tail:       "all",
	}
	reader, err := dp.client.ContainerLogs(context.Background(), containerId, opts)
	if err != nil {
		log.Errorf("Failed to retrieve logs for container %s: %v", containerId, err)
		return nil, err
	}
	log.Debugf("Logs retrieved successfully for container ID: %s", containerId)
	return reader, nil
}

func initClient() (*dockerclient.Client, error) {
	log.Debug("Initializing Docker client with configuration")
	cli, err := dockerclient.NewClientWithOpts(dockerclient.FromEnv, dockerclient.WithAPIVersionNegotiation())
	if err != nil {
		log.Fatalf("Failed to create Docker client: %v", err)
		return nil, err
	}
	defer cli.Close()
	log.Debug("Docker client initialized successfully")
	return cli, nil
}
