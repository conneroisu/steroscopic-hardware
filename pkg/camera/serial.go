package camera

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"sync"
	"time"

	"go.bug.st/serial"
	"go.bug.st/serial/enumerator"
)

// SerialCameraConfig holds configuration for a camera connected via serial port
type SerialCameraConfig struct {
	PortName           string        // Serial port name (e.g., "COM3" on Windows, "/dev/ttyUSB0" on Linux)
	BaudRate           int           // Serial port baud rate (e.g., 115200)
	ImageWidth         int           // Expected image width in pixels
	ImageHeight        int           // Expected image height in pixels
	StartDelimiter     []byte        // Byte sequence indicating start of image data
	EndDelimiter       []byte        // Byte sequence indicating end of image data
	Timeout            time.Duration // Read timeout
	CaptureCommand     []byte        // Command to trigger image capture
	UseJPEGCompression bool          // Whether the camera sends JPEG compressed images
}

// SerialCamera represents a camera connected via serial port
type SerialCamera struct {
	config    SerialCameraConfig
	port      serial.Port
	isOpen    bool
	lastFrame []byte
	mutex     sync.Mutex
}

// DefaultStartDelimiter is the default start marker for image data
var DefaultStartDelimiter = []byte{0xFF, 0xD8} // JPEG SOI marker

// DefaultEndDelimiter is the default end marker for image data
var DefaultEndDelimiter = []byte{0xFF, 0xD9} // JPEG EOI marker

// DefaultCaptureCommand is the default command to trigger image capture
var DefaultCaptureCommand = []byte("CAPTURE\n")

// NewSerialCamera creates a new SerialCamera instance
func NewSerialCamera(config SerialCameraConfig) *SerialCamera {
	// Set defaults if not specified
	if config.StartDelimiter == nil {
		config.StartDelimiter = DefaultStartDelimiter
	}
	if config.EndDelimiter == nil {
		config.EndDelimiter = DefaultEndDelimiter
	}
	if config.CaptureCommand == nil {
		config.CaptureCommand = DefaultCaptureCommand
	}
	if config.Timeout == 0 {
		config.Timeout = 5 * time.Second
	}

	return &SerialCamera{
		config: config,
		isOpen: false,
	}
}

// Open opens the serial connection to the camera
func (sc *SerialCamera) Open() error {
	sc.mutex.Lock()
	defer sc.mutex.Unlock()

	if sc.isOpen {
		return nil
	}

	// Configure serial port
	mode := &serial.Mode{
		BaudRate: sc.config.BaudRate,
		DataBits: 8,
		Parity:   serial.NoParity,
		StopBits: serial.OneStopBit,
	}

	// Open the port
	var err error
	sc.port, err = serial.Open(sc.config.PortName, mode)
	if err != nil {
		return fmt.Errorf("failed to open serial port %s: %v", sc.config.PortName, err)
	}

	// Set read timeout
	err = sc.port.SetReadTimeout(sc.config.Timeout)
	if err != nil {
		sc.port.Close()
		return fmt.Errorf("failed to set read timeout: %v", err)
	}

	sc.isOpen = true
	return nil
}

// Close closes the serial connection
func (sc *SerialCamera) Close() error {
	sc.mutex.Lock()
	defer sc.mutex.Unlock()

	if !sc.isOpen {
		return nil
	}

	err := sc.port.Close()
	if err != nil {
		return err
	}

	sc.isOpen = false
	return nil
}

// CaptureImage triggers an image capture and reads the resulting image data
func (sc *SerialCamera) CaptureImage() (image.Image, error) {
	sc.mutex.Lock()
	defer sc.mutex.Unlock()

	if !sc.isOpen {
		return nil, fmt.Errorf("serial port is not open")
	}

	// Reset input buffer before sending command
	if err := sc.port.ResetInputBuffer(); err != nil {
		return nil, fmt.Errorf("failed to reset input buffer: %v", err)
	}

	// Send capture command
	_, err := sc.port.Write(sc.config.CaptureCommand)
	if err != nil {
		return nil, fmt.Errorf("failed to send capture command: %v", err)
	}

	// Read image data
	imageData, err := sc.readImageData()
	if err != nil {
		return nil, err
	}

	// Store the last frame
	sc.lastFrame = imageData

	// Convert to image
	if sc.config.UseJPEGCompression {
		// If the data is already in JPEG format, decode it
		img, err := jpeg.Decode(bytes.NewReader(imageData))
		if err != nil {
			return nil, fmt.Errorf("failed to decode JPEG data: %v", err)
		}
		return img, nil
	}

	// Otherwise, interpret as raw pixel data
	return sc.convertRawToImage(imageData)
}

// GetLastFrame returns the last captured frame as raw bytes
func (sc *SerialCamera) GetLastFrame() []byte {
	sc.mutex.Lock()
	defer sc.mutex.Unlock()

	if sc.lastFrame == nil {
		return nil
	}

	// Return a copy to avoid data races
	frameCopy := make([]byte, len(sc.lastFrame))
	copy(frameCopy, sc.lastFrame)
	return frameCopy
}

