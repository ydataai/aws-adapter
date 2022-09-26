package marketplace

import "github.com/kelseyhightower/envconfig"

// AWSMarketplaceConfiguration represents the configuration for marketplace client.
type AWSMarketplaceConfiguration struct {
	ProductCode string `envconfig:"AWS_PRODUCT_CODE" required:"true"`
	Region      string `envconfig:"AWS_REGION" required:"true"`
}

// LoadFromEnvVars reads all env vars required for the marketplace client.
func (c *AWSMarketplaceConfiguration) LoadFromEnvVars() error {
	return envconfig.Process("", c)
}
