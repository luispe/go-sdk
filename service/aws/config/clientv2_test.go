package config_test

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/stretchr/testify/assert"

	cfg "github.com/pomelo-la/go-toolkit/service/aws/config"
)

func TestNewWithoutConfig(t *testing.T) {
	type expected struct {
		config *aws.Config
		err    error
	}

	want := expected{config: &aws.Config{Region: "us-east-1"}, err: nil}

	t.Run("success config without options", func(t *testing.T) {
		got, err := cfg.New()
		assert.Equal(t, want.err, err)
		if got != nil {
			assert.Equal(t, want.config.Region, got.Region)
		}
	})
}

func TestNewWithConfig(t *testing.T) {
	type expected struct {
		config   *aws.Config
		endpoint *aws.Endpoint
		err      error
	}

	want := expected{
		config: &aws.Config{
			Region: "us-west-2",
		},
		endpoint: &aws.Endpoint{
			URL: "http://localhost:4566",
		},
		err: nil,
	}

	t.Run("success config with options", func(t *testing.T) {
		got, err := cfg.New(
			cfg.WithEndpoint("http://localhost:4566"),
			cfg.WithRegion("us-west-2"),
			cfg.WithProfile("some-profile"),
		)
		assert.Equal(t, want.err, err)
		if got != nil {
			assert.Equal(t, want.config.Region, got.Region)
		}
	})
}
