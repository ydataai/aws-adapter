// Package usage offers objects and methods to help using usage APIs
package usage

import (
	"log"

	"github.com/ydataai/go-core/pkg/common/logging"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/servicequotas"
)

const (
	vCPUToGPUFactor float64 = 4
)

// QuotaClient defines an interface for service quota client
type QuotaClient interface {
	GetAvailableQuota(string, string) (float64, error)
}

type quotaClient struct {
	logger       logging.Logger
	serviceQuota *servicequotas.ServiceQuotas
}

// NewServiceQuotaClient initializes service quota
func NewQuotaClient(logger logging.Logger, serviceQuota *servicequotas.ServiceQuotas) QuotaClient {
	return quotaClient{
		logger:       logger,
		serviceQuota: serviceQuota,
	}
}

// GetAvailableQuota fetchs available gpu instances in service quota
func (sq quotaClient) GetAvailableQuota(gpuQuotaCode string, gpuQuotaServiceCode string) (float64, error) {
	sq.logger.Infof("Starting to fetch Available %s GPU instances for quota code", gpuQuotaCode)

	quota, err := sq.serviceQuota.GetServiceQuota(&servicequotas.GetServiceQuotaInput{
		QuotaCode:   aws.String(gpuQuotaCode),
		ServiceCode: aws.String(gpuQuotaServiceCode),
	})
	if err != nil {
		log.Fatal(err)
		return 0, err
	}

	availableQuota := *quota.Quota.Value
	sq.logger.Infof("Available quota for %s/%s: %f", gpuQuotaServiceCode, gpuQuotaCode, availableQuota)

	return availableQuota, nil
}
