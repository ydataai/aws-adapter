// Package usage offers objects and methods to help using usage APIs
package usage

import "github.com/kelseyhightower/envconfig"

// ServiceConfiguration defines required configuration for rest service
type ServiceConfiguration struct {
	GPUInstanceType     string `envconfig:"GPU_INSTANCE_TYPE" required:"true"`
	GPUQuotaServiceCode string `envconfig:"GPU_QUOTA_SERVICE_CODE" required:"true"`
}

// LoadFromEnvVars parses the required configuration variables. Throws an error if the validations aren't met
func (c *ServiceConfiguration) LoadFromEnvVars() error {
	return envconfig.Process("", c)
}
