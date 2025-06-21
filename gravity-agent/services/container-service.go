package services

import (
	"bufio"
	"context"
	"errors"
	log "github.com/sirupsen/logrus"
	"gravity/gravity-agent/helpers"
	"gravity/gravity-agent/providers"
	cb "gravity/proto/container/pb"
)

type ContainerService struct {
	dockerProvider *providers.DockerProvider
}

func NewContainerService(dockerProvider *providers.DockerProvider) *ContainerService {
	return &ContainerService{
		dockerProvider: dockerProvider,
	}
}

func (cs *ContainerService) Run(ctx context.Context, req *cb.RunContainerRequest) (*cb.RunContainerResponse, error) {
	log.Infof("Received registration for container ID: %s", req.RequestId)

	// Here you would typically handle the registration logic, such as storing the container info.
	// For now, we just log it and return a success response.

	if err := verifyHostResources(req.Memory, req.Vcpu); err != nil {
		log.Errorf("Failed to verify host resources: %v", err)
		return nil, err
	}

	log.Debugf("Host resources verified successfully for container ID: %s", req.RequestId)
	// Simulate running the container by logging the request details

	resp, err := cs.dockerProvider.LaunchContainer(ctx,
		map[string]string{
			"image": req.Image,
		})

	if err != nil {
		log.Errorf("Failed to launch container: %v", err)
		return nil, err
	}

	return &cb.RunContainerResponse{
		ContainerId: resp.ID,
	}, nil
}

func (cs *ContainerService) StreamContainerLogs(req *cb.StreamContainerLogsRequest,
	stream cb.ContainerService_StreamContainerLogsServer) error {
	log.Debugf("Streaming logs for container ID: %s", req.ContainerId)
	if req.ContainerId == "" {
		log.Error("Container ID is empty")
		return errors.New("container ID cannot be empty")
	}
	logReader, err := cs.dockerProvider.RetrieveLogs(req.ContainerId)
	if err != nil {
		log.Errorf("Failed to retrieve logs for container ID %s: %v", req.ContainerId, err)
	}
	defer logReader.Close()
	scanner := bufio.NewScanner(logReader)
	for scanner.Scan() {
		select {
		case <-stream.Context().Done():
			log.Infof("Stream context done for container ID %s", req.ContainerId)
			return stream.Context().Err()
		default:
			logLine := scanner.Text()
			if err := stream.Send(&cb.StreamContainerLogsResponse{
				ContainerId: req.ContainerId,
				Log:         logLine,
			}); err != nil {
				log.Errorf("Failed to send log stream response for container ID %s: %v", req.ContainerId, err)
				return err
			}
		}
	}
	if err := scanner.Err(); err != nil {
		log.Errorf("Error reading logs for container ID %s: %v", req.ContainerId, err)

		return err
	}

	return errors.New("no logs to stream")
}

func verifyHostResources(memory int64, vCpu float64) error {
	log.Debugf("Verifying host resources for container")
	if err := checkMemory(memory); err != nil {
		log.Errorf("Memory check failed: %v", err)
		return err
	}
	if err := checkVCpu(vCpu); err != nil { // Assuming 1 vCPU per GB of memory
		log.Errorf("vCPU check failed: %v", err)
		return err
	}
	return nil
}

func checkMemory(memory int64) error {
	log.Debugf("Verifying host resources")
	availableMemory, err := helpers.RetrieveHostAvailableMemory()
	if err != nil {
		log.Errorf("Failed to retrieve host available memory: %v", err)
		return err
	}
	if availableMemory < memory {
		log.Errorf("Not enough memory available on host: %d bytes available, %d bytes requested", availableMemory, memory)
		return errors.New("not enough memory available on host")
	}
	return nil
}

func checkVCpu(vCpu float64) error {
	log.Debugf("Checking vcpu")
	availableVCpu, err := helpers.RetrieveHostAvailableVCpu()
	if err != nil {
		log.Errorf("Failed to retrieve host available vCPU: %v", err)
		return err
	}
	if vCpu > availableVCpu {
		log.Errorf("Not enough vCPU available on host: %f vCPUs available, %f vCPUs requested", availableVCpu, vCpu)
		return errors.New("not enough vCPU available on host")
	}
	return nil
}
