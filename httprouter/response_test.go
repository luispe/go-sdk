package httprouter_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/pomelo-la/go-toolkit/httprouter"
)

func TestRespondJSONNoContent(t *testing.T) {
	w := httptest.NewRecorder()

	err := httprouter.RespondJSON(w, http.StatusNoContent, nil)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, w.Code)
	assert.Empty(t, w.Body)
}

func TestRespondJSONNilBody(t *testing.T) {
	w := httptest.NewRecorder()

	err := httprouter.RespondJSON(w, http.StatusOK, nil)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Empty(t, w.Body)
}

func TestRespondJSON(t *testing.T) {
	tests := []struct {
		name     string
		input    any
		expected string
	}{
		{
			name:     "it should respond using a map",
			input:    make(map[string]interface{}),
			expected: "{}",
		},
		{
			name:     "it should respond using a map",
			input:    []byte(`{}`),
			expected: "{}",
		},
		{
			name:     "it should respond using a map",
			input:    strings.NewReader("{}"),
			expected: "{}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()

			err := httprouter.RespondJSON(w, http.StatusOK, tt.input)

			assert.NoError(t, err)
			assert.Equal(t, "application/json", w.Header().Get("Content-type"))
			assert.Equal(t, tt.expected, w.Body.String())
			assert.Equal(t, 200, w.Code)
		})
	}
}
