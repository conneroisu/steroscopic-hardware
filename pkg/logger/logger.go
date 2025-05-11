package logger

import (
	"bytes"
	"io"
	"log/slog"
	"os"
	"strings"
	"time"
)

// Logger is a slog.Logger that sends logs to a channel and also to the console.
type Logger struct {
	*slog.Logger
	buffer *bytes.Buffer
}

// Bytes returns the buffered log.
func (l Logger) Bytes() []byte {
	return l.buffer.Bytes()
}

// NewLogger creates a new Logger.
func NewLogger() Logger {
	var buffer bytes.Buffer
	multiHandler := NewMultiHandler(
		NewLogWriter(&buffer),
		NewLogWriter(os.Stdout))
	logger := slog.New(multiHandler)
	slog.SetDefault(logger)
	return Logger{
		buffer: &buffer,
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

// NewLogWriter returns a slog.Handler that writes to a buffer.
func NewLogWriter(w io.Writer) slog.Handler {
	// consoleHandler is a default logger.
	return slog.NewTextHandler(w, &slog.HandlerOptions{
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
				a.Key = "src"
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
}
