package camera

import (
	"context"
	"image"
)

// ImageChannel is a typed channel for image.Gray frames.
type ImageChannel chan *image.Gray

// Type identifies the type of camera.
type Type string

const (
	// LeftCameraType is the type for the left camera.
	LeftCameraType Type = "left"
	// RightCameraType is the type for the right camera.
	RightCameraType Type = "right"
	// OutputCameraType is the type for the output camera.
	OutputCameraType Type = "output"
)

// Config represents all configurable camera parameters.
type Config struct {
	Port        string
	BaudRate    int
	Compression int
}

// Camera defines the interface that all camera types must implement.
type Camera interface {
	// Stream reads images and sends them to the provided channel
	Stream(ctx context.Context, outCh ImageChannel)
	// Close releases all resources and stops any ongoing streaming
	Close() error
	// Config returns the current configuration of the camera
	Config() *Config
	// Pause temporarily stops streaming
	Pause()
	// Resume restarts streaming after being paused
	Resume()
	// Type returns the camera type (left, right, output)
	Type() Type
}
