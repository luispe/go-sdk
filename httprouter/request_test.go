package httprouter_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/pomelo-la/go-toolkit/httprouter"
)

func TestParams(t *testing.T) {
	m := make(httprouter.URIParams)
	ctx := context.WithValue(context.Background(), httprouter.URIParamsCtxKey{}, m)
	r, _ := http.NewRequest("GET", "/", nil)
	assert.Equal(t, m, httprouter.Params(r.WithContext(ctx)))
}

func TestParamsNotFoundInContext(t *testing.T) {
	r, _ := http.NewRequest("GET", "/", nil)
	assert.Nil(t, httprouter.Params(r))
}

func TestParamsConversionFunctionsOK(t *testing.T) {
	params := httprouter.URIParams{
		"string": "go-toolkit rules!",
		"int":    "-1",
		"bool":   "true",
		"noint":  "wish I'd be a number",
		"nobool": "fuzzy",
		"uint":   "1",
		"nouint": "spread love anywhere you go",
	}

	type testParams struct {
		expectValue any
		assertFunc  func(string) (any, error)
		name        string
	}

	tt := map[string]testParams{
		"string": {
			expectValue: "go-toolkit rules!",
			assertFunc: func(p string) (any, error) {
				return params.String(p)
			},
			name: "convert to string",
		},
		"int": {
			expectValue: -1,
			assertFunc: func(p string) (any, error) {
				return params.Int(p)
			},
			name: "convert to int",
		},
		"uint": {
			expectValue: uint(1),
			assertFunc: func(p string) (any, error) {
				return params.Uint(p)
			},
			name: "convert to uint",
		},
		"bool": {
			expectValue: true,
			assertFunc: func(p string) (any, error) {
				return params.Bool(p)
			},
			name: "convert to bool",
		},
	}

	for k, tt := range tt {
		t.Run(tt.name, func(t *testing.T) {
			actual, err := tt.assertFunc(k)
			require.Nil(t, err)
			require.EqualValues(t, tt.expectValue, actual)
		})
	}
}

func TestParamsConversionFunctionsParamsNotFound(t *testing.T) {
	params := httprouter.URIParams{}

	type testParams struct {
		assertFunc func(string) error
		name       string
	}

	tests := []testParams{
		{
			assertFunc: func(p string) error {
				_, err := params.String(p)
				return err
			},
			name: "convert to string should fail",
		},
		{
			assertFunc: func(p string) error {
				_, err := params.Int(p)
				return err
			},
			name: "convert to int should fail",
		},
		{
			assertFunc: func(p string) error {
				_, err := params.Uint(p)
				return err
			},
			name: "convert to uint should fail",
		},
		{
			assertFunc: func(p string) error {
				_, err := params.Bool(p)
				return err
			},
			name: "convert to bool should fail",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.EqualValues(t, httprouter.NewErrorf(500, "uri param is not found: %s", "non_existing_param"), tt.assertFunc("non_existing_param"))
		})
	}
}

func TestParamsConversionFunctionsInvalidValueType(t *testing.T) {
	params := httprouter.URIParams{
		"noint":  "wish I'd be a number",
		"nobool": "fuzzy",
		"nouint": "spread love anywhere you go",
	}

	type testParams struct {
		assertFunc func(string) error
		name       string
		errMsgFmt  string
	}

	tests := map[string]testParams{
		"noint": {
			assertFunc: func(p string) error {
				_, err := params.Int(p)
				return err
			},
			name:      "invalid int type conversion",
			errMsgFmt: "uri param %s is not an int value: %s",
		},
		"nouint": {
			assertFunc: func(p string) error {
				_, err := params.Uint(p)
				return err
			},
			name:      "invalid uint type conversion",
			errMsgFmt: "uri param %s is not an uint value: %s",
		},
		"nobool": {
			assertFunc: func(p string) error {
				_, err := params.Bool(p)
				return err
			},
			name:      "invalid bool type conversion",
			errMsgFmt: "uri param %s is not an bool value: %s",
		},
	}

	for k, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, httprouter.NewErrorf(400, tt.errMsgFmt, k, params[k]), tt.assertFunc(k))
		})
	}
}
