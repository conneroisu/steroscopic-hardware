package handlers

import (
	"context"
	"fmt"
	"image"
	"image/jpeg"
	"log"
	"net/http"

	"github.com/conneroisu/steroscopic-hardware/pkg/camera"
)

// StreamHandlerFn returns a handler for streaming camera images
func StreamHandlerFn(
	ctx context.Context,
	camera camera.Camer,
	outCh chan *image.Gray,
) APIFn {
	return func(w http.ResponseWriter, r *http.Request) error {
		// Set headers for MJPEG stream
		w.Header().Set("Content-Type", "multipart/x-mixed-replace; boundary=frame")

		var imgCh = make(chan *image.Gray)
		go camera.Stream(ctx, imgCh)
		// Continuously send frames
		for {
			select {
			case <-r.Context().Done():
				return nil
			case img := <-imgCh:
				for {
					select {
					case <-r.Context().Done():
						return nil
					case outCh <- img:
						// Write frame boundary
						fmt.Fprintf(w, "--frame\r\nContent-Type: image/jpeg\r\n\r\n")

						// Encode and write JPEG data
						if err := jpeg.Encode(w, img, nil); err != nil {
							log.Printf("Error encoding JPEG: %v", err)
							return fmt.Errorf("error encoding JPEG: %v", err)
						}
					}
				}
			}
		}
	}
}

// GetMapHandler is a handler for the websocket connection.
func GetMapHandler(leftOutCh, rightOutCh chan *image.Gray) APIFn {
	return func(w http.ResponseWriter, r *http.Request) error {
		// TODO: implement
		return nil
	}
}
