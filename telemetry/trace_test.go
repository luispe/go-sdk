package telemetry_test

import (
	"context"
	"encoding/hex"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

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

func TestGetTraceID(t *testing.T) {
	t.Run("Test GetTraceID without existing trace ID", func(t *testing.T) {
		result, err := telemetry.GetTraceID(context.Background())
		if err != nil {
			t.Errorf("Expected no error, but got %v", err)
		}

		// Check if the generated trace ID is a valid hexadecimal string
		_, err = hex.DecodeString(result)
		if err != nil {
			t.Errorf("Generated trace ID is not a valid hexadecimal string: %v", err)
		}
	})
}
