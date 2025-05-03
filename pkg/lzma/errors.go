package lzma

import "fmt"

// HeaderError is returned when the header is corrupt.
type HeaderError struct {
	msg string
}

// Error returns the error message.
func (e *HeaderError) Error() string {
	return fmt.Sprintf("header error: %s", e.msg)
}

// StreamError is returned when the stream is corrupt.
type StreamError struct {
	msg string
}

// Error returns the error message and implements the error interface.
func (e *StreamError) Error() string {
	return fmt.Sprintf("stream error: %s", e.msg)
}

// NWriteError is returned when the number of bytes returned by Writer.Write() didn't meet expectances.
type NWriteError struct {
	msg string
}

// Error returns the error message and implements the error interface.
func (e *NWriteError) Error() string {
	return fmt.Sprintf("number of bytes returned by Writer.Write() didn't meet expectances: %s", e.msg)
}
