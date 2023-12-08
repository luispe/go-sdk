package httprouter_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/pomelo-la/go-toolkit/httprouter"
)

func TestNewConfigHandlers(t *testing.T) {
	tests := []struct {
		name           string
		profilerActive bool
		inputURL       string
		wantCode       int
		wantMsg        string
	}{
		{
			name:     "liveness",
			inputURL: "/liveness",
			wantCode: http.StatusNoContent,
			wantMsg:  "",
		},
		{
			name:     "readiness",
			inputURL: "/readiness",
			wantCode: http.StatusNoContent,
			wantMsg:  "",
		},
		{
			name:     "not found handler",
			inputURL: "/notfound",
			wantCode: http.StatusNotFound,
			wantMsg:  "handler not found",
		},
		{
			name:     "profiler is not active",
			inputURL: "/debug",
			wantCode: http.StatusNotFound,
		},
		{
			name:           "profiler is active",
			profilerActive: true,
			inputURL:       "/debug",
			wantCode:       http.StatusOK,
		},
		{
			name:     "with middleware",
			inputURL: "/",
			wantCode: http.StatusNotFound,
			wantMsg:  "",
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			livenessHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNoContent)
			})
			readinessHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNoContent)
			})
			notFoundHandler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				w.WriteHeader(http.StatusNotFound)
				writeLen, err := w.Write([]byte("handler not found"))
				assert.NoError(t, err)
				assert.EqualValues(t, 17, writeLen)
			})

			mw := func(f http.Handler) http.Handler {
				return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					f.ServeHTTP(w, r)
				})
			}

			r := httprouter.New(
				httprouter.WithGlobalMiddlewares(mw),
				httprouter.WithHealthCheckLivenessHandler(livenessHandler),
				httprouter.WithHealthCheckReadinessHandler(readinessHandler),
				httprouter.WithNotFoundHandler(notFoundHandler),
				httprouter.WithEnableProfiler(tc.profilerActive),
			)

			server := httptest.NewServer(r)
			defer server.Close()

			res, err := http.Get(fmt.Sprintf("%s%s", server.URL, tc.inputURL))
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, tc.wantCode, res.StatusCode)

			if tc.wantMsg != "" {
				b, err := io.ReadAll(res.Body)
				assert.Nil(t, err)
				assert.Equal(t, tc.wantMsg, string(b))
			}
		})
	}
}

func TestRouterMethod(t *testing.T) {
	tests := []struct {
		name     string
		shortcut func(r *httprouter.Router, path string, handler httprouter.Handler)
		method   string
	}{
		{
			name:     "get",
			shortcut: (*httprouter.Router).Get,
			method:   http.MethodGet,
		},
		{
			name:     "head",
			shortcut: (*httprouter.Router).Head,
			method:   http.MethodHead,
		},
		{
			name:     "options",
			shortcut: (*httprouter.Router).Options,
			method:   http.MethodOptions,
		},
		{
			name:     "post",
			shortcut: (*httprouter.Router).Post,
			method:   http.MethodPost,
		},
		{
			name:     "put",
			shortcut: (*httprouter.Router).Put,
			method:   http.MethodPut,
		},
		{
			name:     "patch",
			shortcut: (*httprouter.Router).Patch,
			method:   http.MethodPatch,
		},
		{
			name:     "delete",
			shortcut: (*httprouter.Router).Delete,
			method:   http.MethodDelete,
		},
		{
			name:     "trace",
			shortcut: (*httprouter.Router).Trace,
			method:   http.MethodTrace,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			r := httprouter.New()

			server := httptest.NewServer(r)
			defer server.Close()

			test.shortcut(r, "/", func(w http.ResponseWriter, r *http.Request) error {
				w.WriteHeader(http.StatusOK)
				return nil
			})

			req, err := http.NewRequest(test.method, server.URL, nil)
			if err != nil {
				t.Fatal(err)
			}

			client := http.Client{}
			resp, err := client.Do(req)
			assert.Nil(t, err)
			assert.Equal(t, http.StatusOK, resp.StatusCode)
		})
	}
}

func TestRouterRoutes(t *testing.T) {
	r := httprouter.New()
	h := func(w http.ResponseWriter, r *http.Request) error {
		return nil
	}
	r.Get("/test", h)

	routes, err := r.Routes()

	assert.Nil(t, err)
	assert.Len(t, routes, 1)
	assert.Equal(t, "/test", routes[0].Route)
	assert.Equal(t, "GET", routes[0].Method)
	assert.Len(t, routes[0].Middlewares, 0)
	assert.NotNil(t, routes[0].Handler)
}

