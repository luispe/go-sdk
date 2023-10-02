package config

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
)

// ErrorLoadSdkConfig indicates that the aws sdk unable to load.
var ErrorLoadSdkConfig = errors.New("unable to load aws SDK config")

// New instantiates an aws config with sane defaults.
func New(optFns ...func(opts *Config)) (*aws.Config, error) {
	var opts Config
	for _, fn := range optFns {
		fn(&opts)
	}

	if opts.Region == "" {
		opts.Region = _defaultAwsRegion
	}

	customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		if opts.Endpoint != "" {
			return aws.Endpoint{
				PartitionID:   "aws",
				URL:           opts.Endpoint,
				SigningRegion: opts.Region,
			}, nil
		}

		return aws.Endpoint{}, nil
	})

	cfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithRegion(opts.Region),
		config.WithEndpointResolverWithOptions(customResolver),
		config.WithSharedConfigProfile(opts.Profile),
	)
	if err != nil {
		return nil, ErrorLoadSdkConfig
	}

	return &cfg, nil
}
