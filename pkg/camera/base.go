// Package camera provides interfaces and implementations for different camera types
// (e.g., static, serial, output) and their management in a stereoscopic hardware system.
package camera

import (
	"context"
	"sync"
)

// BaseCamera provides common functionality for camera implementations, including
// context management, pausing, and configuration storage. It is intended to be embedded
// in concrete camera types.
type BaseCamera struct {
	ctx    context.Context    // Context for cancellation and lifecycle management
	cancel context.CancelFunc // Function to cancel the context
	paused bool               // Indicates if the camera is paused
	cType  Type               // The camera type (left, right, output)
	mu     sync.Mutex         // Mutex for synchronizing access
	config Config             // Current camera configuration
}

// NewBaseCamera creates a new BaseCamera with the specified type and parent context.
// The returned BaseCamera has its own cancellable context.
func NewBaseCamera(ctx context.Context, cType Type) BaseCamera {
	childCtx, cancel := context.WithCancel(ctx)

	return BaseCamera{
		ctx:    childCtx,
		cancel: cancel,
		cType:  cType,
		config: Config{},
	}
}

// Type returns the camera type (left, right, output).
func (b *BaseCamera) Type() Type {
	return b.cType
}

// Config returns a pointer to the current configuration of the camera.
func (b *BaseCamera) Config() *Config {
	return &b.config
}

// Pause sets the camera's paused state to true, temporarily stopping streaming.
func (b *BaseCamera) Pause() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.paused = true
}

// Resume sets the camera's paused state to false, resuming streaming if paused.
func (b *BaseCamera) Resume() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.paused = false
}

// IsPaused returns whether the camera is currently paused.
func (b *BaseCamera) IsPaused() bool {
	b.mu.Lock()
	defer b.mu.Unlock()

	return b.paused
}

// Context returns the camera's context for cancellation and lifecycle control.
func (b *BaseCamera) Context() context.Context {
	return b.ctx
}

// Cancel cancels the camera's context, stopping all operations.
func (b *BaseCamera) Cancel() {
	b.cancel()
}

// SetConfig sets the camera configuration to the provided value.
func (b *BaseCamera) SetConfig(cfg Config) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.config = cfg
}
