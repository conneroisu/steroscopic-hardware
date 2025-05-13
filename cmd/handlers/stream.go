package handlers

import (
	"fmt"
	"image"
	"image/jpeg"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/conneroisu/steroscopic-hardware/pkg/camera"
)

const (
	streamTimeout = 60 * time.Minute
	frameRate     = time.Second / 10 // 10 FPS max
)

var encodeOpts = &jpeg.Options{Quality: 75}

func processImg(img *image.Gray, w io.Writer) error {
	// Write frame boundary to buffer
	_, err := fmt.Fprintf(w, "--frame\r\nContent-Type: image/jpeg\r\n\r\n")
	if err != nil {
		return err
	}

	// Encode image to buffer - consider using a worker pool for this
	err = jpeg.Encode(w, img, encodeOpts)
	if err != nil {
		log.Printf("Error encoding JPEG: %v", err)
		// Skip this frame instead of failing
		return nil
	}

	// Flush after writing complete frame
	f, ok := w.(http.Flusher)
	if ok {
		f.Flush()
	}

	return nil
}

// HandleLeftStream returns a handler for streaming the left camera image to a
// client.
func HandleLeftStream(w http.ResponseWriter, r *http.Request) error {
	// Set headers for MJPEG stream
	w.Header().Set("Content-Type", "multipart/x-mixed-replace; boundary=frame")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "close")
	w.Header().Set("Pragma", "no-cache")
	ticker := time.NewTicker(frameRate)
	defer ticker.Stop()
	for {
		select {
		case <-r.Context().Done():
			return nil
		case <-ticker.C:
			// Only process on tick to control frame rate
			select {
			case img, ok := <-camera.LeftCh():
				if !ok {
					return nil // Channel closed
				}
				err := processImg(img, w)
				if err != nil {
					return err
				}
				camera.LeftOutputCh() <- img
			case <-r.Context().Done():
				return nil
			}
		}
	}
}

// HandleRightStream returns a handler for streaming the right camera image to a
// client.
func HandleRightStream(w http.ResponseWriter, r *http.Request) error {
	// Set headers for MJPEG stream
	w.Header().Set("Content-Type", "multipart/x-mixed-replace; boundary=frame")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "close")
	w.Header().Set("Pragma", "no-cache")
	ticker := time.NewTicker(frameRate)
	defer ticker.Stop()
	for {
		select {
		case <-r.Context().Done():
			return nil
		case <-ticker.C:
			// Only process on tick to control frame rate
			select {
			case img, ok := <-camera.RightCh():
				if !ok {
					return nil // Channel closed
				}
				err := processImg(img, w)
				if err != nil {
					return err
				}
				camera.RightOutputCh() <- img
			case <-r.Context().Done():
				return nil
			}
		}
	}
}

// HandleOutputStream returns a handler for streaming the output camera image to a
// client.
func HandleOutputStream(w http.ResponseWriter, r *http.Request) error {
	// Set headers for MJPEG stream
	w.Header().Set("Content-Type", "multipart/x-mixed-replace; boundary=frame")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "close")
	w.Header().Set("Pragma", "no-cache")
	for {
		select {
		case <-r.Context().Done():
			return nil
		case img, ok := <-camera.OutputCh():
			if !ok {
				return nil // Channel closed
			}
			err := processImg(img, w)
			if err != nil {
				return err
			}
		}
	}
}
