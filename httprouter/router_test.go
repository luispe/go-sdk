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
		name     string
		config   httprouter.Config
		inputURL string
		wantCode int
		wantMsg  string
	}{
		{
			name:     "liveness",
			inputURL: "/liveness",
			config: httprouter.Config{
				HealthCheckLivenessHandler: http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
					w.WriteHeader(http.StatusNoContent)
				}),
			},
			wantCode: http.StatusNoContent,
			wantMsg:  "",
		},
		{
			name:     "readiness",
			inputURL: "/readiness",
			config: httprouter.Config{
				HealthCheckReadinessHandler: http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
					w.WriteHeader(http.StatusNoContent)
				}),
			},
			wantCode: http.StatusNoContent,
			wantMsg:  "",
		},
		{
			name:     "not found handler",
			inputURL: "/notfound",
			config: httprouter.Config{
				NotFoundHandler: http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
					w.WriteHeader(http.StatusNotFound)
					writeLen, err := w.Write([]byte("handler not found"))
					assert.NoError(t, err)
					assert.EqualValues(t, 17, writeLen)
				}),
			},
			wantCode: http.StatusNotFound,
			wantMsg:  "handler not found",
		},
		{
			name:     "profiler is not active",
			inputURL: "/debug",
			wantCode: http.StatusNotFound,
			wantMsg:  "404 page not found\n",
		},
		{
			name:     "profiler is active",
			config:   httprouter.Config{EnableProfiler: true},
			inputURL: "/debug",
			wantCode: http.StatusOK,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			r := httprouter.New(test.config)

			server := httptest.NewServer(r)
			defer server.Close()

			res, err := http.Get(fmt.Sprintf("%s%s", server.URL, test.inputURL))
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, test.wantCode, res.StatusCode)

			if test.wantMsg != "" {
				b, err := io.ReadAll(res.Body)
				assert.Nil(t, err)
				assert.Equal(t, test.wantMsg, string(b))
			}
		})
	}
}

