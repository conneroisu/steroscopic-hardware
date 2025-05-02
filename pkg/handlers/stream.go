package handlers

import (
	"fmt"
	"image"
	"image/jpeg"
	"log"
	"net/http"
	"time"

	"github.com/conneroisu/steroscopic-hardware/pkg/camera"
	"github.com/conneroisu/steroscopic-hardware/pkg/logger"
	datastar "github.com/starfederation/datastar/sdk/go"
)

// StreamHandlerFn returns a handler for streaming camera images to multiple clients
func StreamHandlerFn(manager *camera.StreamManager) APIFn {
	// Make sure manager is running
	manager.Start()

	return func(w http.ResponseWriter, r *http.Request) error {
		// Set headers for MJPEG stream
		w.Header().Set("Content-Type", "multipart/x-mixed-replace; boundary=frame")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "close")
		w.Header().Set("Pragma", "no-cache")

		clientChan := make(chan *image.Gray, 10) // Buffer a few frames

		manager.Register <- clientChan

		defer func() {
			manager.Unregister <- clientChan
		}()

		for {
			select {
			case <-r.Context().Done():
				return nil

			case img, ok := <-clientChan:
				if !ok {
					// Channel closed
					return nil
				}

				// Write frame boundary
				fmt.Fprintf(w, "--frame\r\nContent-Type: image/jpeg\r\n\r\n")

				// Encode and write JPEG data
				if err := jpeg.Encode(w, img, nil); err != nil {
					log.Printf("Error encoding JPEG: %v", err)
					return fmt.Errorf("error encoding JPEG: %v", err)
				}

				// Flush
				if f, ok := w.(http.Flusher); ok {
					f.Flush()
				}
			}
		}
	}
}

// LogHandler returns a handler for streaming logs to the browser console
func LogHandler() APIFn {
	ch := logger.Default()
	return func(w http.ResponseWriter, r *http.Request) error {
		sse := datastar.NewSSE(w, r)
		for {
			select {
			case <-r.Context().Done():
				return nil
			case log := <-ch:
				println("logging using slog")
				// err := sse.ConsoleLogf("%s %s - %s", log.Level, log.Time.Format(time.RFC3339), log.Message)
				// if err != nil {
				// 	slog.Error("failed to log to browser console", "err", err.Error())
				// }
				err := sse.ConsoleError(fmt.Errorf(
					"%s %s - %s",
					log.Level,
					log.Time.Format(time.RFC3339),
					log.Message,
				), datastar.WithExecuteScriptAutoRemove(false))
				if err != nil {
					return err
				}
			}
		}
	}
}
