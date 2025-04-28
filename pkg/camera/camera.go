package camera

import (
	"fmt"
	"image"
	"log"
	"time"
)

// ZedBoardCamera implements a camera interface for ZedBoard FPGAs
type ZedBoardCamera struct {
	SerialCamera *SerialCamera
	Ready        bool
}

// ZedBoardCameraConfig contains ZedBoard specific configuration
type ZedBoardCameraConfig struct {
	SerialConfig SerialCameraConfig
	InitCommands [][]byte // Commands to initialize the camera
}

// DefaultZedBoardConfig returns a default configuration for ZedBoard cameras
func DefaultZedBoardConfig(portName string) ZedBoardCameraConfig {
	return ZedBoardCameraConfig{
		SerialConfig: SerialCameraConfig{
			PortName:           portName,
			BaudRate:           115200,             // ZedBoard typically uses this baud rate
			ImageWidth:         640,                // Adjust based on your camera
			ImageHeight:        480,                // Adjust based on your camera
			StartDelimiter:     []byte{0xAA, 0xBB}, // Custom start delimiter for ZedBoard
			EndDelimiter:       []byte{0xCC, 0xDD}, // Custom end delimiter for ZedBoard
			Timeout:            5 * time.Second,
			CaptureCommand:     []byte("CAPTURE\n"),
			UseJPEGCompression: true, // Adjust based on your implementation
		},
		InitCommands: [][]byte{
			[]byte("INIT\n"),               // Command to initialize camera
			[]byte("RESOLUTION 640x480\n"), // Set resolution
			[]byte("FORMAT JPEG\n"),        // Set format to JPEG
		},
	}
}

// NewZedBoardCamera creates a new ZedBoard camera instance
func NewZedBoardCamera(config ZedBoardCameraConfig) *ZedBoardCamera {
	return &ZedBoardCamera{
		SerialCamera: NewSerialCamera(config.SerialConfig),
		Ready:        false,
	}
}

// Initialize opens the connection and sends initialization commands
func (zbc *ZedBoardCamera) Initialize() error {
	err := zbc.SerialCamera.Open()
	if err != nil {
		return fmt.Errorf("failed to open ZedBoard camera: %v", err)
	}

	// Send initialization commands
	for _, cmd := range zbc.GetConfig().InitCommands {
		_, err := zbc.SerialCamera.port.Write(cmd)
		if err != nil {
			zbc.SerialCamera.Close()
			return fmt.Errorf("failed to send init command: %v", err)
		}

		// Wait for acknowledgment response (optional)
		// This depends on your ZedBoard firmware implementation
		buffer := make([]byte, 64)
		_, err = zbc.SerialCamera.port.Read(buffer)
		if err != nil {
			// Just log this, don't fail initialization if no response
			log.Printf("No response to init command %q: %v", cmd, err)
		}

		// Allow time for command to process
		time.Sleep(100 * time.Millisecond)
	}

	zbc.Ready = true
	return nil
}

// GetConfig returns the ZedBoard camera configuration
func (zbc *ZedBoardCamera) GetConfig() ZedBoardCameraConfig {
	return ZedBoardCameraConfig{
		SerialConfig: zbc.SerialCamera.config,
		InitCommands: [][]byte{}, // We don't store this after initialization
	}
}

// Capture captures an image from the ZedBoard camera
func (zbc *ZedBoardCamera) Capture() (image.Image, error) {
	if !zbc.Ready {
		return nil, fmt.Errorf("ZedBoard camera is not initialized")
	}

	return zbc.SerialCamera.CaptureImage()
}

// Close closes the ZedBoard camera connection
func (zbc *ZedBoardCamera) Close() error {
	zbc.Ready = false
	return zbc.SerialCamera.Close()
}
