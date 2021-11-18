package service

import "github.com/kelseyhightower/envconfig"

// RESTServiceConfiguration defines required configuration for rest service
type RESTServiceConfiguration struct {
	GPUInstanceType     string `envconfig:"GPU_INSTANCE_TYPE" required:"true"`
	GPUQuotaServiceCode string `envconfig:"GPU_QUOTA_SERVICE_CODE" required:"true"`
}

// LoadFromEnvVars parses the required configuration variables. Throws an error if the validations aren't met
func (c *RESTServiceConfiguration) LoadFromEnvVars() error {
	return envconfig.Process("", c)
}
