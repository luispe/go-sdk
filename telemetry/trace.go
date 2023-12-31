package telemetry

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"os"
	"strings"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdkTrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"go.opentelemetry.io/otel/trace"
	"go.opentelemetry.io/otel/trace/noop"
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

var _patternReplacer = strings.NewReplacer(
	"{", "_",
	"}", "",
)

// SanitizeMetricTagValue sanitizes the given value in a standard way. It:
//   - Trims suffix "/".
//   - Replace "{" with "_"
//   - Remove  "}".
func SanitizeMetricTagValue(value string) string {
	value = strings.TrimRight(value, "/")
	return _patternReplacer.Replace(value)
}

// GetTraceID returns the trace id from the context.
func GetTraceID(ctx context.Context) (string, error) {
	v, ok := ctx.Value(key).(*Values)
	if ok {
		return v.TraceID, nil
	}

	id, err := generateTraceID()
	if err != nil {
		return "", err
	}

	traceID, err := trace.TraceIDFromHex(id)
	if err != nil {
		return "", err
	}

	return traceID.String(), nil
}

// GetValues returns the values from the context.
func GetValues(ctx context.Context) (*Values, error) {
	v, ok := ctx.Value(key).(*Values)
	if ok {
		return v, nil
	}

	id, err := generateTraceID()
	if err != nil {
		return nil, err
	}

	traceID, err := trace.TraceIDFromHex(id)
	if err != nil {
		return nil, err
	}

	return &Values{
		TraceID: traceID.String(),
		Tracer:  noop.NewTracerProvider().Tracer(""),
		Now:     time.Now(),
	}, nil
}

func generateTraceID() (string, error) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
