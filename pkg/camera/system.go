package camera

import (
	"fmt"
	"image"
	"sync"
)

// ZedBoardStereoSystem represents a stereoscopic system using two ZedBoard cameras
type ZedBoardStereoSystem struct {
	LeftCamera  *ZedBoardCamera
	RightCamera *ZedBoardCamera
}

// NewZedBoardStereoSystem creates a new ZedBoard stereoscopic camera system
func NewZedBoardStereoSystem(leftPortName, rightPortName string) *ZedBoardStereoSystem {
	leftConfig := DefaultZedBoardConfig(leftPortName)
	rightConfig := DefaultZedBoardConfig(rightPortName)

	return &ZedBoardStereoSystem{
		LeftCamera:  NewZedBoardCamera(leftConfig),
		RightCamera: NewZedBoardCamera(rightConfig),
	}
}

// Initialize initializes both ZedBoard cameras
func (zbs *ZedBoardStereoSystem) Initialize() error {
	leftErr := zbs.LeftCamera.Initialize()
	if leftErr != nil {
		return fmt.Errorf("failed to initialize left camera: %v", leftErr)
	}

	rightErr := zbs.RightCamera.Initialize()
	if rightErr != nil {
		// Close the left camera if the right one fails
		zbs.LeftCamera.Close()
		return fmt.Errorf("failed to initialize right camera: %v", rightErr)
	}

	return nil
}

// CaptureStereoImages captures images from both ZedBoard cameras
func (zbs *ZedBoardStereoSystem) CaptureStereoImages() (leftImg, rightImg image.Image, err error) {
	// Use a WaitGroup to ensure both captures complete
	var wg sync.WaitGroup
	var leftError, rightError error

	wg.Add(2)

	// Capture from left camera
	go func() {
		defer wg.Done()
		leftImg, leftError = zbs.LeftCamera.Capture()
	}()

	// Capture from right camera
	go func() {
		defer wg.Done()
		rightImg, rightError = zbs.RightCamera.Capture()
	}()

	// Wait for both captures to complete
	wg.Wait()

	// Check for errors
	if leftError != nil {
		return nil, nil, fmt.Errorf("left camera capture failed: %v", leftError)
	}

	if rightError != nil {
		return nil, nil, fmt.Errorf("right camera capture failed: %v", rightError)
	}

	return leftImg, rightImg, nil
}

// Close closes both ZedBoard cameras
func (zbs *ZedBoardStereoSystem) Close() error {
	leftErr := zbs.LeftCamera.Close()
	rightErr := zbs.RightCamera.Close()

	if leftErr != nil {
		return fmt.Errorf("failed to close left camera: %v", leftErr)
	}

	if rightErr != nil {
		return fmt.Errorf("failed to close right camera: %v", rightErr)
	}

	return nil
}
