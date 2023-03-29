// Package metering provides objects to interact with metering API
package metering

import "github.com/kelseyhightower/envconfig"

// Configuration represents the configuration for marketplace client.
type Configuration struct {
	ProductCode string `envconfig:"PRODUCT_CODE" required:"true"`
}

// LoadFromEnvVars reads all env vars required for the marketplace client.
func (c *Configuration) LoadFromEnvVars() error {
	return envconfig.Process("", c)
}
