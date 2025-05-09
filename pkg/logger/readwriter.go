package logger

import (
	"io"
	"log/slog"
	"time"
)

// LoggingReadWriter wraps an io.ReadWriter and logs all Read/Write operations
type LoggingReadWriter struct {
	wrapped io.ReadWriter
	logger  *slog.Logger
	prefix  string
}

// NewLoggingReadWriter creates a new LoggingReadWriter
func NewLoggingReadWriter(
	wrapped io.ReadWriter,
	logger *slog.Logger,
	prefix string,
) *LoggingReadWriter {
	return &LoggingReadWriter{
		wrapped: wrapped,
		logger:  logger,
		prefix:  prefix,
	}
}

// Read implements io.Reader
func (l *LoggingReadWriter) Read(p []byte) (n int, err error) {
	n, err = l.wrapped.Read(p)

	// Log the read operation
	l.logger.Info("data read",
		"prefix", l.prefix,
		"bytes", n,
		"error", err,
		"data", string(p[:n]),
		"timestamp", time.Now(),
	)

	return n, err
}

// Write implements io.Writer
func (l *LoggingReadWriter) Write(p []byte) (n int, err error) {
	// Log the write operation before writing
	l.logger.Info("data write",
		"prefix", l.prefix,
		"bytes", len(p),
		"data", string(p),
		"timestamp", time.Now(),
	)

	n, err = l.wrapped.Write(p)

	// Log any errors
	if err != nil {
		l.logger.Error("write error",
			"prefix", l.prefix,
			"error", err,
			"bytes_attempted", len(p),
			"bytes_written", n,
		)
	}

	return n, err
}
