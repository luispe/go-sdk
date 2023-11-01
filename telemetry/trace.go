package telemetry

import (
	"context"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdkTrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

const _shutdownTraceTimeout = 10 * time.Second

type ctxKey int

const key ctxKey = 1

// Values represent state for each request.
type Values struct {
	TraceID    string
	Tracer     trace.Tracer
	Now        time.Time
	StatusCode int
}

// Trace represents a trace provider.
type Trace struct {
	trace *sdkTrace.TracerProvider
}

// NewTrace creates a new trace provider.
func NewTrace(ctx context.Context) (*Trace, error) {
	exp, err := otlptracegrpc.New(ctx)
	if err != nil {
		return nil, err
	}

	res := resource.Default()

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

// AddSpan adds a OpenTelemetry span to the trace and context.
func (tp Trace) AddSpan(ctx context.Context, spanName string, keyValues ...attribute.KeyValue) (context.Context, trace.Span) {
	v, ok := ctx.Value(key).(*Values)
	if !ok || v.Tracer == nil {
		return ctx, trace.SpanFromContext(ctx)
	}

	ctx, span := v.Tracer.Start(ctx, spanName)
	for _, kv := range keyValues {
		span.SetAttributes(kv)
	}

	return ctx, span
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
