package httprouter

import (
	"context"
	"net/http"
	"strconv"
)

// URIParamsCtxKey is the context key used to store and retrieve URI parameters.
type URIParamsCtxKey struct{}

// URIParams contains the key-value combination of parameters from the URI.
type URIParams map[string]string

// Params returns a map from the given context containing every URI parameter defined in the route.
// The key represents the name of the route variable for the current request, if any.
func Params(r *http.Request) URIParams {
	if v, ok := r.Context().Value(URIParamsCtxKey{}).(URIParams); ok {
		return v
	}

	return nil
}

// WithParams returns a new Context that carries the provided params.
func WithParams(ctx context.Context, params URIParams) context.Context {
	return context.WithValue(ctx, URIParamsCtxKey{}, params)
}

// String gets parameter as a string.
// If the parameter is not found, it returns InternalServerError(500).
func (p URIParams) String(param string) (string, error) {
	v, ok := p[param]
	if !ok {
		return "", NewErrorf(http.StatusInternalServerError, "uri param is not found: %s", param)
	}

	return v, nil
}

// Int gets parameter as an int.
// If the parameter is not found, it returns InternalServerError(500).
// If the parameter type value is a not an int, it returns BadRequestError(400).
func (p URIParams) Int(param string) (int, error) {
	v, ok := p[param]
	if !ok {
		return 0, NewErrorf(http.StatusInternalServerError, "uri param is not found: %s", param)
	}

	paramParsed, err := strconv.Atoi(v)
	if err != nil {
		return 0, NewErrorf(http.StatusBadRequest, "uri param %s is not an int value: %s", param, v)
	}

	return paramParsed, nil
}

// Uint gets parameter as an uint.
// If the parameter is not found, it returns InternalServerError(500).
// If the parameter type value is a not an uint, it returns BadRequestError(400).
func (p URIParams) Uint(param string) (uint, error) {
	v, ok := p[param]
	if !ok {
		return 0, NewErrorf(http.StatusInternalServerError, "uri param is not found: %s", param)
	}

	paramParsed, err := strconv.ParseUint(v, 10, 0)
	if err != nil {
		return 0, NewErrorf(http.StatusBadRequest, "uri param %s is not an uint value: %s", param, v)
	}

	return uint(paramParsed), nil
}

// Bool gets parameter as a bool.
// If the parameter is not found, it returns InternalServerError(500).
// If the parameter type value is a not a bool, it returns BadRequestError(400).
func (p URIParams) Bool(param string) (bool, error) {
	v, ok := p[param]
	if !ok {
		return false, NewErrorf(http.StatusInternalServerError, "uri param is not found: %s", param)
	}

	parsedValue, err := strconv.ParseBool(v)
	if err != nil {
		return false, NewErrorf(http.StatusBadRequest, "uri param %s is not an bool value: %s", param, v)
	}

	return parsedValue, nil
}
