package server

import (
	"github.com/ydataai/aws-quota-provider/pkg/common"
)

// Configuration defines a struct with required environment variables for server
type Configuration struct {
	Host string
	Port string
}

// LoadEnvVars reads all env vars required for the server package
func (c *Configuration) LoadEnvVars() error {
	host, _ := common.VariableFromEnvironment("HOST")

	port, err := common.VariableFromEnvironment("PORT")
	if err != nil {
		return err
	}

	c.Host = host
	c.Port = port

	return nil
}
