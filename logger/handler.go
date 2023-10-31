package logger

import (
	"context"
	"log/slog"
)

// logHandler provides a wrapper around the slog handler to capture which
// logger level is being logged for event handling.
type logHandler struct {
	handler slog.Handler
	events  Events
}

func newLogHandler(handler slog.Handler, events Events) *logHandler {
	return &logHandler{
		handler: handler,
		events:  events,
	}
}

// Enabled reports whether the handler handles records at the given level.
// The handler ignores records whose level is lower.
func (h *logHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.handler.Enabled(ctx, level)
}

// WithAttrs returns a newLogger JSONHandler whose attributes consists
// of h's attributes followed by attrs.
func (h *logHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &logHandler{handler: h.handler.WithAttrs(attrs), events: h.events}
}

// WithGroup returns a newLogger Handler with the given group appended to the receiver's
// existing groups. The keys of all subsequent attributes, whether added by With
// or in a Record, should be qualified by the sequence of group names.
func (h *logHandler) WithGroup(name string) slog.Handler {
	return &logHandler{handler: h.handler.WithGroup(name), events: h.events}
}

// Handle looks to see if an event function needs to be executed for a given
// logger level and then formats its argument Record.
func (h *logHandler) Handle(ctx context.Context, r slog.Record) error {
	logHandlers := map[slog.Level]EventFunc{
		slog.LevelDebug: h.events.Debug,
		slog.LevelError: h.events.Error,
		slog.LevelWarn:  h.events.Warn,
		slog.LevelInfo:  h.events.Info,
	}

	if handler, ok := logHandlers[r.Level]; ok {
		handler(ctx, toRecord(r))
	}

	return h.handler.Handle(ctx, r)
}
