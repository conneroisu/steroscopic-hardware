package camera

import (
	"context"
	"fmt"
	"sync"
)

// Manager handles all cameras and their associated channels. It provides methods to
// retrieve, configure, and control cameras and their image channels.
type Manager interface {
	// GetCamera retrieves a camera by type (left, right, output).
	GetCamera(typ Type) Camera
	// SetCamera configures and starts a new camera of the specified type.
	SetCamera(ctx context.Context, typ Type, cam Camera) error
	// CloseAll closes all cameras and releases their resources.
	CloseAll() error
}

// manager implements the Manager interface for camera management.
type manager struct {
	cameras map[Type]Camera // Map of camera type to camera instance
	mu      sync.RWMutex    // Mutex for concurrent access
}

// NewManager creates a new camera manager instance with initialized channels.
func NewManager() Manager {
	m := &manager{
		cameras: make(map[Type]Camera),
	}

	return m
}

// GetCamera retrieves a camera by type (left, right, output).
func (m *manager) GetCamera(typ Type) Camera {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return m.cameras[typ]
}

// SetCamera configures and starts a new camera of the specified type. It pauses other cameras,
// closes any existing camera of the same type, drains channels, and starts the new camera.
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
	oldCam, exists := m.cameras[typ]
	if exists {
		err := oldCam.Close()
		if err != nil {
			// Resume other cameras before returning
			for t, c := range m.cameras {
				if t != typ {
					c.Resume()
				}
			}

			return fmt.Errorf("failed to close existing %s camera: %w", typ, err)
		}
	}

	// Store and start new camera
	m.cameras[typ] = cam
	go cam.Stream(ctx, nil)

	// Resume other cameras
	for t, c := range m.cameras {
		if t != typ {
			c.Resume()
		}
	}

	return nil
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

// Global manager instance for default usage.
var defaultManager = NewManager()

// GetCamera gets a camera by type from the default manager.
func GetCamera(typ Type) Camera {
	return defaultManager.GetCamera(typ)
}

// SetCamera sets a camera by type in the default manager.
func SetCamera(ctx context.Context, typ Type, cam Camera) error {
	return defaultManager.SetCamera(ctx, typ, cam)
}

// CloseAll closes all cameras in the default manager.
func CloseAll() error {
	return defaultManager.CloseAll()
}
