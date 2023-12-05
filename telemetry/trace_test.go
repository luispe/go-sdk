package telemetry_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

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

	mockTrace, _ := telemetry.NewTrace(context.Background(), "my-service-name")
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
			got, err := telemetry.NewTrace(tt.args.ctx, "my-service-name")
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

	mockTrace, _ := telemetry.NewTrace(context.Background(), "my-service-name")
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
			client, err := telemetry.NewTrace(tt.args.ctx, "my-service-name")
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

	mockTrace, _ := telemetry.NewTrace(context.Background(), "my-service-name")
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
			client, err := telemetry.NewTrace(tt.args.ctx, "my-service-name")
			assert.NoError(t, err)

			err = client.ShutdownTraceProvider(tt.args.ctx, telemetry.WithTraceShutdown(tt.args.timeout))
			assert.Equal(t, tt.want.err, err)
		})
	}
}
