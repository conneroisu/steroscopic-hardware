package camera

import (
	"context"
	"fmt"
	"sync"
)

// Manager handles all cameras and their associated channels.
type Manager interface {
	// GetCamera retrieves a camera by type
	GetCamera(typ Type) Camera
	// SetCamera configures and starts a new camera of the specified type
	SetCamera(ctx context.Context, typ Type, cam Camera) error
	// GetChannel returns the input channel for the specified camera type
	GetChannel(typ Type) ImageChannel
	// GetOutputChannel returns the output channel for the specified camera type
	GetOutputChannel(typ Type) ImageChannel
	// CloseAll closes all cameras and releases their resources
	CloseAll() error
	// DrainAll empties all camera channels
	DrainAll()
}

// channelBufferSize defines how many images each channel can buffer.
const channelBufferSize = 5

// manager implements the Manager interface.
type manager struct {
	cameras     map[Type]Camera
	channels    map[Type]ImageChannel
	outChannels map[Type]ImageChannel
	mu          sync.RWMutex
}

// NewManager creates a new camera manager instance.
func NewManager() Manager {
	m := &manager{
		cameras:     make(map[Type]Camera),
		channels:    make(map[Type]ImageChannel),
		outChannels: make(map[Type]ImageChannel),
	}

	// Initialize channels
	m.channels[LeftCameraType] = make(ImageChannel, channelBufferSize)
	m.channels[RightCameraType] = make(ImageChannel, channelBufferSize)
	m.channels[OutputCameraType] = make(ImageChannel, channelBufferSize)

	// Initialize output channels (not needed for output camera)
	m.outChannels[LeftCameraType] = make(ImageChannel, channelBufferSize)
	m.outChannels[RightCameraType] = make(ImageChannel, channelBufferSize)

	return m
}

// GetCamera retrieves a camera by type.
func (m *manager) GetCamera(typ Type) Camera {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.cameras[typ]
}

// SetCamera configures and starts a new camera of the specified type.
func (m *manager) SetCamera(ctx context.Context, typ Type, cam Camera) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Pause other cameras to prevent race conditions
	for t, c := range m.cameras {
		if t != typ {
			c.Pause()
		}
	}

	// Close existing camera if any
	if oldCam, exists := m.cameras[typ]; exists {
		if err := oldCam.Close(); err != nil {
			// Resume other cameras before returning
			for t, c := range m.cameras {
				if t != typ {
					c.Resume()
				}
			}
			return fmt.Errorf("failed to close existing %s camera: %w", typ, err)
		}
	}

	// Drain channels
	m.drainChannels()

	// Store and start new camera
	m.cameras[typ] = cam
	go cam.Stream(ctx, m.channels[typ])

	// Resume other cameras
	for t, c := range m.cameras {
		if t != typ {
			c.Resume()
		}
	}

	return nil
}

// GetChannel returns the input channel for the specified camera type.
func (m *manager) GetChannel(typ Type) ImageChannel {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.channels[typ]
}

// GetOutputChannel returns the output channel for the specified camera type.
func (m *manager) GetOutputChannel(typ Type) ImageChannel {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.outChannels[typ]
}

// CloseAll closes all cameras and releases their resources.
func (m *manager) CloseAll() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	for typ, cam := range m.cameras {
		if err := cam.Close(); err != nil {
			return fmt.Errorf("failed to close %s camera: %w", typ, err)
		}
		delete(m.cameras, typ)
	}

	return nil
}

// DrainAll empties all camera channels.
func (m *manager) DrainAll() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.drainChannels()
}

// drainChannels non-blocking drain of all channels.
func (m *manager) drainChannels() {
	for {
		allEmpty := true

		// Try to read from each channel
		select {
		case <-m.channels[LeftCameraType]:
			allEmpty = false
		case <-m.channels[RightCameraType]:
			allEmpty = false
		case <-m.channels[OutputCameraType]:
			allEmpty = false
		case <-m.outChannels[LeftCameraType]:
			allEmpty = false
		case <-m.outChannels[RightCameraType]:
			allEmpty = false
		default:
			// No more data in any channel
		}

		if allEmpty {
			break
		}
	}
}

// Global manager instance.
var defaultManager = NewManager()

// DefaultManager returns the default camera manager.
func DefaultManager() Manager {
	return defaultManager
}

// SetDefaultManager replaces the default camera manager.
func SetDefaultManager(m Manager) {
	defaultManager = m
}

// GetCamera gets a camera by type from the default manager.
func GetCamera(typ Type) Camera {
	return defaultManager.GetCamera(typ)
}

// SetCamera sets a camera by type in the default manager.
func SetCamera(ctx context.Context, typ Type, cam Camera) error {
	return defaultManager.SetCamera(ctx, typ, cam)
}

// GetChannel gets a channel by type from the default manager.
func GetChannel(typ Type) ImageChannel {
	return defaultManager.GetChannel(typ)
}

// GetOutputChannel gets an output channel by type from the default manager.
func GetOutputChannel(typ Type) ImageChannel {
	return defaultManager.GetOutputChannel(typ)
}

// CloseAll closes all cameras in the default manager.
func CloseAll() error {
	return defaultManager.CloseAll()
}

// DrainAll drains all channels in the default manager.
func DrainAll() {
	defaultManager.DrainAll()
}
