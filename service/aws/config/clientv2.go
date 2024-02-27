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

	if opts.Endpoint != "" {
		cfg, err := buildLocalStackAwsConfig(opts)
		if err != nil {
			return nil, err
		}

		return cfg, nil
	}

	cfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithRegion(opts.Region),
	)
	if err != nil {
		return nil, ErrorLoadSdkConfig
	}

	return &cfg, nil
}

//revive:disable:unused-parameter
func buildLocalStackAwsConfig(opts Config) (*aws.Config, error) {
	customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			PartitionID:   "aws",
			URL:           opts.Endpoint,
			SigningRegion: opts.Region,
		}, nil
	})

	cfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithRegion(opts.Region),
		config.WithEndpointResolverWithOptions(customResolver),
	)
	if err != nil {
		return nil, ErrorLoadSdkConfig
	}

	return &cfg, nil
}

//revive:enable:unused-parameter
