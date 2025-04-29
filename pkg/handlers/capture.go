package handlers

import (
	"encoding/json"
	"fmt"
	"image"
	"image/jpeg"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/conneroisu/steroscopic-hardware/pkg/camera"
)

// CameraSystem manages the stereoscopic camera system
type CameraSystem struct {
	stereoSystem   *camera.ZedBoardStereoSystem
	initialized    bool
	mutex          sync.Mutex
	imageSavePath  string
	leftImagePath  string
	rightImagePath string
	depthMapPath   string
	lastCapture    time.Time
	parameters     *Parameters
}

// NewCameraSystem creates a new camera system
func NewCameraSystem(
	leftPort, rightPort, imageSavePath string,
	params *Parameters,
) *CameraSystem {
	return &CameraSystem{
		stereoSystem:  camera.NewZedBoardStereoSystem(leftPort, rightPort),
		initialized:   false,
		imageSavePath: imageSavePath,
		parameters:    params,
	}
}

// Initialize initializes the camera system
func (cs *CameraSystem) Initialize() error {
	cs.mutex.Lock()
	defer cs.mutex.Unlock()

	if cs.initialized {
		return nil
	}

	// Create image directory if it doesn't exist
	if err := os.MkdirAll(cs.imageSavePath, 0755); err != nil {
		return fmt.Errorf("failed to create image directory: %v", err)
	}

	// Initialize the camera system
	if err := cs.stereoSystem.Initialize(); err != nil {
		return err
	}

	cs.initialized = true
	return nil
}

// Capture captures stereo images and processes the depth map
func (cs *CameraSystem) Capture() error {
	cs.mutex.Lock()
	defer cs.mutex.Unlock()

	if !cs.initialized {
		return fmt.Errorf("camera system not initialized")
	}

	// Capture stereo images
	leftImg, rightImg, err := cs.stereoSystem.CaptureStereoImages()
	if err != nil {
		return err
	}

	// Generate timestamp for filenames
	timestamp := time.Now().Format("20060102_150405")
	cs.leftImagePath = filepath.Join(cs.imageSavePath, fmt.Sprintf("left_%s.jpg", timestamp))
	cs.rightImagePath = filepath.Join(cs.imageSavePath, fmt.Sprintf("right_%s.jpg", timestamp))
	cs.depthMapPath = filepath.Join(cs.imageSavePath, fmt.Sprintf("depth_%s.png", timestamp))

	// Save images
	if err := saveJPEGImage(leftImg, cs.leftImagePath); err != nil {
		return fmt.Errorf("failed to save left image: %v", err)
	}

	if err := saveJPEGImage(rightImg, cs.rightImagePath); err != nil {
		return fmt.Errorf("failed to save right image: %v", err)
	}

	// Process depth map (this would call oyr SAD algorithm)
	print("Processing depth map")

	cs.lastCapture = time.Now()
	return nil
}

// Close closes the camera system
func (cs *CameraSystem) Close() error {
	cs.mutex.Lock()
	defer cs.mutex.Unlock()

	if !cs.initialized {
		return nil
	}

	err := cs.stereoSystem.Close()
	cs.initialized = false
	return err
}

// GetStatus returns the status of the camera system
func (cs *CameraSystem) GetStatus() map[string]any {
	cs.mutex.Lock()
	defer cs.mutex.Unlock()

	status := map[string]any{
		"initialized": cs.initialized,
		"leftCamera":  "disconnected",
		"rightCamera": "disconnected",
		"depthMap":    "Not available",
	}

	if cs.initialized {
		status["leftCamera"] = "connected"
		status["rightCamera"] = "connected"

		if cs.lastCapture.IsZero() {
			status["depthMap"] = "Not available"
		} else {
			status["depthMap"] = "Available"
		}
	}

	return status
}

// GetImagePaths returns the paths to the latest captured images
func (cs *CameraSystem) GetImagePaths() (string, string, string) {
	cs.mutex.Lock()
	defer cs.mutex.Unlock()

	return cs.leftImagePath, cs.rightImagePath, cs.depthMapPath
}

// saveJPEGImage saves an image as JPEG
func saveJPEGImage(img image.Image, path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	return jpeg.Encode(file, img, &jpeg.Options{Quality: 90})
}

// CameraHandler handles HTTP requests for camera operations
func CameraHandler(cameraSystem *CameraSystem) APIFn {
	return func(w http.ResponseWriter, r *http.Request) error {
		w.Header().Set("Content-Type", "application/json")

		// Initialize camera system if needed
		if !cameraSystem.initialized {
			if err := cameraSystem.Initialize(); err != nil {
				http.Error(w, fmt.Sprintf(`{"success":false,"message":"Failed to initialize camera system: %v"}`, err), http.StatusInternalServerError)
				return fmt.Errorf("failed to initialize camera system: %v", err)
			}
		}

		// Handle capture request
		if r.Method == http.MethodPost {
			if err := cameraSystem.Capture(); err != nil {
				http.Error(w, fmt.Sprintf(`{"success":false,"message":"Failed to capture images: %v"}`, err), http.StatusInternalServerError)
				return fmt.Errorf("failed to capture images: %v", err)
			}

			// Get image paths
			leftPath, rightPath, depthPath := cameraSystem.GetImagePaths()

			// Create response with image paths
			response := map[string]any{
				"success": true,
				"message": "Capture successful",
				"images": map[string]string{
					"left":  filepath.Base(leftPath),
					"right": filepath.Base(rightPath),
					"depth": filepath.Base(depthPath),
				},
			}

			err := json.NewEncoder(w).Encode(response)
			return err
		}

		// Handle GET request (status)
		status := cameraSystem.GetStatus()
		status["success"] = true
		err := json.NewEncoder(w).Encode(status)
		return err
	}
}

// GetStreamHandler returns a handler for streaming camera images
func GetStreamHandler(cameraSystem *CameraSystem, side string) APIFn {
	return func(w http.ResponseWriter, r *http.Request) error {
		// Set headers for MJPEG stream
		w.Header().Set("Content-Type", "multipart/x-mixed-replace; boundary=frame")

		// Continuously send frames
		for {
			select {
			case <-r.Context().Done():
				return nil
			default:
				// Continue
			}

			// Get the appropriate camera
			var camera *camera.ZedBoardCamera
			if side == "left" {
				camera = cameraSystem.stereoSystem.LeftCamera
			} else {
				camera = cameraSystem.stereoSystem.RightCamera
			}

			// Capture an image
			img, err := camera.Capture()
			if err != nil {
				log.Printf("Error capturing image from %s camera: %v", side, err)
				time.Sleep(500 * time.Millisecond)
				continue
			}

			// Write frame boundary
			fmt.Fprintf(w, "--frame\r\nContent-Type: image/jpeg\r\n\r\n")

			// Encode and write JPEG data
			if err := jpeg.Encode(w, img, nil); err != nil {
				log.Printf("Error encoding JPEG: %v", err)
				return fmt.Errorf("error encoding JPEG: %v", err)
			}

			// Add small delay between frames
			time.Sleep(100 * time.Millisecond)
		}
	}
}

// GetMapHandler is a handler for the websocket connection.
func GetMapHandler(cameraSystem *CameraSystem) APIFn {
	return func(w http.ResponseWriter, r *http.Request) error {
		// TODO: implement
		return nil
	}
}
