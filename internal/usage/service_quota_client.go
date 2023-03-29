// Package usage offers objects and methods to help using usage APIs
package usage

import (
	"log"

	"github.com/ydataai/go-core/pkg/common/logging"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/servicequotas"
)

const (
	gpuQuotaCode    string  = "L-417A185B"
	vCPUToGPUFactor float64 = 4
)

// ServiceQuotaClient defines an interface for service quota client
type ServiceQuotaClient interface {
	GetAvailableGPUInstances(string, string) (GPU, error)
}

type serviceQuotaClient struct {
	log          logging.Logger
	serviceQuota *servicequotas.ServiceQuotas
}

// NewServiceQuotaClient initializes service quota
func NewServiceQuotaClient(log logging.Logger, serviceQuota *servicequotas.ServiceQuotas) ServiceQuotaClient {
	return serviceQuotaClient{
		log:          log,
		serviceQuota: serviceQuota,
	}
}

// GetAvailableGPUInstances fetchs available gpu instances in service quota
func (sq serviceQuotaClient) GetAvailableGPUInstances(
	gpuInstanceType string,
	gpuQuotaServiceCode string,
) (GPU, error) {
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

	return GPU(availableInstances), nil
}
