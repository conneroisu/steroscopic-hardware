package logger

import (
	"log/slog"
	"os"
	"strings"
	"time"
)

// Logger is a slog.Logger that sends logs to a channel and also to the console.
type Logger struct {
	ch chan LogEntry
	*slog.Logger
}

// Channel returns the channel to which logs are sent to the browser.
func (l *Logger) Channel() chan LogEntry {
	return l.ch
}

// SetLogger sets the default logger to a slog.Logger.
func SetLogger() chan LogEntry {
	ch := make(chan LogEntry, 100)
	channelHandler := NewChannelHandler(ch, slog.LevelInfo)

	// Combine handlers
	multiHandler := NewMultiHandler(channelHandler, consoleHandler)
	logger := slog.New(multiHandler)
	slog.SetDefault(logger)
	return ch
}

// NewLogger creates a new Logger.
func NewLogger() Logger {
	logger := slog.New(consoleHandler)
	return Logger{
		ch:     make(chan LogEntry, 100),
		Logger: logger,
	}
}

// LogEntry represents a structured log entry
type LogEntry struct {
	Level   slog.Level
	Time    time.Time
	Message string
	Attrs   []slog.Attr
}

// consoleHandler is a default logger.
var consoleHandler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
	AddSource: true,
	Level:     slog.LevelDebug,
	ReplaceAttr: func(_ []string, a slog.Attr) slog.Attr {
		if a.Key == "time" {
			return slog.Attr{}
		}
		if a.Key == "level" {
			return slog.Attr{}
		}
		if a.Key == slog.SourceKey {
			str := a.Value.String()
			split := strings.Split(str, "/")
			if len(split) > 2 {
				a.Value = slog.StringValue(
					strings.Join(split[len(split)-2:], "/"),
				)
				a.Value = slog.StringValue(
					strings.ReplaceAll(a.Value.String(), "}", ""),
				)
			}
		}
		return a
	}})
