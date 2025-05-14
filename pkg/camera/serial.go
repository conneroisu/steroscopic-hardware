package camera

import (
	"context"
	"errors"
	"fmt"
	"image"
	"image/color"
	"log/slog"
	"sync"
	"time"

	"github.com/conneroisu/steroscopic-hardware/pkg/homedir"
	"go.bug.st/serial"
)

var (
	// DefaultStartSeq is the default start marker for image data.
	DefaultStartSeq = []byte{0xff, 0xd8}

	// DefaultEndSeq is the default end marker for image data.
	DefaultEndSeq = []byte{0xff, 0xd9}

	// DefaultImageWidth is the default expected image width in pixels.
	DefaultImageWidth = 1920

	// DefaultImageHeight is the default expected image height in pixels.
	DefaultImageHeight = 1080
)

// SerialCamera represents a camera connected via serial port. It handles communication
// with hardware cameras over a serial interface, including image acquisition and streaming.
type SerialCamera struct {
	*BaseCamera
	port        serial.Port  // Serial port interface
	startSeq    []byte       // Start sequence for image data
	endSeq      []byte       // End sequence for image data
	imageWidth  int          // Expected image width in pixels
	imageHeight int          // Expected image height in pixels
	logger      *slog.Logger // Logger for serial camera events
	onClose     func()       // Cleanup function for closing the camera
	streamMu    sync.Mutex   // Mutex for synchronizing streaming
}

// NewSerialCamera creates a new SerialCamera instance for the given type, port, baud rate,
// and compression setting. It opens the serial port and prepares the camera for streaming.
func NewSerialCamera(ctx context.Context, typ Type, portName string, baudRate int, compression int) (*SerialCamera, error) {
	base := NewBaseCamera(ctx, typ)

	// Configure the camera
	base.SetConfig(Config{
		Port:        portName,
		BaudRate:    baudRate,
		Compression: compression,
	})

	sc := &SerialCamera{
		BaseCamera:  &base,
		startSeq:    DefaultStartSeq,
		endSeq:      DefaultEndSeq,
		imageWidth:  DefaultImageWidth,
		imageHeight: DefaultImageHeight,
		logger:      slog.Default().WithGroup(fmt.Sprintf("serial-camera-%s", typ)),
	}

	// Configure serial port
	mode := &serial.Mode{
		BaudRate: baudRate,
		DataBits: 8,
		Parity:   serial.NoParity,
		StopBits: serial.OneStopBit,
	}

	sc.logger.Info("opening serial port", "port", portName, "baudRate", baudRate)
	port, err := serial.Open(portName, mode)
	if err != nil {
		return nil, fmt.Errorf("failed to open serial port %s: %w", portName, err)
	}

	sc.port = port

	// Set read timeout
	err = sc.port.SetReadTimeout(serial.NoTimeout)
	if err != nil {
		sc.port.Close()
		return nil, fmt.Errorf("failed to set read timeout: %w", err)
	}

	return sc, nil
}

// Stream reads images from the camera and sends them to the provided channel. It manages
// the streaming lifecycle, error handling, and reconnection logic.
func (sc *SerialCamera) Stream(ctx context.Context, outCh ImageChannel) {
	sc.logger.Debug("SerialCamera.Stream started")
	defer sc.logger.Debug("SerialCamera.Stream completed")

	// Create error channel for internal communication
	errChan := make(chan error, 1)

	// Start the initial stream
	readFn, err := sc.initializeStream(ctx, errChan, outCh)
	if err != nil {
		sc.logger.Error("failed to initialize image stream", "err", err)
		return
	}

	// Launch the reading goroutine
	go readFn()

	// Monitor for errors and cancellation
	for {
		select {
		case <-ctx.Done():
			sc.logger.Debug("context canceled, stopping stream")
			return
		case <-sc.Context().Done():
			sc.logger.Debug("camera context canceled, stopping stream")
			return
		case err := <-errChan:
			sc.logger.Error("error in image stream", "err", err)
			// Could implement reconnection logic here if needed
		}
	}
}

