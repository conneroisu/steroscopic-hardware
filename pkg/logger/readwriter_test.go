package logger

import (
	"io"
	"log/slog"
	"os"
	"testing"
)

// Example of a simple in-memory ReadWriter for demonstration
type MemoryReadWriter struct {
	buffer []byte
	pos    int
}

func NewMemoryReadWriter() *MemoryReadWriter {
	return &MemoryReadWriter{
		buffer: make([]byte, 0, 1024),
		pos:    0,
	}
}

func (m *MemoryReadWriter) Read(p []byte) (n int, err error) {
	if m.pos >= len(m.buffer) {
		return 0, io.EOF
	}

	n = copy(p, m.buffer[m.pos:])
	m.pos += n
	return n, nil
}

func (m *MemoryReadWriter) Write(p []byte) (n int, err error) {
	m.buffer = append(m.buffer, p...)
	return len(p), nil
}
func TestReadWriter(t *testing.T) {
	t.Log("Testing ReadWriter")
	// Create a logger
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	// Create the underlying ReadWriter
	rw := NewMemoryReadWriter()

	// Wrap it with logging
	loggingRW := NewLoggingReadWriter(rw, logger, "demo-rw")

	// Example usage
	data := []byte("Hello, world!")

	// Write data (will be logged)
	n, err := loggingRW.Write(data)
	if err != nil {
		logger.Error("write failed", "error", err)
		return
	}
	logger.Info("write completed", "bytes", n)

	// Read data (will be logged)
	readBuf := make([]byte, 20)
	n, err = loggingRW.Read(readBuf)
	if err != nil && err != io.EOF {
		logger.Error("read failed", "error", err)
		return
	}
	logger.Info("read completed", "bytes", n, "data", string(readBuf[:n]))

	// Example with a file
	file, err := os.Create("example.txt")
	if err != nil {
		logger.Error("failed to create file", "error", err)
		return
	}
	defer file.Close()

	// Create a ReadWriter that logs and writes to both file and original destination
	// This demonstrates a "tee" pattern by writing to multiple destinations
	teeWriter := io.MultiWriter(file, rw)
	teeLoggingWriter := NewLoggingReadWriter(
		struct {
			io.Reader
			io.Writer
		}{
			Reader: rw,
			Writer: teeWriter,
		},
		logger,
		"tee-example",
	)

	// Write using the tee setup
	_, err = teeLoggingWriter.Write([]byte("This goes to file and memory, with logging!"))
	if err != nil {
		t.Error("failed to write to tee")
		return
	}
}