// readImageData reads image data from the serial port
func (sc *SerialCamera) readImageData() ([]byte, error) {
	// Buffer to store image data
	var buffer bytes.Buffer

	// Read data until timeout or end delimiter is found
	inImageData := false
	startDelimiter := sc.config.StartDelimiter
	endDelimiter := sc.config.EndDelimiter

	// Temporary read buffer
	tempBuf := make([]byte, 1024)

	for {
		n, err := sc.port.Read(tempBuf)
		if err != nil {
			return nil, fmt.Errorf("error reading from serial port: %v", err)
		}

		if n == 0 {
			// Timeout occurred
			if inImageData {
				// If we were in the middle of reading image data, we may have finished
				// (some cameras might not send an end delimiter)
				break
			}
			return nil, fmt.Errorf("timeout waiting for image data")
		}

		// Process the received data
		for i := 0; i < n; i++ {
			if !inImageData {
				// Look for start delimiter
				if i+len(startDelimiter) <= n && bytes.Equal(tempBuf[i:i+len(startDelimiter)], startDelimiter) {
					inImageData = true
					buffer.Write(startDelimiter)
					i += len(startDelimiter) - 1 // -1 because the loop will increment i
					continue
				}
			} else {
				// Look for end delimiter
				if i+len(endDelimiter) <= n && bytes.Equal(tempBuf[i:i+len(endDelimiter)], endDelimiter) {
					buffer.Write(endDelimiter)
					return buffer.Bytes(), nil
				}
				buffer.WriteByte(tempBuf[i])
			}
		}

		// Safety check for buffer size
		if buffer.Len() > sc.config.ImageWidth*sc.config.ImageHeight*3 {
			// If buffer size exceeds expected image size by a lot, something might be wrong
			return nil, fmt.Errorf("received data exceeds expected image size")
		}
	}

	// If we're here, we either timed out or didn't find an end delimiter
	if buffer.Len() == 0 {
		return nil, fmt.Errorf("no image data received")
	}

	return buffer.Bytes(), nil
}

// convertRawToImage converts raw pixel data to an image.Image
func (sc *SerialCamera) convertRawToImage(data []byte) (image.Image, error) {
	expectedSize := sc.config.ImageWidth * sc.config.ImageHeight

	// Check if we have reasonable data size for grayscale or RGB format
	if len(data) != expectedSize && len(data) != expectedSize*3 {
		return nil, fmt.Errorf("unexpected data size: got %d bytes, expected %d (grayscale) or %d (RGB)",
			len(data), expectedSize, expectedSize*3)
	}

	// Create a new RGBA image
	img := image.NewRGBA(image.Rect(0, 0, sc.config.ImageWidth, sc.config.ImageHeight))

	if len(data) == expectedSize {
		// Grayscale format (1 byte per pixel)
		for y := range sc.config.ImageHeight { // y := 0; y < sc.config.ImageHeight; y++
			for x := range sc.config.ImageWidth { // x := 0; x < sc.config.ImageWidth; x++
				i := y*sc.config.ImageWidth + x
				gray := data[i]
				img.Set(x, y, color.RGBA{gray, gray, gray, 255})
			}
		}
	} else {
		// RGB format (3 bytes per pixel)
		for y := range sc.config.ImageHeight { // y := 0; y < sc.config.ImageHeight; y++
			for x := range sc.config.ImageWidth { // x := 0; x < sc.config.ImageWidth; x++
				i := (y*sc.config.ImageWidth + x) * 3
				r, g, b := data[i], data[i+1], data[i+2]
				img.Set(x, y, color.RGBA{r, g, b, 255})
			}
		}
	}

	return img, nil
}

// FindCameras searches for available serial cameras
func FindCameras() ([]string, error) {
	ports, err := enumerator.GetDetailedPortsList()
	if err != nil {
		return nil, fmt.Errorf("failed to get serial ports list: %v", err)
	}

	var cameras []string
	for _, port := range ports {
		// Add ports with specific VID/PID for ZedBoard (if known)
		// or add all serial ports for the user to select from
		cameras = append(cameras, port.Name)
	}

	return cameras, nil
}

// StereoCameraSystem represents a system with two cameras for stereoscopic imaging
type StereoCameraSystem struct {
	LeftCamera  *SerialCamera
	RightCamera *SerialCamera
}

// NewStereoCameraSystem creates a new stereoscopic camera system
func NewStereoCameraSystem(leftConfig, rightConfig SerialCameraConfig) *StereoCameraSystem {
	return &StereoCameraSystem{
		LeftCamera:  NewSerialCamera(leftConfig),
		RightCamera: NewSerialCamera(rightConfig),
	}
}

// Open opens connections to both cameras
func (scs *StereoCameraSystem) Open() error {
	err := scs.LeftCamera.Open()
	if err != nil {
		return fmt.Errorf("failed to open left camera: %v", err)
	}

	err = scs.RightCamera.Open()
	if err != nil {
		// Close the left camera if the right one fails
		scs.LeftCamera.Close()
		return fmt.Errorf("failed to open right camera: %v", err)
	}

	return nil
}

// Close closes connections to both cameras
func (scs *StereoCameraSystem) Close() error {
	leftErr := scs.LeftCamera.Close()
	rightErr := scs.RightCamera.Close()

	if leftErr != nil {
		return fmt.Errorf("failed to close left camera: %v", leftErr)
	}

	if rightErr != nil {
		return fmt.Errorf("failed to close right camera: %v", rightErr)
	}

	return nil
}

// CaptureStereoImages captures images from both cameras as simultaneously as possible
func (scs *StereoCameraSystem) CaptureStereoImages() (leftImg, rightImg image.Image, err error) {
	// Use a WaitGroup to ensure both captures complete
	var wg sync.WaitGroup
	var leftError, rightError error

	wg.Add(2)

	// Capture from left camera
	go func() {
		defer wg.Done()
		leftImg, leftError = scs.LeftCamera.CaptureImage()
	}()

	// Capture from right camera
	go func() {
		defer wg.Done()
		rightImg, rightError = scs.RightCamera.CaptureImage()
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
