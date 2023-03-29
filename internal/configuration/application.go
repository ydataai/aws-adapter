// Package configuration provides objects to configure adapter objects
package configuration

import (
	"github.com/kelseyhightower/envconfig"
)

// Application defines all env vars required for the application
type Application struct {
	Region string `envconfig:"REGION" required:"true"`
}

// LoadFromEnvVars reads all env vars required for the quota package
func (c *Application) LoadFromEnvVars() error {
	return envconfig.Process("", c)
}
