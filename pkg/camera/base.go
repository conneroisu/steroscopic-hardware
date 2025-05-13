package camera

import (
	"context"
	"sync"
)

// BaseCamera provides common functionality for camera implementations.
type BaseCamera struct {
	ctx    context.Context
	cancel context.CancelFunc
	paused bool
	cType  Type
	mu     sync.Mutex
	config Config
}

// NewBaseCamera creates a new BaseCamera with the specified type.
func NewBaseCamera(ctx context.Context, cType Type) BaseCamera {
	childCtx, cancel := context.WithCancel(ctx)
	return BaseCamera{
		ctx:    childCtx,
		cancel: cancel,
		cType:  cType,
		config: Config{},
	}
}

// Type returns the camera type.
func (b *BaseCamera) Type() Type {
	return b.cType
}

// Config returns the current configuration of the camera.
func (b *BaseCamera) Config() *Config {
	return &b.config
}

// Pause pauses the camera.
func (b *BaseCamera) Pause() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.paused = true
}

// Resume resumes the camera.
func (b *BaseCamera) Resume() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.paused = false
}

// IsPaused returns whether the camera is paused.
func (b *BaseCamera) IsPaused() bool {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.paused
}

// Context returns the camera's context.
func (b *BaseCamera) Context() context.Context {
	return b.ctx
}

// Cancel cancels the camera's context.
func (b *BaseCamera) Cancel() {
	b.cancel()
}

// SetConfig sets the camera configuration.
func (b *BaseCamera) SetConfig(cfg Config) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.config = cfg
}
