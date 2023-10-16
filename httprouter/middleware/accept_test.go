package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/pomelo-la/go-toolkit/httprouter"
	"github.com/pomelo-la/go-toolkit/httprouter/middleware"
)

func TestAcceptHeaderMid(t *testing.T) {
	app := httprouter.New(httprouter.Config{})
	app.Method(http.MethodGet, "/",
		func(w http.ResponseWriter, r *http.Request) error {
			return nil
		}, middleware.Accept("^.+/json", "image/*"))

	tests := []struct {
		mt         string
		statusCode int
	}{
		{
			mt:         "application/json",
			statusCode: http.StatusOK,
		},
		{
			mt:         "application/json+aws",
			statusCode: http.StatusOK,
		},
		{
			mt:         "image/jpeg",
			statusCode: http.StatusOK,
		},
		{
			mt:         "application/json;q=0.9",
			statusCode: http.StatusOK,
		},
		{
			mt:         "application/xml",
			statusCode: http.StatusNotAcceptable,
		},
		{
			mt:         "*/*",
			statusCode: http.StatusOK,
		},
		{
			mt:         "",
			statusCode: http.StatusOK,
		},
	}

	for _, tt := range tests {
		recorder := httptest.NewRecorder()
		request, err := http.NewRequest("GET", "/", nil)
		assert.NoError(t, err)
		request.Header.Add("Accept", tt.mt)

		app.ServeHTTP(recorder, request)

		assert.EqualValues(t, tt.statusCode, recorder.Result().StatusCode)
	}
}

func TestMatchJSONAcceptHeader(t *testing.T) {
	app := httprouter.New(httprouter.Config{})
	app.Method(http.MethodGet, "/",
		func(w http.ResponseWriter, r *http.Request) error {
			return nil
		}, middleware.AcceptJSON())

	tests := []struct {
		mt         string
		statusCode int
	}{
		{
			mt:         "application/json",
			statusCode: http.StatusOK,
		},
		{
			mt:         "application/json+aws",
			statusCode: http.StatusNotAcceptable,
		},
	}

	for _, tt := range tests {
		recorder := httptest.NewRecorder()
		request, err := http.NewRequest("GET", "/", nil)
		assert.NoError(t, err)
		request.Header.Add("Accept", tt.mt)

		app.ServeHTTP(recorder, request)

		assert.EqualValues(t, tt.statusCode, recorder.Result().StatusCode)
	}
}
