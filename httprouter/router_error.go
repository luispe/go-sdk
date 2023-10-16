package httprouter

import (
	"errors"
)

// HandlerError represents an error returned by an ErrorHandler.
type HandlerError struct {
	StatusCode int
	Error      any
	Notify     bool
}

// ErrorHandlerFunc is used to define centralized error handler for your application.
type ErrorHandlerFunc func(err error, defaultHandlerError func(error) HandlerError) HandlerError

// DefaultHandlerError is a default implementation of an error handler.
// It converts the given error into a HandlerError object.
func DefaultHandlerError(err error) HandlerError {
	var webErr *Error
	if !errors.As(err, &webErr) {
		webErr = NewErrorf(500, err.Error()).(*Error)
	}

	return HandlerError{
		StatusCode: webErr.StatusCode,
		Error:      webErr,
		Notify:     webErr.StatusCode >= 500 && webErr.StatusCode <= 599,
	}
}
