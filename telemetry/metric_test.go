package telemetry_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/pomelo-la/go-toolkit/telemetry"
)

func TestNewMetric(t *testing.T) {
	ctx := context.Background()

	metric, err := telemetry.NewMetric(ctx, "my-service-name")
	assert.NotNil(t, metric)
	assert.NoError(t, err)
}

func TestNewMetricWithOptions(t *testing.T) {
	ctx := context.Background()

	metric, err := telemetry.NewMetric(ctx,
		"my-service-name",
		telemetry.WithMetricInterval(3*time.Second),
	)
	assert.NotNil(t, metric)
	assert.NoError(t, err)
}
