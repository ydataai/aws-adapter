package service

import (
	"context"

	"github.com/ydataai/aws-quota-provider/pkg/clients"
	"github.com/ydataai/aws-quota-provider/pkg/common"
	"github.com/ydataai/go-core/pkg/common/logging"
)

// RESTServiceInterface defines rest service interface
type RESTServiceInterface interface {
	AvailableGPU(ctx context.Context) (common.GPU, error)
}

// RESTService defines a struct with required dependencies for rest service
type RESTService struct {
	log                logging.Logger
	ec2Client          clients.EC2ClientInterface
	serviceQuotaClient clients.ServiceQuotaClientInterface
	configuration      RESTServiceConfiguration
}

// NewRESTService initializes rest service
func NewRESTService(
	log logging.Logger,
	ec2Client clients.EC2ClientInterface,
	serviceQuotaClient clients.ServiceQuotaClientInterface,
	configuration RESTServiceConfiguration,
) RESTService {
	return RESTService{
		log:                log,
		ec2Client:          ec2Client,
		serviceQuotaClient: serviceQuotaClient,
		configuration:      configuration,
	}
}

// AvailableGPU ..
func (rs RESTService) AvailableGPU(ctx context.Context) (common.GPU, error) {
	rs.log.Info("Starting to featch available GPU")

	gpuInstances, err := rs.ec2Client.GetGPUInstances(rs.configuration.GPUInstanceType)
	if err != nil {
		rs.log.Infof("while fetching gpu instances. Error: %v", gpuInstances)
		return common.GPU(0), err
	}

	availableGPUInstances, err := rs.serviceQuotaClient.GetAvailableGPUInstances(
		rs.configuration.GPUInstanceType,
		rs.configuration.GPUQuotaServiceCode,
	)
	if err != nil {
		rs.log.Infof("while fetching available gpu instances. Error: %v", availableGPUInstances)
		return common.GPU(0), err
	}

	gpus := availableGPUInstances - gpuInstances

	return gpus, nil
}
