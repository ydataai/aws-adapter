package clients

import (
	"log"

	"github.com/ydataai/aws-adapter/pkg/common"
	"github.com/ydataai/go-core/pkg/common/logging"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/servicequotas"
)

const (
	gpuQuotaCode    string  = "L-417A185B"
	vCPUToGPUFactor float64 = 4
)

// ServiceQuotaClientInterface defines an interface for service quota client
type ServiceQuotaClientInterface interface {
	GetAvailableGPUInstances(string, string) (common.GPU, error)
}

// ServiceQuotaClient is the service quota client
type ServiceQuotaClient struct {
	log          logging.Logger
	serviceQuota *servicequotas.ServiceQuotas
}

// NewServiceQuotaClient initializes service quota
func NewServiceQuotaClient(log logging.Logger, serviceQuota *servicequotas.ServiceQuotas) ServiceQuotaClient {
	return ServiceQuotaClient{
		log:          log,
		serviceQuota: serviceQuota,
	}
}

// GetAvailableGPUInstances fetchs available gpu instances in service quota
func (sq ServiceQuotaClient) GetAvailableGPUInstances(
	gpuInstanceType string,
	gpuQuotaServiceCode string,
) (common.GPU, error) {
	sq.log.Info("Starting to featch Available GPU instances")

	quota, err := sq.serviceQuota.GetServiceQuota(&servicequotas.GetServiceQuotaInput{
		QuotaCode:   aws.String(gpuQuotaCode),
		ServiceCode: aws.String(gpuQuotaServiceCode),
	})
	if err != nil {
		log.Fatal(err)
		return 0, err
	}

	availableInstances := *quota.Quota.Value / vCPUToGPUFactor

	return common.GPU(availableInstances), nil
}
