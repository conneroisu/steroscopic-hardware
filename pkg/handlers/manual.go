package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/conneroisu/steroscopic-hardware/pkg/despair"
)

// ManualCalcDepthMapHandler is a handler for the manual depth map calculation endpoint.
func ManualCalcDepthMapHandler(cameraSystem *CameraSystem) APIFn {
	return func(w http.ResponseWriter, r *http.Request) error {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return fmt.Errorf("method not allowed")
		}

		// Set content type for JSON response
		w.Header().Set("Content-Type", "application/json")

		// Create temp directory for uploaded files if it doesn't exist
		uploadDir := filepath.Join(cameraSystem.imageSavePath, "uploads")
		if err := os.MkdirAll(uploadDir, 0755); err != nil {
			return fmt.Errorf("failed to create upload directory: %v", err)
		}

		// The request should have multipart/form-data content type
		if err := r.ParseMultipartForm(10 << 20); err != nil {
			return fmt.Errorf("failed to parse multipart form: %v", err)
		}

		// Get the parameters from the form
		blockSize := cameraSystem.parameters.BlockSize
		maxDisparity := cameraSystem.parameters.MaxDisparity

		// Try to get parameters from the form if provided
		if val := r.FormValue("blockSize"); val != "" {
			var err error
			blockSize, err = parseIntParam(val, 3, 31)
			if err != nil {
				return fmt.Errorf("invalid blockSize parameter: %v", err)
			}
		}

		if val := r.FormValue("maxDisparity"); val != "" {
			var err error
			maxDisparity, err = parseIntParam(val, 16, 256)
			if err != nil {
				return fmt.Errorf("invalid maxDisparity parameter: %v", err)
			}
		}

		// Get uploaded files
		leftFile, leftHeader, err := r.FormFile("leftImage")
		if err != nil {
			return fmt.Errorf("failed to get left image: %v", err)
		}
		defer leftFile.Close()

		rightFile, rightHeader, err := r.FormFile("rightImage")
		if err != nil {
			return fmt.Errorf("failed to get right image: %v", err)
		}
		defer rightFile.Close()

		// Generate timestamp for filenames
		timestamp := time.Now().Format("20060102_150405")
		leftImagePath := filepath.Join(uploadDir, fmt.Sprintf("left_%s%s", timestamp, filepath.Ext(leftHeader.Filename)))
		rightImagePath := filepath.Join(uploadDir, fmt.Sprintf("right_%s%s", timestamp, filepath.Ext(rightHeader.Filename)))
		depthMapPath := filepath.Join(uploadDir, fmt.Sprintf("depth_%s.png", timestamp))

		// Save uploaded files
		err = saveUploadedFile(leftFile, leftImagePath)
		if err != nil {
			return fmt.Errorf("failed to save left image: %v", err)
		}

		err = saveUploadedFile(rightFile, rightImagePath)
		if err != nil {
			return fmt.Errorf("failed to save right image: %v", err)
		}

		// Process the depth map using the despair package
		err = despair.RunSadPaths(leftImagePath, rightImagePath, blockSize, maxDisparity)
		if err != nil {
			return fmt.Errorf("failed to generate depth map: %v", err)
		}

		// Get the resulting depth map path (the despair package saves it to a fixed location)
		// We need to move it to our desired location
		err = os.Rename("disparity_map.png", depthMapPath)
		if err != nil {
			return fmt.Errorf("failed to move depth map: %v", err)
		}

		// Return the path to the depth map image
		response := map[string]interface{}{
			"success":     true,
			"message":     "Depth map generated successfully",
			"depthMapUrl": "/static/images/uploads/" + filepath.Base(depthMapPath),
		}

		return json.NewEncoder(w).Encode(response)
	}
}

// Helper function to save an uploaded file
func saveUploadedFile(file io.Reader, destinationPath string) error {
	out, err := os.Create(destinationPath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, file)
	return err
}

// Helper function to parse integer parameters with validation
func parseIntParam(val string, min, max int) (int, error) {
	var result int
	if _, err := fmt.Sscanf(val, "%d", &result); err != nil {
		return 0, err
	}
	if result < min || result > max {
		return 0, fmt.Errorf("value must be between %d and %d", min, max)
	}
	return result, nil
}