func TestRouterHandlerReturnNoError(t *testing.T) {
	var mwWasCalled bool
	mw := func(f http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			f.ServeHTTP(w, r)
			mwWasCalled = true
		})
	}

	r := httprouter.New()

	r.Use(mw)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) error {
		w.WriteHeader(http.StatusCreated)
		return nil
	})

	server := httptest.NewServer(r)
	defer server.Close()

	req, err := http.NewRequest("GET", server.URL, nil)
	if err != nil {
		t.Fatal(err)
	}

	client := http.Client{}
	resp, err := client.Do(req)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	assert.True(t, mwWasCalled)
}

// handlerSpy that captures the error returned by the handler so that we can validate
// that it is the same being notified.
type handlerSpy struct {
	h   httprouter.Handler
	err error
}

func (hs *handlerSpy) spy() httprouter.Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		hs.err = hs.h(w, r)
		return hs.err
	}
}

func TestRouterErrorHandler(t *testing.T) {
	tests := []struct {
		name               string
		errorHandlerFunc   httprouter.ErrorHandlerFunc
		handler            *handlerSpy
		expectedCode       int
		expectedErrMessage string
		shouldNotify       bool
	}{
		{
			name: "custom error with default error handler",
			handler: &handlerSpy{h: func(w http.ResponseWriter, r *http.Request) error {
				return errors.New("custom error")
			}},
			expectedCode:       http.StatusInternalServerError,
			expectedErrMessage: "custom error",
			shouldNotify:       true,
		},
		{
			name: "web error with default error handler",
			handler: &handlerSpy{h: func(w http.ResponseWriter, r *http.Request) error {
				return httprouter.NewErrorf(http.StatusBadRequest, "bad request")
			}},
			expectedCode:       http.StatusBadRequest,
			expectedErrMessage: "bad request",
			shouldNotify:       false,
		},
		{
			name: "custom error with custom error handler",
			errorHandlerFunc: func(err error, defaultHandlerError func(error) httprouter.HandlerError) httprouter.HandlerError {
				if err.Error() == "something went wrong" {
					e := httprouter.Error{
						Message: err.Error(),
					}
					return httprouter.HandlerError{
						StatusCode: http.StatusInternalServerError,
						Error:      e,
						Notify:     false,
					}
				}

				panic("should not be here")
			},
			handler: &handlerSpy{h: func(w http.ResponseWriter, r *http.Request) error {
				return errors.New("something went wrong")
			}},
			expectedCode:       http.StatusInternalServerError,
			expectedErrMessage: "something went wrong",
			shouldNotify:       false,
		},
		{
			name: "custom notifiable error with custom error handler",
			errorHandlerFunc: func(err error, defaultHandlerError func(error) httprouter.HandlerError) httprouter.HandlerError {
				if err.Error() == "something went wrong" {
					e := httprouter.Error{
						Message: err.Error(),
					}
					return httprouter.HandlerError{
						StatusCode: http.StatusInternalServerError,
						Error:      e,
						Notify:     true,
					}
				}

				panic("should not be here")
			},
			handler: &handlerSpy{h: func(w http.ResponseWriter, r *http.Request) error {
				return errors.New("something went wrong")
			}},
			expectedCode:       http.StatusInternalServerError,
			expectedErrMessage: "something went wrong",
			shouldNotify:       true,
		},
		{
			name: "custom notifiable error with custom error handler but using default callback",
			errorHandlerFunc: func(err error, defaultHandlerError func(error) httprouter.HandlerError) httprouter.HandlerError {
				return httprouter.DefaultHandlerError(err)
			},
			handler: &handlerSpy{h: func(w http.ResponseWriter, r *http.Request) error {
				return errors.New("something went wrong")
			}},
			expectedCode:       http.StatusInternalServerError,
			expectedErrMessage: "something went wrong",
			shouldNotify:       true,
		},
	}

	for _, tt := range tests {
		var mwWasCalled bool
		mw := func(f http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				f.ServeHTTP(w, r)
				mwWasCalled = true
			})
		}

		t.Run(tt.name, func(t *testing.T) {
			config := httprouter.Config{}

			if tt.errorHandlerFunc != nil {
				config.ErrorHandlerFunc = tt.errorHandlerFunc
			}

			router := httprouter.New()
			router.Use(mw)
			router.Get("/{id}", tt.handler.spy())

			rr := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/1", nil)

			ctx := httprouter.WithParams(req.Context(), map[string]string{"id": "1"})
			req = req.WithContext(ctx)

			router.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedCode, rr.Code)
			assert.True(t, mwWasCalled)

			var webErr httprouter.Error
			if err := json.Unmarshal(rr.Body.Bytes(), &webErr); err != nil {
				assert.NoError(t, err)
			}

			assert.Equal(t, tt.expectedErrMessage, webErr.Message)
		})

	}
}
