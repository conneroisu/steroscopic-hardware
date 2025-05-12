// Package camera provides interfaces and implementations for cameras.
package camera

import (
	"context"
	"image"
)

// Camer is the interface for a camera.
type Camer interface {
	Stream(context.Context, chan *image.Gray)
	Close() error
	Port() string
}

// Config represents all configurable camera parameters.
type Config struct {
	Port        string
	BaudRate    int
	Compression int
}

var defaultConfig = Config{
	Port:        "/dev/ttyUSB0",
	BaudRate:    115200,
	Compression: 0,
}

// DefaultCameraConfig returns default camera configuration.
func DefaultCameraConfig() Config { return defaultConfig }
