package main

import (
	"github.com/kelseyhightower/envconfig"
)

// Configuration defines all env vars required for the application
type Configuration struct {
	AWSRegion string `envconfig:"REGION" required:"true"`
}

// LoadEnvVars reads all env vars required for the server package
func (c *Configuration) LoadFromEnvVars() error {
	return envconfig.Process("", c)
}
