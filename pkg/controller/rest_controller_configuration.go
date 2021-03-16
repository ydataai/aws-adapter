package controller

import (
	"time"

	"github.com/ydataai/aws-quota-provider/pkg/common"
)

// RESTControllerConfiguration defines a struct with required environment variables for rest controller
type RESTControllerConfiguration struct {
	timeout time.Duration
}

// LoadEnvVars reads all env vars required for the server package
func (c *RESTControllerConfiguration) LoadEnvVars() error {
	timeout, err := common.IntVariableFromEnvironment("REQUEST_TIMEOUT")
	if err != nil {
		return err
	}
	c.timeout = time.Minute * time.Duration(timeout)

	return nil
}
