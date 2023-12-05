package telemetry

import (
	"context"
	"os"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdkTrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

const _shutdownTraceTimeout = 10 * time.Second

// Trace represents a trace provider.
type Trace struct {
	trace *sdkTrace.TracerProvider
}

// NewTrace creates a new trace provider.
func NewTrace(ctx context.Context, serviceName string) (*Trace, error) {
	exp, err := otlptracegrpc.New(ctx)
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

	traceProvider := sdkTrace.NewTracerProvider(
		sdkTrace.WithSampler(sdkTrace.AlwaysSample()),
		sdkTrace.WithBatcher(exp),
		sdkTrace.WithResource(res),
	)

	otel.SetTracerProvider(traceProvider)
	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{},
			propagation.Baggage{},
		))

	return &Trace{trace: traceProvider}, nil
}

// ShutdownTraceProvider shuts down the TraceProvider gracefully.
func (tp Trace) ShutdownTraceProvider(ctx context.Context, optFns ...func(options *TraceOptions)) error {
	var opt TraceOptions
	for _, fn := range optFns {
		fn(&opt)
	}

	if opt.shutdownTimeout <= 0 {
		opt.shutdownTimeout = _shutdownTraceTimeout
	}

	ctxDeadline, cancel := context.WithTimeout(ctx, opt.shutdownTimeout)
	defer cancel()

	return tp.trace.Shutdown(ctxDeadline)
}

// TraceOptions represents the options for the Trace functionality.
type TraceOptions struct {
	shutdownTimeout time.Duration
}

// WithTraceShutdown allows you to configure the shutdown (in seconds)
// that the shutdown trace provider.
func WithTraceShutdown(duration time.Duration) func(options *TraceOptions) {
	return func(opt *TraceOptions) {
		opt.shutdownTimeout = duration
	}
}
