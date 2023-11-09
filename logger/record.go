package logger

import (
	"log/slog"
	"time"
)

// Level represents different logging levels.
type Level slog.Level

// A set of possible logging levels.
const (
	LevelDebug = Level(slog.LevelDebug)
	LevelInfo  = Level(slog.LevelInfo)
	LevelWarn  = Level(slog.LevelWarn)
	LevelError = Level(slog.LevelError)
)

// Record represents the data that is being logged.
type Record struct {
	Time       time.Time
	Message    string
	Level      Level
	Attributes map[string]any
}

func toRecord(r slog.Record) Record {
	attrs := make(map[string]any, r.NumAttrs())

	f := func(attr slog.Attr) bool {
		attrs[attr.Key] = attr.Value.Any()
		return true
	}
	r.Attrs(f)

	return Record{
		Time:       r.Time,
		Message:    r.Message,
		Level:      Level(r.Level),
		Attributes: attrs,
	}
}

// LevelToString converts a given Level value to its corresponding string representation.
func (l Level) LevelToString() string {
	switch {
	case l == LevelDebug:
		return "DEBUG"
	case l == LevelInfo:
		return "INFO"
	case l == LevelWarn:
		return "WARN"
	case l == LevelError:
		return "ERROR"
	default:
		return ""
	}
}

// StringToLogLevel converts a given string log level value to its corresponding
// Level representation. By default, return LevelInfo
func StringToLogLevel(stringLevel string) Level {
	switch stringLevel {
	case "DEBUG":
		return LevelDebug
	case "WARN":
		return LevelWarn
	case "ERROR":
		return LevelError
	default:
		return LevelInfo
	}
}
