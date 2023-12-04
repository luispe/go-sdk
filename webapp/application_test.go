package webapp_test

import (
	"net"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/pomelo-la/go-toolkit/httprouter"
	"github.com/pomelo-la/go-toolkit/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/pomelo-la/go-toolkit/webapp"
)

func TestNewWebApplication(t *testing.T) {
	t.Run("default app", func(t *testing.T) {
		app, err := webapp.New("test-app")
		require.NoError(t, err)
		require.NotNil(t, app)
		require.NotNil(t, app.Logger)
		require.NotNil(t, app.Router)
		require.NotNil(t, app.Tracer)
		require.Equal(t, "local", app.Environment.Name)
	})

	t.Run("err app name empty", func(t *testing.T) {
		_, err := webapp.New("")
		assert.Equal(t, err, webapp.ErrInvalidAppName)
	})

	t.Run("err app name blank space", func(t *testing.T) {
		_, err := webapp.New(" ")
		assert.Equal(t, err, webapp.ErrInvalidAppName)
	})

	t.Run("err app name blank space between name", func(t *testing.T) {
		_, err := webapp.New("my app")
		assert.Equal(t, err, webapp.ErrInvalidAppName)
	})

	t.Run("web app with configure log level", func(t *testing.T) {
		app, err := webapp.New("test-app", webapp.WithLogLevel(logger.LevelWarn))
		require.NoError(t, err)
		require.NotNil(t, app)
		require.NotNil(t, app.Logger)
		require.NotNil(t, app.Router)
		require.NotNil(t, app.Tracer)
		require.Equal(t, "local", app.Environment.Name)
	})

	t.Run("web app with configure log level from env", func(t *testing.T) {
		err := os.Setenv("LOG_LEVEL", "ERROR")
		assert.NoError(t, err)

		app, err := webapp.New("test-app")
		require.NoError(t, err)
		require.NotNil(t, app)
		require.NotNil(t, app.Logger)
		require.NotNil(t, app.Router)
		require.NotNil(t, app.Tracer)
		require.Equal(t, "local", app.Environment.Name)
	})

	t.Run("web app with configure timeouts", func(t *testing.T) {
		timeOuts := httprouter.Timeouts{
			ReadTimeout:       5 * time.Second,
			ReadHeaderTimeout: 5 * time.Second,
			WriteTimeout:      10 * time.Second,
			ShutdownTimeout:   10 * time.Second,
		}
		app, err := webapp.New("test-app", webapp.WithTimeouts(timeOuts))
		require.NoError(t, err)
		require.NotNil(t, app)
		require.NotNil(t, app.Logger)
		require.NotNil(t, app.Router)
		require.NotNil(t, app.Tracer)
		require.Equal(t, "local", app.Environment.Name)
	})

	t.Run("web app with configure listener", func(t *testing.T) {
		ln, err := net.Listen("tcp", ":9090")
		require.NoError(t, err)
		app, err := webapp.New("test-app", webapp.WithListener(ln))
		require.NoError(t, err)
		require.NotNil(t, app)
		require.NotNil(t, app.Logger)
		require.NotNil(t, app.Router)
		require.NotNil(t, app.Tracer)
		require.Equal(t, "local", app.Environment.Name)
	})

	t.Run("web app with configure err handler func", func(t *testing.T) {
		errHandler := func(err error, defaultHandlerError func(error) httprouter.HandlerError) httprouter.HandlerError {
			return httprouter.HandlerError{
				Error:      err,
				Notify:     false,
				StatusCode: http.StatusInternalServerError,
			}
		}
		app, err := webapp.New("test-app", webapp.WithErrorHandler(errHandler))
		require.NoError(t, err)
		require.NotNil(t, app)
		require.NotNil(t, app.Logger)
		require.NotNil(t, app.Router)
		require.NotNil(t, app.Tracer)
		require.Equal(t, "local", app.Environment.Name)
	})

	t.Run("web app with configure environment", func(t *testing.T) {
		app, err := webapp.New("test-app", webapp.WithEnvironment("production"))
		require.NoError(t, err)
		require.NotNil(t, app)
		require.NotNil(t, app.Logger)
		require.NotNil(t, app.Router)
		require.NotNil(t, app.Tracer)
		require.Equal(t, "production", app.Environment.Name)
	})

	mw := func(f http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			f.ServeHTTP(w, r)
		})
	}

	t.Run("web app with configure middleware", func(t *testing.T) {
		app, err := webapp.New("test-app", webapp.WithGlobalMiddlewares(mw))
		require.NoError(t, err)
		require.NotNil(t, app)
		require.NotNil(t, app.Logger)
		require.NotNil(t, app.Router)
		require.NotNil(t, app.Tracer)
		require.Equal(t, "local", app.Environment.Name)
	})
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

			app, err := webapp.New("test-app")
			require.NoError(t, err)
			require.NotEmpty(t, app)

			err = app.Run()
			require.Equal(t, tc.wantErr, err)
		})
	}
}
