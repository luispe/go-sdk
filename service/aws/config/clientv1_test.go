package config_test

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/stretchr/testify/assert"

	cfg "github.com/pomelo-la/go-toolkit/service/aws/config"
)

func TestNewV1WithoutConfig(t *testing.T) {
	type expected struct {
		sess *session.Session
		err  error
	}

	want := expected{sess: &session.Session{Config: &aws.Config{Region: aws.String("us-east-1")}}, err: nil}

	t.Run("success config without options", func(t *testing.T) {
		got, err := cfg.NewV1()
		assert.Equal(t, want.err, err)
		if got != nil {
			assert.Equal(t, want.sess.Config.Region, got.Config.Region)
		}
	})
}

func TestNewV1WithConfig(t *testing.T) {
	type expected struct {
		sess *session.Options
		err  error
	}

	want := expected{
		sess: &session.Options{
			Config: aws.Config{
				Endpoint: aws.String("http://localhost:4566"),
				Region:   aws.String("us-west-2"),
			},
			Profile: "some-profile",
		},
		err: nil,
	}

	t.Run("success config with options", func(t *testing.T) {
		got, err := cfg.NewV1(
			cfg.WithEndpoint("http://localhost:4566"),
			cfg.WithRegion("us-west-2"),
			cfg.WithProfile("some-profile"),
		)
		assert.Equal(t, want.err, err)
		if got != nil {
			assert.Equal(t, want.sess.Config.Endpoint, got.Config.Endpoint)
			assert.Equal(t, want.sess.Config.Region, got.Config.Region)
		}
	})
}
