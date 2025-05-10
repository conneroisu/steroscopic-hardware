// Package camera provides interfaces and implementations for cameras.
package camera

import (
	"context"
	"image"
	"log/slog"
	"sync"

	"github.com/conneroisu/steroscopic-hardware/pkg/logger"
)

// Camer is the interface for a camera.
type Camer interface {
	Stream(context.Context, chan *image.Gray)
	Close() error
	Port() string
}

// Config represents all configurable camera parameters
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

// DefaultCameraConfig returns default camera configuration
func DefaultCameraConfig() Config {
	return defaultConfig
}

// StreamManager manages multiple client connections to a single camera stream
type StreamManager struct {
	clients    map[chan *image.Gray]bool
	Register   chan chan *image.Gray
	Unregister chan chan *image.Gray
	camera     Camer
	config     *Config
	frames     chan *image.Gray
	mu         sync.Mutex
	ctx        context.Context
	logger     *logger.Logger
	cancel     context.CancelFunc
	runCtx     context.Context
	runCancel  context.CancelFunc
	running    bool
}

// NewStreamManager creates a new broadcaster for the given camera
func NewStreamManager(camera Camer, logger *logger.Logger) *StreamManager {
	ctx, cancel := context.WithCancel(context.Background())
	return &StreamManager{
		clients:    make(map[chan *image.Gray]bool),
		Register:   make(chan chan *image.Gray),
		Unregister: make(chan chan *image.Gray),
		camera:     camera,
		frames:     make(chan *image.Gray),
		logger:     logger,
		ctx:        ctx,
		cancel:     cancel,
		config:     &defaultConfig,
		running:    false,
	}
}

// Lock locks the mutex
func (b *StreamManager) Lock() { b.mu.Lock() }

// Unlock unlocks the mutex
func (b *StreamManager) Unlock() { b.mu.Unlock() }

// Start begins streaming from the camera and broadcasting to clients
func (b *StreamManager) Start() {
	b.logger.Info("StreamManager.Start()")
	defer b.logger.Info("StreamManager.Start() done")
	b.mu.Lock()
	if b.running {
		b.mu.Unlock()
		return
	}
	b.running = true
	b.runCtx, b.runCancel = context.WithCancel(b.ctx)
	b.mu.Unlock()

	slog.Info("starting camera stream")
	// Start camera stream
	go b.camera.Stream(b.ctx, b.frames)

	slog.Info("starting broadcasting loop")
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
					case <-b.runCtx.Done():
						continue
					case <-b.ctx.Done():
						continue
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

			case <-b.runCtx.Done():
				// Run context canceled, clean up
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

// Configure configures the camera owned by this StreamManager.
func (b *StreamManager) Configure(config Config) error {
	var (
		err    error
		camera Camer
		opts   = []SerialCameraOption{}
	)
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.camera.Port() == config.Port {
		b.logger.Info("camera port already configured closing it to reconfigure", "port", config.Port)
		b.camera.Close()
	}
	opts = append(opts, WithLogger(b.logger))
	b.logger.Info(
		"opening new camera",
		"port",
		config.Port,
		"baud",
		config.BaudRate,
		"compression",
		config.Compression == 1,
	)
	camera, err = NewSerialCamera(
		config.Port,
		config.BaudRate,
		config.Compression == 1,
		opts...,
	)
	if err != nil {
		return err
	}

	b.runCancel()

	b.logger.Info(
		"closing old camera",
		"port",
		config.Port,
		"baud",
		config.BaudRate,
		"compression",
		config.Compression == 1,
	)
	err = b.camera.Close()
	if err != nil {
		return err
	}
	b.camera = camera
	b.logger.Info(
		"configured camera",
		"port",
		config.Port,
		"baud",
		config.BaudRate,
		"compression",
		config.Compression == 1,
	)
	go b.Start()
	b.config = &config
	return nil
}

// Config returns the current configuration of the camera owned by this StreamManager.
func (b *StreamManager) Config() *Config {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.config
}
