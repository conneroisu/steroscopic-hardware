package handlers

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"io"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/conneroisu/steroscopic-hardware/pkg/camera"
)

// StreamHandlerFn returns a handler for streaming camera images to multiple clients
func StreamHandlerFn(manager *camera.StreamManager) APIFn {
	// Make sure manager is running
	manager.Start()
	var jpegPool = sync.Pool{
		New: func() any {
			return &jpeg.Options{Quality: 75} // Lower quality for faster encoding
		},
	}

	return func(w http.ResponseWriter, r *http.Request) error {
		// Set headers for MJPEG stream
		w.Header().Set("Content-Type", "multipart/x-mixed-replace; boundary=frame")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "close")
		w.Header().Set("Pragma", "no-cache")

		clientChan := make(chan *image.Gray, 10) // Buffer a few frames

		// Register client and defer unregistering
		manager.Register <- clientChan
		defer func() {
			manager.Unregister <- clientChan
			// Drain the channel to prevent goroutine leaks
			for range clientChan {
				// Just drain
				print("")
			}
		}()

		// Set a reasonable timeout for the connection
		timeout := time.After(30 * time.Minute)

		// Create a buffer to avoid reallocating for each frame
		buffer := new(bytes.Buffer)
		buffer.Grow(1024 * 1024) // Pre-allocate 1MB
		// Control frame rate - don't send more than X frames per second
		ticker := time.NewTicker(time.Second / 10) // 10 FPS max
		defer ticker.Stop()

		for {
			select {
			case <-timeout:
				return nil
			case <-r.Context().Done():
				return nil
			case <-ticker.C:
				// Only process on tick to control frame rate
				select {
				case img, ok := <-clientChan:
					if !ok {
						return nil // Channel closed
					}
					err := processImg(
						img,
						buffer,
						&jpegPool,
						w,
					)
					if err != nil {
						return err
					}
				default:
					// No new frame available, continue
				}
			}
		}
	}
}

func processImg(
	img *image.Gray,
	buffer *bytes.Buffer,
	jpegPool *sync.Pool,
	w io.Writer,
) error {
	// Clear buffer and reuse
	buffer.Reset()

	// Get encoder options from pool
	opts := jpegPool.Get().(*jpeg.Options)

	// Write frame boundary to buffer
	_, err := fmt.Fprintf(buffer, "--frame\r\nContent-Type: image/jpeg\r\n\r\n")
	if err != nil {
		return err
	}

	// Encode image to buffer - consider using a worker pool for this
	err = jpeg.Encode(buffer, img, opts)
	if err != nil {
		jpegPool.Put(opts)
		log.Printf("Error encoding JPEG: %v", err)
		// Skip this frame instead of failing
		return nil
	}

	// Return options to pool
	jpegPool.Put(opts)

	// Write the entire frame at once instead of in chunks
	_, err = w.Write(buffer.Bytes())
	if err != nil {
		log.Printf("Error writing to client: %v", err)
		return fmt.Errorf("error writing to client: %v", err)
	}

	// Flush after writing complete frame
	f, ok := w.(http.Flusher)
	if ok {
		f.Flush()
	}

	return nil
}
