package handlers

import (
	"context"
	"fmt"
	"image"
	"image/jpeg"
	"io"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"github.com/conneroisu/steroscopic-hardware/pkg/camera"
)

const (
	streamTimeout = 60 * time.Minute
	frameRate     = time.Second / 10 // 10 FPS max
)

var (
	// Pool of JPEG encoding options to avoid constant allocations.
	encodeOptsPool = sync.Pool{
		New: func() any {
			return &jpeg.Options{Quality: 75}
		},
	}
)

// processImg encodes and writes an image frame to the response writer.
func processImg(img *image.Gray, w io.Writer) error {
	// Write frame boundary
	if _, err := fmt.Fprintf(w, "--frame\r\nContent-Type: image/jpeg\r\n\r\n"); err != nil {
		return err
	}

	// Get encoding options from pool
	encodeOpts := encodeOptsPool.Get().(*jpeg.Options)
	defer encodeOptsPool.Put(encodeOpts)

	// Encode image
	if err := jpeg.Encode(w, img, encodeOpts); err != nil {
		slog.Error("Error encoding JPEG", "err", err)

		return nil // Skip this frame instead of failing the stream
	}

	// Flush after writing complete frame
	if f, ok := w.(http.Flusher); ok {
		f.Flush()
	}

	return nil
}

// HandleCameraStream is a generic handler for streaming camera images.
func HandleCameraStream(camType camera.Type, useOutputChannel bool) APIFn {
	logger := slog.Default().WithGroup(fmt.Sprintf("stream-%s", camType))

	return func(w http.ResponseWriter, r *http.Request) error {
		logger.Debug("stream requested")

		// Set MJPEG stream headers
		w.Header().Set("Content-Type", "multipart/x-mixed-replace; boundary=frame")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "close")
		w.Header().Set("Pragma", "no-cache")

		// Create a context with timeout
		ctx, cancel := context.WithTimeout(r.Context(), streamTimeout)
		defer cancel()

		// Set up ticker for frame rate control
		ticker := time.NewTicker(frameRate)
		defer ticker.Stop()

		// Get input channel to read from
		var inputCh camera.ImageChannel
		if useOutputChannel && (camType == camera.LeftCameraType || camType == camera.RightCameraType) {
			inputCh = camera.GetOutputChannel(camType)
		} else {
			inputCh = camera.GetChannel(camType)
		}

		// Get output channel to write to (if using input channel)
		var outputCh camera.ImageChannel
		if !useOutputChannel && (camType == camera.LeftCameraType || camType == camera.RightCameraType) {
			outputCh = camera.GetOutputChannel(camType)
		}

		// Stream images
		for {
			select {
			case <-ctx.Done():
				logger.Debug("stream context done", "reason", ctx.Err())

				return nil

			case <-ticker.C:
				// Only process on tick to control frame rate
				select {
				case img, ok := <-inputCh:
					if !ok {
						logger.Debug("camera channel closed")

						return nil
					}

					// Process the image
					err := processImg(img, w)
					if err != nil {
						logger.Error("error processing image", "err", err)

						return err
					}

					// Send to output channel if needed
					if outputCh != nil {
						select {
						case outputCh <- img:
						default:
							logger.Debug("output channel full, dropping frame")
						}
					}

				case <-ctx.Done():
					logger.Debug("stream context done while waiting for frame", "reason", ctx.Err())

					return nil

				default:
					// No image available, wait for next tick
				}
			}
		}
	}
}

// HandleLeftStream returns a handler for streaming the left camera.
func HandleLeftStream(w http.ResponseWriter, r *http.Request) error {
	return HandleCameraStream(camera.LeftCameraType, false)(w, r)
}

// HandleRightStream returns a handler for streaming the right camera.
func HandleRightStream(w http.ResponseWriter, r *http.Request) error {
	return HandleCameraStream(camera.RightCameraType, false)(w, r)
}

// HandleOutputStream returns a handler for streaming the output camera.
func HandleOutputStream(w http.ResponseWriter, r *http.Request) error {
	return HandleCameraStream(camera.OutputCameraType, false)(w, r)
}
