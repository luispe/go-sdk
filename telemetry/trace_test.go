package telemetry_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
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
