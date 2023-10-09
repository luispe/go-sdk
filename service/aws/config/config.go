package config

var _defaultAwsRegion = "us-east-1"

// Config holds the configuration options for aws config.
type Config struct {
	// Endpoint specifies the endpoint override.
	Endpoint string
	// Region specifies the aws region.
	Region string
}

// WithEndpoint allows you to configure the endpoint.
func WithEndpoint(endpointOverride string) func(*Config) {
	return func(config *Config) {
		config.Endpoint = endpointOverride
	}
}

// WithRegion allows you to configure the different region for the aws client config.
//
// Default behavior is us-east-1.
func WithRegion(region string) func(*Config) {
	return func(config *Config) {
		config.Region = region
	}
}
