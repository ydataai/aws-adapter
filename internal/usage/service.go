// Package usage offers objects and methods to help using usage APIs
package usage

import (
	"context"

	"github.com/ydataai/go-core/pkg/common/logging"
)

// Service defines rest service interface
type Service interface {
	AvailableGPU(ctx context.Context) (GPU, error)
}

type service struct {
	log                logging.Logger
	ec2Client          EC2Client
	serviceQuotaClient ServiceQuotaClient
	configuration      ServiceConfiguration
}

// NewService initializes rest service
func NewService(
	log logging.Logger,
	ec2Client EC2Client,
	serviceQuotaClient ServiceQuotaClient,
	configuration ServiceConfiguration,
) Service {
	return service{
		log:                log,
		ec2Client:          ec2Client,
		serviceQuotaClient: serviceQuotaClient,
		configuration:      configuration,
	}
}

// AvailableGPU ..
func (rs service) AvailableGPU(ctx context.Context) (GPU, error) {
	rs.log.Info("Starting to featch available GPU")

	gpuInstances, err := rs.ec2Client.GetGPUInstances(rs.configuration.GPUInstanceType)
	if err != nil {
		rs.log.Infof("while fetching gpu instances. Error: %v", gpuInstances)
		return GPU(0), err
	}

	availableGPUInstances, err := rs.serviceQuotaClient.GetAvailableGPUInstances(
		rs.configuration.GPUInstanceType,
		rs.configuration.GPUQuotaServiceCode,
	)
	if err != nil {
		rs.log.Infof("while fetching available gpu instances. Error: %v", availableGPUInstances)
		return GPU(0), err
	}

	gpus := availableGPUInstances - gpuInstances

	return gpus, nil
}
