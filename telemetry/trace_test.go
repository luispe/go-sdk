package telemetry_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/attribute"

	"github.com/pomelo-la/go-toolkit/telemetry"
)

func TestNewTrace(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	type expected struct {
		trace *telemetry.Trace
		err   error
	}

	mockTrace, _ := telemetry.NewTrace(context.Background())
	tests := []struct {
		name string
		args args
		want expected
	}{
		{
			name: "new client",
			args: args{ctx: context.Background()},
			want: expected{trace: mockTrace, err: nil},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := telemetry.NewTrace(tt.args.ctx)
			assert.ObjectsAreEqual(tt.want.trace, got)
			assert.Equal(t, tt.want.err, err)
		})
	}
}

func TestShutdownTimeout(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	type expected struct {
		trace *telemetry.Trace
		err   error
	}

	mockTrace, _ := telemetry.NewTrace(context.Background())
	tests := []struct {
		name string
		args args
		want expected
	}{
		{
			name: "new client",
			args: args{ctx: context.Background()},
			want: expected{trace: mockTrace, err: nil},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := telemetry.NewTrace(tt.args.ctx)
			assert.NoError(t, err)

			err = client.ShutdownTraceProvider(tt.args.ctx)
			assert.Equal(t, tt.want.err, err)
		})
	}
}

func TestConfigShutdownTimeout(t *testing.T) {
	type args struct {
		ctx     context.Context
		timeout time.Duration
	}
	type expected struct {
		trace *telemetry.Trace
		err   error
	}

	mockTrace, _ := telemetry.NewTrace(context.Background())
	tests := []struct {
		name string
		args args
		want expected
	}{
		{
			name: "new client",
			args: args{ctx: context.Background(), timeout: 5 * time.Second},
			want: expected{trace: mockTrace, err: nil},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := telemetry.NewTrace(tt.args.ctx)
			assert.NoError(t, err)

			err = client.ShutdownTraceProvider(tt.args.ctx, telemetry.WithTraceShutdown(tt.args.timeout))
			assert.Equal(t, tt.want.err, err)
		})
	}
}

func TestAddSpan(t *testing.T) {
	ctx := context.Background()
	trace, err := telemetry.NewTrace(ctx)
	assert.NoError(t, err)

	_, span := trace.AddSpan(ctx, "go-toolkit.telemetry", attribute.String("key1", "value1"))
	if span == nil {
		t.Error("Expected a valid span, but got nil")
	}
}

func TestSanitizeMetricTagValue(t *testing.T) {
	tt := []struct {
		tag  string
		want string
	}{
		{
			tag:  "",
			want: "",
		},
		{
			tag:  "/fixed/{first_param}",
			want: "/fixed/_first_param",
		},
		{
			tag:  "/{first_param}/{second_param}/fixed",
			want: "/_first_param/_second_param/fixed",
		},
		{
			tag:  "/fixed/fixed",
			want: "/fixed/fixed",
		},
		{
			tag:  "/{first_param:[a-z0-9]+}/fixed",
			want: "/_first_param:[a-z0-9]+/fixed",
		},
		{
			tag:  "/fixed/fixed/",
			want: "/fixed/fixed",
		},
	}

	for _, tc := range tt {
		t.Run(tc.tag, func(t *testing.T) {
			name := telemetry.SanitizeMetricTagValue(tc.tag)
			require.Equal(t, tc.want, name)
		})
	}
}

func TestSanitizeMetricTagValue_MultipleTimes(t *testing.T) {
	input := telemetry.SanitizeMetricTagValue("/{first_param}/{second_param}/fixed")
	for i := 0; i < 10; i++ {
		tmp := telemetry.SanitizeMetricTagValue(input)
		require.Equal(t, input, tmp)
		input = tmp
	}
}
