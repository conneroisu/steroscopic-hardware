package handlers

import (
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/conneroisu/steroscopic-hardware/pkg/camera"
	"github.com/conneroisu/steroscopic-hardware/pkg/homedir"
)

// HandleCameraStream is a generic handler for streaming camera images.
func HandleCameraStream(camType camera.Type) APIFn {
	return func(w http.ResponseWriter, _ *http.Request) error {
		// read $HOME/{type}.png
		dir, err := homedir.Dir()
		if err != nil {
			return err
		}
		f, err := os.Open(filepath.Join(dir, string(camType)+".png"))
		if err != nil {
			return err
		}
		defer f.Close()
		bdy, err := io.ReadAll(f)
		if err != nil {
			return err
		}
		_, err = w.Write(bdy)
		if err != nil {
			return err
		}

		return nil
	}
}

// HandleLeftStream returns a handler for streaming the left camera.
func HandleLeftStream(w http.ResponseWriter, r *http.Request) error {
	return HandleCameraStream(camera.LeftCameraType)(w, r)
}

// HandleRightStream returns a handler for streaming the right camera.
func HandleRightStream(w http.ResponseWriter, r *http.Request) error {
	return HandleCameraStream(camera.RightCameraType)(w, r)
}

// HandleOutputStream returns a handler for streaming the output camera.
func HandleOutputStream(w http.ResponseWriter, r *http.Request) error {
	return HandleCameraStream(camera.OutputCameraType)(w, r)
}
