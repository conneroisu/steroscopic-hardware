// Package camera provides interfaces and implementations for cameras.
package camera

import (
	"context"
	"image"
	"sync/atomic"
)

type (
	// Camera represents a camera.
	Camera struct {
		Camer
	}
	// Camer is the interface for a camera.
	Camer interface {
		// Stream reads images from the camera and
		// sends them to the channel.
		Stream(context.Context, chan *image.Gray)
		// Close closes the camera.
		Close() error
		// Config returns the current configuration
		// of the camera.
		Config() *Config
	}
	// Config represents all configurable camera parameters.
	Config struct {
		Port        string
		BaudRate    int
		Compression int
	}
)

var (
	defaultLeftCamera    = atomic.Pointer[Camera]{}
	defaultLeftCh        = atomic.Pointer[chan *image.Gray]{}
	defaultLeftOutputCh  = atomic.Pointer[chan *image.Gray]{}
	defaultRightCamera   = atomic.Pointer[Camera]{}
	defaultRightCh       = atomic.Pointer[chan *image.Gray]{}
	defaultRightOutputCh = atomic.Pointer[chan *image.Gray]{}
	defaultOutputCamera  = atomic.Pointer[Camera]{}
	defaultOutputCh      = atomic.Pointer[chan *image.Gray]{}
)

// Buffer size for camera channels. Adjust as needed.
const channelBufferSize = 5

func init() {
	leftCh := make(chan *image.Gray, channelBufferSize)
	leftOutputCh := make(chan *image.Gray, channelBufferSize) // Still potentially problematic logic, but usable
	rightCh := make(chan *image.Gray, channelBufferSize)
	rightOutputCh := make(chan *image.Gray, channelBufferSize) // Still potentially problematic logic, but usable
	outputCh := make(chan *image.Gray, channelBufferSize)
	defaultLeftCh.Store(&leftCh)
	defaultLeftOutputCh.Store(&leftOutputCh)
	defaultRightCh.Store(&rightCh)
	defaultRightOutputCh.Store(&rightOutputCh)
	defaultOutputCh.Store(&outputCh)
}

// LeftCh returns the left camera channel.
func LeftCh() chan *image.Gray { return *defaultLeftCh.Load() }

// LeftOutputCh returns the left camera output channel.
//
// This is the channel that the output camera reads from.
func LeftOutputCh() chan *image.Gray { return *defaultLeftOutputCh.Load() }

// RightCh returns the right camera channel.
func RightCh() chan *image.Gray { return *defaultRightCh.Load() }

// RightOutputCh returns the right camera output channel.
//
// This is the channel that the output camera reads from.
func RightOutputCh() chan *image.Gray { return *defaultRightOutputCh.Load() }

// CloseAll closes all cameras.
func CloseAll() error {
	err := defaultLeftCamera.Load().Close()
	if err != nil {
		return err
	}
	err = defaultRightCamera.Load().Close()
	if err != nil {
		return err
	}
	err = defaultOutputCamera.Load().Close()
	if err != nil {
		return err
	}
	return nil
}

// OutputCh returns the output channel.
func OutputCh() chan *image.Gray { return *defaultOutputCh.Load() }

// SetOutputCamera sets the output camera.
func SetOutputCamera(
	ctx context.Context,
	cam Camer,
) {
	defaultOutputCamera.Store(&Camera{cam})
	go cam.Stream(ctx, OutputCh())
}

// SetLeftCamera sets the left camera.
func SetLeftCamera(
	ctx context.Context,
	cam Camer,
) {
	defaultLeftCamera.Store(&Camera{cam})
	go cam.Stream(ctx, LeftCh())
}

// SetRightCamera sets the right camera.
func SetRightCamera(
	ctx context.Context,
	cam Camer,
) {
	defaultRightCamera.Store(&Camera{cam})
	go cam.Stream(ctx, RightCh())
}
