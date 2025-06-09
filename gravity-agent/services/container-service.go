package services

import (
	"context"
	"errors"
	log "github.com/sirupsen/logrus"
	"gravity/gravity-agent/helpers"
	cb "gravity/proto/container/pb"
)

func Register(ctx context.Context, req *cb.RunContainerRequest) (*cb.RunContainerResponse, error) {
	log.Infof("Received registration for container ID: %s", req.RequestId)

	// Here you would typically handle the registration logic, such as storing the container info.
	// For now, we just log it and return a success response.

	if err := verifyHostResources(req.Memory); err != nil {
		log.Errorf("Failed to verify host resources: %v", err)
		return nil, err
	}

	return &cb.RunContainerResponse{
		ContainerId: "id",
	}, nil
}

func verifyHostResources(memory int64) error {
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
