package marketplace

import "github.com/kelseyhightower/envconfig"

// AWSMarketplaceConfiguration represents the configuration for marketplace client.
type AWSMarketplaceConfiguration struct {
	RegistrationToken string `envconfig:"AWS_CUSTOMER_REGISTRATIO_TOKEN" required:"true"`
}

// LoadFromEnvVars reads all env vars required for the marketplace client.
func (c *AWSMarketplaceConfiguration) LoadFromEnvVars() error {
	return envconfig.Process("", c)
}
