package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/pomelo-la/go-toolkit/httprouter"
	"github.com/pomelo-la/go-toolkit/httprouter/middleware"
)

type breaker struct{ mock.Mock }

func (b *breaker) Allow() bool {
	args := b.Called()
	return args.Bool(0)
}
func (b *breaker) Success() { b.Called() }
func (b *breaker) Failure() { b.Called() }

func TestMidBreaker(t *testing.T) {
	tests := []struct {
		Name          string
		HandlerStatus int
		SetupMock     func(b *breaker)
		AssertFunc    func(t *testing.T, status int)
	}{
		{
			Name:          "Success on handler <500",
			HandlerStatus: http.StatusTeapot,
			SetupMock: func(b *breaker) {
				b.On("Allow").Return(true).Once()
				b.On("Success").Once()
			},
			AssertFunc: func(t *testing.T, status int) {
				require.EqualValues(t, http.StatusTeapot, status)
			},
		},
		{
			Name:          "Failure on handler >=500",
			HandlerStatus: http.StatusInternalServerError,
			SetupMock: func(b *breaker) {
				b.On("Allow").Return(true).Once()
				b.On("Failure").Once()
			},
			AssertFunc: func(t *testing.T, status int) {
				require.EqualValues(t, http.StatusInternalServerError, status)
			},
		},
		{
			Name:          "Server Unavailable on open circuit",
			HandlerStatus: http.StatusOK,
			SetupMock: func(b *breaker) {
				b.On("Allow").Return(false).Once()
			},
			AssertFunc: func(t *testing.T, status int) {
				require.EqualValues(t, http.StatusServiceUnavailable, status)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			app := httprouter.New(httprouter.Config{})

			var cb breaker
			defer cb.AssertExpectations(t)
			tt.SetupMock(&cb)

			handler := func(w http.ResponseWriter, r *http.Request) error {
				w.WriteHeader(tt.HandlerStatus)
				return nil
			}
			mdl := middleware.Breaker(&cb, middleware.DefaultBreakerValidator)

			app.Method(http.MethodGet, "/", handler, mdl)

			recorder := httptest.NewRecorder()
			request, err := http.NewRequest("GET", "/", nil)
			assert.NoError(t, err)

			app.ServeHTTP(recorder, request)

			tt.AssertFunc(t, recorder.Code)
		})
	}
}
