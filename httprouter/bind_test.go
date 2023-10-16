package httprouter_test

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/pomelo-la/go-toolkit/httprouter"
)

func TestBind_JSON(t *testing.T) {
	type s struct {
		Field1 string `json:"field1" validate:"required"`
	}

	type s2 struct {
		s
		Field2 []string `json:"field2" validate:"required"`
	}

	type s3 []json.RawMessage

	tt := []struct {
		name               string
		input              string
		expectedErr        string
		expectedStatusCode int
		destination        interface{}
		assertFunc         func(d interface{})
	}{
		{
			name:               "should return bad request when body is nil",
			expectedErr:        "400 bad_request: Request body is empty",
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:               "should return validation error when missing field value",
			input:              `{"field1":""}`,
			expectedErr:        "422 unprocessable_entity: validation_error: invalid fields: Field1",
			destination:        &s{},
			expectedStatusCode: http.StatusUnprocessableEntity,
		},
		{
			name:               "should return validation error when missing multiple field value",
			input:              `{"field1":""}`,
			expectedErr:        "422 unprocessable_entity: validation_error: invalid fields: Field1,Field2",
			destination:        &s2{},
			expectedStatusCode: http.StatusUnprocessableEntity,
		},
		{
			name:               "should return unmarshal type error when json is valid but dest is not",
			input:              `{"field1":"1", "field2":"2"}`,
			expectedErr:        "400 bad_request: Unmarshal type error: expected=[]string, got=string, field=field2, offset=27",
			destination:        &s2{},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:               "should return syntax error when json is not valid",
			input:              "invalid content",
			expectedErr:        "400 bad_request: Syntax error: offset=1, error=invalid character 'i' looking for beginning of value",
			destination:        &s{},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:               "should return json unmarshal error",
			input:              `{}`,
			expectedErr:        "400 bad_request: json: Unmarshal(non-pointer chan string)",
			destination:        make(chan string),
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:        "should bind s",
			input:       `{"field1":"1"}`,
			destination: &s{},
			assertFunc: func(d interface{}) {
				s, ok := d.(*s)
				require.True(t, ok)
				require.Equal(t, s.Field1, "1")
			},
		},
		{
			name:        "should bind s2",
			input:       `{"field1":"1", "field2": ["1","2"]}`,
			destination: &s2{},
			assertFunc: func(d interface{}) {
				s, ok := d.(*s2)
				require.True(t, ok)
				require.Equal(t, s.Field1, "1")
				require.Len(t, s.Field2, 2)
			},
		},
		{
			name:        "unsupported type",
			input:       `[{"field1":"1"},{"field1":"2"}]`,
			destination: &s3{},
			assertFunc: func(d interface{}) {
				s, ok := d.(*s3)
				require.True(t, ok)
				require.Len(t, *s, 2)
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			req := &http.Request{
				Body: io.NopCloser(strings.NewReader(tc.input)),
			}

			err := httprouter.Bind(req, tc.destination)
			if tc.expectedErr != "" {
				require.EqualError(t, err, tc.expectedErr)
				webErr := err.(*httprouter.Error)
				require.Equal(t, tc.expectedStatusCode, webErr.StatusCode)
			} else {
				require.NoError(t, err)
				tc.assertFunc(tc.destination)
			}
		})
	}
}

func TestBind_UnsupportedMediaType(t *testing.T) {
	h := http.Header{}
	h.Add("Content-Type", "application/xml")
	r := http.Request{
		Header: h,
	}

	err := httprouter.Bind(&r, nil)
	webErr, ok := err.(*httprouter.Error)
	require.True(t, ok)
	require.Equal(t, http.StatusUnsupportedMediaType, webErr.StatusCode)
}
