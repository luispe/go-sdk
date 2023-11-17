package telemetry

import (
	"context"
	"os"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

const (
	_metricInterval        = 60 * time.Second
	_shutdownMetricTimeout = 5 * time.Second
)

// Metric represents a metric provider.
type Metric struct {
	meter *sdkmetric.MeterProvider
}

func newRelicTemporalitySelector(kind sdkmetric.InstrumentKind) metricdata.Temporality {
	if kind == sdkmetric.InstrumentKindUpDownCounter || kind == sdkmetric.InstrumentKindObservableUpDownCounter {
		return metricdata.CumulativeTemporality
	}
	return metricdata.DeltaTemporality
}

// NewMetric creates a new metric provider.
func NewMetric(ctx context.Context, serviceName string, optFns ...func(options *MetricOptions)) (*Metric, error) {
	var opt MetricOptions
	for _, fn := range optFns {
		fn(&opt)
	}

	if opt.interval <= 0 {
		opt.interval = _metricInterval
	}

	exp, err := otlpmetricgrpc.New(ctx, otlpmetricgrpc.WithTemporalitySelector(newRelicTemporalitySelector))
	if err != nil {
		return nil, err
	}

	hostname, err := os.Hostname()
	if err != nil {
		return nil, err
	}
	res := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceInstanceIDKey.String(hostname),
		semconv.ServiceName(serviceName),
	)

	metricProvider := sdkmetric.NewMeterProvider(
		sdkmetric.WithResource(res),
		sdkmetric.WithReader(
			sdkmetric.NewPeriodicReader(
				exp,
				sdkmetric.WithInterval(opt.interval),
			)))
	otel.SetMeterProvider(metricProvider)

	return &Metric{meter: metricProvider}, err
}

// ShutdownMetricProvider shuts down the MetricProvider gracefully.
func (mp Metric) ShutdownMetricProvider(ctx context.Context, optFns ...func(options *MetricOptions)) error {
	var opt MetricOptions
	for _, fn := range optFns {
		fn(&opt)
	}

	if opt.shutdownTimeout <= 0 {
		opt.shutdownTimeout = _shutdownMetricTimeout
	}
	ctx, cancel := context.WithTimeout(ctx, opt.shutdownTimeout)
	defer cancel()

	return mp.meter.Shutdown(ctx)
}

// MetricOptions represents the options for the Metric functionality.
type MetricOptions struct {
	interval        time.Duration
	shutdownTimeout time.Duration
}

// WithMetricInterval configures the intervening time between exports for a PeriodicReader.
// If this option is not used or d is less than or equal to zero, 60 seconds is used as the default.
func WithMetricInterval(interval time.Duration) func(options *MetricOptions) {
	return func(opt *MetricOptions) {
		opt.interval = interval
	}
}

// WithMetricShutdown allows you to configure the shutdown (in seconds)
// that the shutdown metric provider.
func WithMetricShutdown(duration time.Duration) func(options *MetricOptions) {
	return func(opt *MetricOptions) {
		opt.shutdownTimeout = duration
	}
}
