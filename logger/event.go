package logger

import (
	"context"
)

// EventFunc is a function to be executed when configured against a logger level.
type EventFunc func(ctx context.Context, r Record)

// Events contains an assignment of an event function to a logger level.
type Events struct {
	Debug EventFunc
	Info  EventFunc
	Warn  EventFunc
	Error EventFunc
}
