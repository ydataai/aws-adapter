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
	logger        logging.Logger
	ec2Client     EC2Client
	quotaClient   QuotaClient
	configuration ServiceConfiguration
}

// NewService initializes rest service
func NewService(
	logger logging.Logger,
	ec2Client EC2Client,
	quotaClient QuotaClient,
	configuration ServiceConfiguration,
) Service {
	return service{
		logger:        logger,
		ec2Client:     ec2Client,
		quotaClient:   quotaClient,
		configuration: configuration,
	}
}

// AvailableGPU ..
func (rs service) AvailableGPU(ctx context.Context) (GPU, error) {
	rs.logger.Infof("Fetching available GPUs to %v", rs.configuration.GPUInstanceType)

	runningGPUInstances, err := rs.ec2Client.GetGPUInstances(ctx, rs.configuration.GPUInstanceType)
	if err != nil {
		rs.logger.Errorf("while fetching gpu instances. Error: %+v", err)
		return GPU(0), err
	}

	rs.logger.Infof("Running GPUs %s instances: %f", rs.configuration.GPUInstanceType, runningGPUInstances)

	availableInstances, err := rs.quotaClient.GetAvailableQuota(
		ctx, rs.configuration.GPUQuotaCode, rs.configuration.GPUQuotaServiceCode)
	if err != nil {
		rs.logger.Errorf("while fetching available instances. Error: %+v", err)
		return GPU(0), err
	}

	availableGPUInstances := GPU(availableInstances / float64(rs.configuration.GPUvCPUFactor))

	rs.logger.Infof("Available GPUs %s instances: %f", rs.configuration.GPUInstanceType, availableGPUInstances)

	gpus := availableGPUInstances - runningGPUInstances

	return gpus, nil
}
