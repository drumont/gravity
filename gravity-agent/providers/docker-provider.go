package providers

import (
	"context"
	log "github.com/sirupsen/logrus"

	"github.com/docker/docker/api/types/container"
	dockerclient "github.com/docker/docker/client"
)

type DockerProvider struct {
	config map[string]string
	client *dockerclient.Client
}

func (dp *DockerProvider) NewDockerProvider(config map[string]string) *DockerProvider {
	dp.config = config
	if dp.client == nil {
		log.Debug("Docker client is not initialized, initializing now")
		var err error
		dp.client, err = initClient()
		if err != nil {
			log.Fatalf("Failed to initialize Docker client: %v", err)
			return nil
		}
		log.Debug("Docker client initialized successfully")
	} else {
		log.Debug("Docker client already initialized, reusing existing client")
	}
	return dp
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
