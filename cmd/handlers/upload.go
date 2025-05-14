package handlers

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"

	"github.com/conneroisu/steroscopic-hardware/pkg/camera"
)

// UploadHandler handles the upload of static image files for camera simulation.
func UploadHandler(appCtx context.Context, typ camera.Type) APIFn {
	logger := slog.Default().WithGroup(fmt.Sprintf("upload-handler-%s", typ))

	return func(w http.ResponseWriter, r *http.Request) error {
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

		logger.Info("file upload started", "filename", header.Filename, "type", typ)

		// Read file content
		body, err := io.ReadAll(file)
		if err != nil {
			return fmt.Errorf("failed to read file content: %w", err)
		}

		// Save file to temporary directory
		dir := os.TempDir()
		path := filepath.Join(dir, header.Filename)
		err = os.WriteFile(path, body, 0644)
		if err != nil {
			return fmt.Errorf("failed to save file: %w", err)
		}

		logger.Info("file saved", "path", path, "size", len(body))

		old := camera.GetCamera(typ)
		err = old.Close()
		if err != nil {
			return err
		}
		// Create static camera - Using the application context
		staticCam := camera.NewStaticCamera(appCtx, path, typ)

		// Set camera in manager - Using the application context
		err = camera.SetCamera(appCtx, typ, staticCam)
		if err != nil {
			return fmt.Errorf("failed to set static camera: %w", err)
		}

		logger.Info("camera configured from uploaded file", "type", typ)

		// Return success HTML for HTMX to replace the form
		successHTML := fmt.Sprintf(`
		<div id="%s-upload-form-container" class="space-y-2">
			<div class="flex items-center justify-between mb-2">
				<span class="text-sm text-green-400">Image uploaded successfully: %s</span>
			</div>
			<div class="mt-4">
				<div class="w-full bg-green-700 rounded-full h-2 mb-2">
					<div class="bg-green-500 h-2 rounded-full w-full"></div>
				</div>
			</div>
			<div class="flex justify-between mt-2 items-center">
				<span class="text-sm text-green-400">Camera now streaming from static image</span>
				<button 
					hx-get="/"
					hx-push-url="true"
					hx-target="#app"
					class="bg-blue-600 hover:bg-blue-700 text-white rounded px-3 py-1 text-sm"
				>
					Reload UI
				</button>
			</div>
		</div>`, string(typ), header.Filename)

		_, err = w.Write([]byte(successHTML))
		if err != nil {
			return fmt.Errorf("failed to write success HTML: %w", err)
		}

		return nil
	}
}
