package service

import "github.com/ydataai/aws-quota-provider/pkg/common"

// RESTServiceConfiguration defines required configuration for rest service
type RESTServiceConfiguration struct {
	gpuInstanceType    string
	gpuCodeServiceCode string
}

// LoadEnvVars parses the required configuration variables. Throws an error if the validations aren't met
func (c *RESTServiceConfiguration) LoadEnvVars() error {
	gpuInstanceType, err := common.VariableFromEnvironment("GPU_INSTANCE_TYPE")
	if err != nil {
		return err
	}

	gpuCodeServiceCode, err := common.VariableFromEnvironment("GPU_QUOTA_SERVICE_CODE")
	if err != nil {
		return err
	}

	c.gpuInstanceType = gpuInstanceType
	c.gpuCodeServiceCode = gpuCodeServiceCode

	return nil
}
