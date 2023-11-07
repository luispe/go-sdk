package webapp_test

import (
	"net"
	"net/http"
	"os"
	"testing"

	"github.com/pomelo-la/go-toolkit/httprouter"
	"github.com/pomelo-la/go-toolkit/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/pomelo-la/go-toolkit/webapp"
)

func TestNewWebApplication(t *testing.T) {
	type args struct {
		env string
	}
	tests := []struct {
		name        string
		args        args
		err         error
		options     []func(options *webapp.AppOptions)
		wantRuntime webapp.Runtime
	}{
		{
			name: "New web backend with default local scope",
			args: args{env: "local"},
			options: []func(options *webapp.AppOptions){
				func(options *webapp.AppOptions) {
					webapp.WithLogLevel(logger.LevelInfo)
				},
			},
			wantRuntime: webapp.Runtime{
				Environment: "local",
			},
		},
		{
			name: "New web backend with an environment scope value",
			args: args{env: "local"},
			options: []func(options *webapp.AppOptions){
				func(options *webapp.AppOptions) {
					webapp.WithEnvironmentRuntime("local")
				},
			},
			wantRuntime: webapp.Runtime{
				Environment: "local",
			},
		},
		{
			name: "New web backend with an error handling function",
			args: args{env: "local"},
			options: []func(options *webapp.AppOptions){
				func(options *webapp.AppOptions) {
					webapp.WithErrorHandler(func(err error, defaultHandlerError func(error) httprouter.HandlerError) httprouter.HandlerError {
						return httprouter.HandlerError{
							Error:      err,
							Notify:     true,
							StatusCode: http.StatusInternalServerError,
						}
					})
				},
			},
			wantRuntime: webapp.Runtime{
				Environment: "local",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := os.Setenv("RUNTIME", tt.args.env)
			assert.NoError(t, err)

			app, err := webapp.New(tt.options...)
			require.Equal(t, tt.err, err)
			require.NotNil(t, app)
			require.NotNil(t, app.Logger)
			require.NotNil(t, app.Router)
			require.NotNil(t, app.Tracer)
			require.Equal(t, tt.wantRuntime.Environment, app.Runtime.Environment)
		})
	}
}

func TestApplicationRunError(t *testing.T) {
	tests := []struct {
		name    string
		port    string
		wantErr error
	}{
		{
			name: "invalid port number",
			port: "-9999",
			wantErr: &net.OpError{
				Op:  "listen",
				Net: "tcp",
				Err: &net.AddrError{
					Err:  "invalid port",
					Addr: "-9999",
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := os.Setenv("RUNTIME", "local")
			assert.NoError(t, err)

			originWebappPort := os.Getenv("PORT")
			err = os.Setenv("PORT", tc.port)
			assert.NoError(t, err)
			defer os.Setenv("PORT", originWebappPort)

			app, err := webapp.New()
			require.NoError(t, err)
			require.NotEmpty(t, app)

			err = app.Run()
			require.Equal(t, tc.wantErr, err)
		})
	}
}