// initializeStream sets up the serial connection and prepares for streaming. It returns a
// function that continuously reads frames and sends them to the image channel.
func (sc *SerialCamera) initializeStream(ctx context.Context, errChan chan error, imgCh ImageChannel) (func(), error) {
	var tries = 0

	for {
		sc.logger.Debug("initializing camera stream")

		// Send the start sequence to request an image
		_, err := sc.port.Write(sc.startSeq)
		if err != nil {
			return nil, fmt.Errorf("failed to send start sequence: %w", err)
		}

		// Read acknowledgement byte
		ackBuffer := make([]byte, 1)
		length, err := sc.port.Read(ackBuffer)
		if err != nil {
			return nil, fmt.Errorf("failed to read acknowledgement: %w", err)
		}

		if length != 1 {
			// If we didn't receive exactly one byte, try again
			tries++
			if tries > 4 {
				return nil, fmt.Errorf("camera not responding properly after %d attempts", tries)
			}

			// Send end sequence to reset the camera
			_, err := sc.port.Write(sc.endSeq)
			if err != nil {
				return nil, fmt.Errorf("failed to send end sequence: %w", err)
			}

			time.Sleep(100 * time.Millisecond)
			continue
		}

		// Set up the cleanup function
		sc.onClose = func() {
			sc.logger.Debug("sending end sequence to camera")
			_, err := sc.port.Write(sc.endSeq)
			if err != nil {
				sc.logger.Error("failed to send end sequence", "err", err)
			}
		}

		// Return the reading function
		return func() {
			// Backoff parameters
			initialBackoff := 10 * time.Millisecond
			maxBackoff := 1 * time.Second
			backoff := initialBackoff
			consecutiveFailures := 0

			for {
				select {
				case <-ctx.Done():
					return
				case <-sc.Context().Done():
					return
				default:
					if sc.IsPaused() {
						time.Sleep(100 * time.Millisecond)
						continue
					}

					// Lock to ensure we don't have multiple reads happening simultaneously
					sc.streamMu.Lock()
					img, err := sc.readFrame()
					sc.streamMu.Unlock()

					if err != nil {
						errChan <- err
						time.Sleep(500 * time.Millisecond) // Delay before retry
						continue
					}

					select {
					case imgCh <- img:
						sc.logger.Debug("image sent to channel")
						// Reset backoff on success
						backoff = initialBackoff
						consecutiveFailures = 0
					case <-ctx.Done():
						return
					case <-sc.Context().Done():
						return
					default:
						// Channel is full, apply backoff
						consecutiveFailures++

						// Only log at certain thresholds to prevent log spam
						if consecutiveFailures == 1 || consecutiveFailures%10 == 0 {
							sc.logger.Debug("output channel full, applying backoff",
								"consecutiveFailures", consecutiveFailures,
								"currentBackoff", backoff)
						}

						// Apply backoff delay
						time.Sleep(backoff)

						// Exponential backoff with a maximum cap
						backoff = time.Duration(float64(backoff) * 1.5)
						if backoff > maxBackoff {
							backoff = maxBackoff
						}
					}
				}
			}
		}, nil
	}
}

// readFrame reads a single image frame from the serial port, converts it to grayscale,
// and returns it as an image.Gray. It handles timeouts and progress reporting.
func (sc *SerialCamera) readFrame() (*image.Gray, error) {
	sc.logger.Debug("reading image frame")

	// Use a timeout for the read operation
	readCtx, cancel := context.WithTimeout(sc.Context(), 3*time.Minute)
	defer cancel()

	// Buffer to store image data
	var buffer []byte
	expectedLength := sc.imageWidth * sc.imageHeight

	// Monitor the read progress
	progressDone := make(chan struct{})
	go func() {
		start := time.Now()
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				sc.logger.Info("reading image data",
					"progress", fmt.Sprintf("%d/%d bytes (%.1f%%)",
						len(buffer), expectedLength,
						float64(len(buffer))/float64(expectedLength)*100),
					"elapsed", time.Since(start))
			case <-progressDone:
				return
			case <-readCtx.Done():
				return
			}
		}
	}()

	// Read the image data in chunks
	for len(buffer) < expectedLength {
		chunk := make([]byte, 1024)
		n, err := sc.port.Read(chunk)
		if err != nil {
			close(progressDone)
			return nil, fmt.Errorf("error reading from serial port: %w", err)
		}

		if n > 0 {
			buffer = append(buffer, chunk[:n]...)
		}

		// Check if the context has been canceled
		select {
		case <-readCtx.Done():
			close(progressDone)
			return nil, errors.New("read operation timed out or was canceled")
		default:
			// Continue reading
		}
	}

	close(progressDone)
	sc.logger.Info("image data read complete", "size", len(buffer))

	// Create grayscale image from the buffer
	img := image.NewGray(image.Rect(0, 0, sc.imageWidth, sc.imageHeight))

	// Copy buffer data to image
	for y := range sc.imageHeight {
		for x := range sc.imageWidth {
			i := y*sc.imageWidth + x
			if i < len(buffer) {
				img.SetGray(x, y, color.Gray{Y: buffer[i]})
			}
		}
	}

	// Save a copy of the image for debugging
	err := homedir.SaveImage(img)
	if err != nil {
		sc.logger.Error("failed to save debug image", "err", err)
	}

	return img, nil
}

// Close releases all resources used by the camera, including closing the serial port and
// canceling the context.
func (sc *SerialCamera) Close() error {
	sc.logger.Info("closing serial camera")

	// Cancel the context to stop all operations
	sc.Cancel()

	// Execute onClose handler if set
	if sc.onClose != nil {
		sc.onClose()
	}

	// Close the serial port
	if sc.port != nil {
		return sc.port.Close()
	}

	return nil
}
