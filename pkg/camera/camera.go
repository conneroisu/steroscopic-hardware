// Package camera provides interfaces and implementations for cameras.
package camera

import (
	"context"
	"image"
	"log/slog"
	"sync"
)

// Camer is the interface for a camera.
type Camer interface {
	Stream(context.Context, chan *image.Gray)
	Close() error
	ConfigurePort(int) error
}

// StreamManager manages multiple client connections to a single camera stream
type StreamManager struct {
	clients    map[chan *image.Gray]bool
	Register   chan chan *image.Gray
	Unregister chan chan *image.Gray
	camera     Camer
	frames     chan *image.Gray
	mu         sync.Mutex
	ctx        context.Context
	cancel     context.CancelFunc
	running    bool
}

// NewStreamManager creates a new broadcaster for the given camera
func NewStreamManager(camera Camer) *StreamManager {
	ctx, cancel := context.WithCancel(context.Background())
	return &StreamManager{
		clients:    make(map[chan *image.Gray]bool),
		Register:   make(chan chan *image.Gray),
		Unregister: make(chan chan *image.Gray),
		camera:     camera,
		frames:     make(chan *image.Gray),
		ctx:        ctx,
		cancel:     cancel,
		running:    false,
	}
}

// Lock locks the mutex
func (b *StreamManager) Lock() { b.mu.Lock() }

// Unlock unlocks the mutex
func (b *StreamManager) Unlock() { b.mu.Unlock() }

// Start begins streaming from the camera and broadcasting to clients
func (b *StreamManager) Start() {
	b.mu.Lock()
	if b.running {
		b.mu.Unlock()
		return
	}
	b.running = true
	b.mu.Unlock()

	// Start camera stream
	go b.camera.Stream(b.ctx, b.frames)

	// Main broadcasting loop
	go func() {
		for {
			select {
			case client := <-b.Register:
				b.mu.Lock()
				b.clients[client] = true
				slog.Debug("client registered", "total", len(b.clients))
				b.mu.Unlock()

			case client := <-b.Unregister:
				b.mu.Lock()
				if _, ok := b.clients[client]; ok {
					delete(b.clients, client)
					close(client)
					slog.Debug("client unregistered", "total", len(b.clients))
				}
				b.mu.Unlock()

			case frame := <-b.frames:
				// Broadcast frame to all clients
				b.mu.Lock()
				for client := range b.clients {
					// Non-blocking send - skip clients that are slow
					select {
					case client <- frame:
					default:
						// Client is too slow, drop frame for this client
					}
				}
				b.mu.Unlock()

			case <-b.ctx.Done():
				// Context canceled, clean up
				b.mu.Lock()
				for client := range b.clients {
					delete(b.clients, client)
					close(client)
				}
				b.running = false
				b.mu.Unlock()
				return
			}
		}
	}()
}

// Stop stops the broadcaster and disconnects all clients
func (b *StreamManager) Stop() {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.running {
		b.running = false
		b.cancel()
		// Create a new context for future clients
		b.ctx, b.cancel = context.WithCancel(context.Background())
	}
}