func TestRouterMethod(t *testing.T) {
	tests := []struct {
		name     string
		shortcut func(r *httprouter.Router, path string, handler httprouter.Handler, mw ...httprouter.Middleware)
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
			r := httprouter.New(httprouter.Config{})

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

func TestRouterAny(t *testing.T) {
	r := httprouter.New(httprouter.Config{})

	methods := []string{
		http.MethodGet,
		http.MethodHead,
		http.MethodOptions,
		http.MethodPost,
		http.MethodPut,
		http.MethodPatch,
		http.MethodDelete,
		http.MethodTrace,
	}

	server := httptest.NewServer(r)
	defer server.Close()

	h := func(w http.ResponseWriter, r *http.Request) error {
		w.WriteHeader(http.StatusOK)
		return nil
	}

	r.Any("/", h)

	for _, method := range methods {
		req, err := http.NewRequest(method, server.URL, nil)
		if err != nil {
			t.Fatal(err)
		}

		client := http.Client{}
		resp, err := client.Do(req)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	}
}

func TestRouterRoutes(t *testing.T) {
	r := httprouter.New(httprouter.Config{})
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

func TestRouteGroupMethod(t *testing.T) {
	tests := []struct {
		name     string
		shortcut func(r *httprouter.RouteGroup, pattern string, handler httprouter.Handler, mw ...httprouter.Middleware)
		method   string
	}{
		{
			name:     "get",
			shortcut: (*httprouter.RouteGroup).Get,
			method:   http.MethodGet,
		},
		{
			name:     "head",
			shortcut: (*httprouter.RouteGroup).Head,
			method:   http.MethodHead,
		},
		{
			name:     "options",
			shortcut: (*httprouter.RouteGroup).Options,
			method:   http.MethodOptions,
		},
		{
			name:     "post",
			shortcut: (*httprouter.RouteGroup).Post,
			method:   http.MethodPost,
		},
		{
			name:     "put",
			shortcut: (*httprouter.RouteGroup).Put,
			method:   http.MethodPut,
		},
		{
			name:     "patch",
			shortcut: (*httprouter.RouteGroup).Patch,
			method:   http.MethodPatch,
		},
		{
			name:     "delete",
			shortcut: (*httprouter.RouteGroup).Delete,
			method:   http.MethodDelete,
		},
		{
			name:     "trace",
			shortcut: (*httprouter.RouteGroup).Trace,
			method:   http.MethodTrace,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			r := httprouter.New(httprouter.Config{})

			server := httptest.NewServer(r)
			defer server.Close()

			test.shortcut(r.Group("/group"), "/", func(w http.ResponseWriter, r *http.Request) error {
				w.WriteHeader(http.StatusOK)
				return nil
			})

			req, err := http.NewRequest(test.method, fmt.Sprintf("%s/group", server.URL), nil)
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

func TestRouterHandlerReturnNoError(t *testing.T) {
	var mwWasCalled bool
	mw := func(f http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			f(w, r)
			mwWasCalled = true
		}
	}

	r := httprouter.New(httprouter.Config{Mw: []httprouter.Middleware{mw}})

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
		mw := func(f http.HandlerFunc) http.HandlerFunc {
			return func(w http.ResponseWriter, r *http.Request) {
				f(w, r)
				mwWasCalled = true
			}
		}

		t.Run(tt.name, func(t *testing.T) {
			config := httprouter.Config{
				Mw: []httprouter.Middleware{mw},
			}

			if tt.errorHandlerFunc != nil {
				config.ErrorHandlerFunc = tt.errorHandlerFunc
			}

			router := httprouter.New(config)
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

func TestRouteGroupAny(t *testing.T) {
	r := httprouter.New(httprouter.Config{})
	g := r.Group("/group")

	// http methods
	methods := []string{
		http.MethodGet,
		http.MethodHead,
		http.MethodOptions,
		http.MethodPost,
		http.MethodPut,
		http.MethodPatch,
		http.MethodDelete,
		http.MethodTrace,
	}

	server := httptest.NewServer(r)
	defer server.Close()

	h := func(w http.ResponseWriter, r *http.Request) error {
		w.WriteHeader(http.StatusOK)
		return nil
	}

	g.Any("/", h)

	for _, method := range methods {
		req, err := http.NewRequest(method, fmt.Sprintf("%s/group", server.URL), nil)
		if err != nil {
			t.Fatal(err)
		}

		client := http.Client{}
		resp, err := client.Do(req)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	}
}

func TestRouteGroupMultiple(t *testing.T) {
	tests := []struct {
		name         string
		path         string
		expectedLogs []string
	}{
		{
			name:         "get request to /ping",
			path:         "/test",
			expectedLogs: []string{"global middleware", "test from /"},
		},
		{
			name:         "get request to /groupA/test",
			path:         "/groupA/test",
			expectedLogs: []string{"global middleware", "groupA middleware", "test from /groupA/test"},
		},
		{
			name:         "get request to /groupB/test",
			path:         "/groupB/test",
			expectedLogs: []string{"global middleware", "test from /groupB/test"},
		},
		{
			name:         "get request to /groupB/groupBB/test",
			path:         "/groupB/groupBB/test",
			expectedLogs: []string{"global middleware", "groupBB middleware", "test from /groupB/groupBB/test"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var logs []string
			mwLog := func(s string) httprouter.Middleware {
				return func(handler http.HandlerFunc) http.HandlerFunc {
					return func(w http.ResponseWriter, r *http.Request) {
						logs = append(logs, s)
						handler(w, r)
					}
				}
			}

			hLog := func(s string) httprouter.Handler {
				return func(w http.ResponseWriter, r *http.Request) error {
					w.WriteHeader(200)
					logs = append(logs, s)
					return nil
				}
			}

			/*	Register the following groups' topology:

				/					=> with global middleware
				/groupA/ 			=> with groupA middleware (inherit global mid)
				/groupB/ 			=> with no middlewares (inherit global mid)
				/groupB/groupBB/ 	=> with groupBB middleware (inherit global mid and group BB)
			*/

			// Registers root
			r := httprouter.New(httprouter.Config{Mw: []httprouter.Middleware{mwLog("global middleware")}})
			r.Get("/test", hLog("test from /"))

			// Registers groupA with groupA middleware
			groupA := r.Group("/groupA", mwLog("groupA middleware"))
			groupA.Get("/test", hLog("test from /groupA/test"))

			// Registers groupB with no middleware
			groupB := r.Group("/groupB")
			groupB.Get("/test", hLog("test from /groupB/test"))

			// Registers groupBB with groupBB middleware
			groupBB := groupB.Group("/groupBB", mwLog("groupBB middleware"))
			groupBB.Get("/test", hLog("test from /groupB/groupBB/test"))

			server := httptest.NewServer(r)
			defer server.Close()

			req, err := http.NewRequest("GET", fmt.Sprintf("%s%s", server.URL, test.path), nil)
			if err != nil {
				t.Fatal(err)
			}

			client := http.Client{}
			resp, err := client.Do(req)
			assert.Nil(t, err)
			assert.Equal(t, 200, resp.StatusCode)
			assert.Equal(t, test.expectedLogs, logs)
		})
	}
}
