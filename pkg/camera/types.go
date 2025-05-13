// Package camera defines types and interfaces for camera devices and their configuration
// in a stereoscopic hardware system.
package camera

import (
	"context"
	"image"
)

// ImageChannel is a typed channel for transmitting image.Gray frames between cameras and processing routines.
type ImageChannel chan *image.Gray

// Type identifies the type of camera (left, right, output).
type Type string

const (
	// LeftCameraType is the type for the left camera.
	LeftCameraType Type = "left"
	// RightCameraType is the type for the right camera.
	RightCameraType Type = "right"
	// OutputCameraType is the type for the output camera (e.g., depth map output).
	OutputCameraType Type = "output"
)

// Config represents all configurable camera parameters, such as serial port, baud rate, and compression.
type Config struct {
	Port        string // Serial port name or identifier
	BaudRate    int    // Baud rate for serial communication
	Compression int    // Compression level or mode
}

// Camera defines the interface that all camera types must implement. It abstracts streaming,
// configuration, and lifecycle management for different camera implementations.
type Camera interface {
	// Stream reads images and sends them to the provided channel.
	Stream(ctx context.Context, outCh ImageChannel)
	// Close releases all resources and stops any ongoing streaming.
	Close() error
	// Config returns the current configuration of the camera.
	Config() *Config
	// Pause temporarily stops streaming.
	Pause()
	// Resume restarts streaming after being paused.
	Resume()
	// Type returns the camera type (left, right, output).
	Type() Type
}
