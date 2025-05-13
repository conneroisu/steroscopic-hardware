package handlers

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/conneroisu/steroscopic-hardware/pkg/camera"
)

// UploadHandler handles the upload of static image files for camera simulation.
func UploadHandler(typ camera.Type) APIFn {
	return func(_ http.ResponseWriter, r *http.Request) error {
		// Parse multipart form
		if err := r.ParseMultipartForm(32 << 20); err != nil { // 32MB max
			return fmt.Errorf("failed to parse multipart form: %w", err)
		}

		// Get uploaded file
		file, header, err := r.FormFile("file")
		if err != nil {
			return fmt.Errorf("failed to get uploaded file: %w", err)
		}
		defer file.Close()

		// Read file content
		body, err := io.ReadAll(file)
		if err != nil {
			return fmt.Errorf("failed to read file content: %w", err)
		}

		// Save file to temporary directory
		dir := os.TempDir()
		path := dir + "/" + header.Filename
		if err := os.WriteFile(path, body, 0644); err != nil {
			return fmt.Errorf("failed to save file: %w", err)
		}

		// Create static camera
		staticCam := camera.NewStaticCamera(path, typ)

		// Set camera in manager
		if err := camera.SetCamera(r.Context(), typ, staticCam); err != nil {
			return fmt.Errorf("failed to set static camera: %w", err)
		}

		return nil
	}
}
