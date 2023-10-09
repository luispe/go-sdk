package config

import (
	awsv1 "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

// NewV1 instantiates an aws config with sane defaults.
func NewV1(optFns ...func(config *Config)) (*session.Session, error) {
	var opts Config
	for _, fn := range optFns {
		fn(&opts)
	}

	if opts.Region == "" {
		opts.Region = _defaultAwsRegion
	}

	// If awsConfig.Endpoint is empty
	// the user wish use default pkg config
	if opts.Endpoint == "" {
		awsSess, err := session.NewSession(&awsv1.Config{
			Region: awsv1.String(opts.Region),
		})
		if err != nil {
			return nil, ErrorLoadSdkConfig
		}

		return awsSess, nil
	}

	awsSess, err := session.NewSessionWithOptions(
		session.Options{
			Config: awsv1.Config{
				Region:   awsv1.String(opts.Region),
				Endpoint: awsv1.String(opts.Endpoint),
			},
		},
	)
	if err != nil {
		return nil, ErrorLoadSdkConfig
	}

	return awsSess, nil
}
