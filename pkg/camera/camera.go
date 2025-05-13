// Package camera provides interfaces and implementations for cameras.
package camera

import (
	"context"
	"fmt"
	"image"
	"sync/atomic"
)

type (
	ImageChannel chan image.Gray
	// Camera represents a camera.
	Camera struct {
		Camer
	}
	// Camer is the interface for a camera.
	Camer interface {
		// Stream reads images from the camera and
		// sends them to the channel.
		Stream(left *image.Gray, right *image.Gray) <-chan image.Gray
		// Close closes the camera.
		Close() error
		// Config returns the current configuration
		// of the camera.
		Config() *Config
		// Pause pauses the camera.
		Pause()
		// Resume resumes the camera.
		Resume()
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
	defaultRightCamera   = atomic.Pointer[Camera]{}
	defaultOutputCamera  = atomic.Pointer[Camera]{}
	defaultLeftCh        = atomic.Pointer[chan image.Gray]{}
	defaultLeftOutputCh  = atomic.Pointer[chan image.Gray]{}
	defaultRightCh       = atomic.Pointer[chan image.Gray]{}
	defaultRightOutputCh = atomic.Pointer[chan image.Gray]{}
	defaultOutputCh      = atomic.Pointer[chan image.Gray]{}
)

// Buffer size for camera channels. Adjust as needed.
const channelBufferSize = 5

func init() {
	leftCh := make(chan image.Gray, channelBufferSize)
	leftOutputCh := make(chan image.Gray, channelBufferSize) // Still potentially problematic logic, but usable
	rightCh := make(chan image.Gray, channelBufferSize)
	rightOutputCh := make(chan image.Gray, channelBufferSize) // Still potentially problematic logic, but usable
	outputCh := make(chan image.Gray, channelBufferSize)
	defaultLeftCh.Store(&leftCh)
	defaultLeftOutputCh.Store(&leftOutputCh)
	defaultRightCh.Store(&rightCh)
	defaultRightOutputCh.Store(&rightOutputCh)
	defaultOutputCh.Store(&outputCh)
}

// LeftCh returns the left camera channel.
func LeftCh() ImageChannel { return *defaultLeftCh.Load() }

// LeftOutputCh returns the left camera output channel.
//
// This is the channel that the output camera reads from.
func LeftOutputCh() ImageChannel { return *defaultLeftOutputCh.Load() }

// RightCh returns the right camera channel.
func RightCh() ImageChannel { return *defaultRightCh.Load() }

// RightOutputCh returns the right camera output channel.
//
// This is the channel that the output camera reads from.
func RightOutputCh() ImageChannel { return *defaultRightOutputCh.Load() }

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

// DrainAll drains all camera channels.
func DrainAll() {
L:
	for {
		select {
		case <-OutputCh():
		case <-LeftCh():
		case <-RightCh():
		case <-LeftOutputCh():
		case <-RightOutputCh():
		default:
			break L
		}
	}
}

// OutputCh returns the output channel.
func OutputCh() chan *image.Gray { return *defaultOutputCh.Load() }

// SetOutputCamera sets the output camera.
func SetOutputCamera(
	ctx context.Context,
	cam Camer,
) {
	left := defaultLeftCamera.Load()
	if left != nil {
		left.Pause()
		defer left.Resume()
	}
	right := defaultRightCamera.Load()
	if right != nil {
		right.Pause()
		defer right.Resume()
	}
	old := defaultOutputCamera.Load()
	if old != nil {
		err := old.Close()
		if err != nil {
			panic(fmt.Errorf("failed to close old output camera: %w", err))
		}
	}
	DrainAll()
	defaultOutputCamera.Store(&Camera{cam})
	go cam.Stream(ctx, OutputCh())
}

// SetLeftCamera sets the left camera.
func SetLeftCamera(
	ctx context.Context,
	cam Camer,
) {
	output := defaultOutputCamera.Load()
	if output != nil {
		output.Pause()
		defer output.Resume()
	}
	right := defaultRightCamera.Load()
	if right != nil {
		right.Pause()
		defer right.Resume()
	}
	old := defaultLeftCamera.Load()
	if old != nil {
		err := old.Close()
		if err != nil {
			panic(fmt.Errorf("failed to close old left camera: %w", err))
		}
	}
	DrainAll()
	defaultLeftCamera.Store(&Camera{cam})
	go cam.Stream(ctx, LeftCh())
}

// SetRightCamera sets the right camera.
func SetRightCamera(
	ctx context.Context,
	cam Camer,
) {
	left := defaultLeftCamera.Load()
	if left != nil {
		left.Pause()
		defer left.Resume()
	}
	output := defaultOutputCamera.Load()
	if output != nil {
		output.Pause()
		defer output.Resume()
	}
	old := defaultRightCamera.Load()
	if old != nil {
		err := old.Close()
		if err != nil {
			panic(fmt.Errorf("failed to close old right camera: %w", err))
		}
	}
	DrainAll()
	defaultRightCamera.Store(&Camera{cam})
	go cam.Stream(ctx, RightCh())

}
