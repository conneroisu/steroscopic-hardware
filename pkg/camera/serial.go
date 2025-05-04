package camera

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/color"
	"log"
	"sync"

	"go.bug.st/serial"
)

// SerialCameraConfig holds configuration for a camera connected via serial port
type SerialCameraConfig struct {
	PortName       string // Serial port name (e.g., "COM3" on Windows, "/dev/ttyUSB0" on Linux)
	BaudRate       int    // Serial port baud rate (e.g., 115200)
	ImageWidth     int    // Expected image width in pixels
	ImageHeight    int    // Expected image height in pixels
	StartDelimiter []byte // Byte sequence indicating start of image data
	UseCompression bool   // Whether the camera sends lzma compressed data
}

// SerialCamera represents a camera connected via serial port
type SerialCamera struct {
	port           serial.Port
	mutex          sync.Mutex
	portID         string
	StartDelimiter []byte // Byte sequence indicating start of image data
	ImageWidth     int    // Expected image width in pixels
	ImageHeight    int    // Expected image height in pixels
}

// DefaultStartDelimiter is the default start marker for image data
var DefaultStartDelimiter = []byte{0xFF, 0xD8} // JPEG SOI marker

// NewSerialCamera creates a new SerialCamera instance
func NewSerialCamera(config SerialCameraConfig) (*SerialCamera, error) {
	// Set defaults if not specified
	if config.StartDelimiter == nil {
		config.StartDelimiter = DefaultStartDelimiter
	}
	sc := SerialCamera{
		StartDelimiter: config.StartDelimiter,
		ImageWidth:     config.ImageWidth,
		ImageHeight:    config.ImageHeight,
		port:           nil,
		mutex:          sync.Mutex{},
		portID:         config.PortName,
	}
	sc.mutex.Lock()
	defer sc.mutex.Unlock()

	// Configure serial port
	mode := &serial.Mode{
		BaudRate: config.BaudRate,
		DataBits: 8,
		Parity:   serial.NoParity,
		StopBits: serial.OneStopBit,
	}

	// Open the port
	var err error
	sc.port, err = serial.Open(sc.portID, mode)
	if err != nil {
		return nil, fmt.Errorf("failed to open serial port %s: %v", config.PortName, err)
	}

	// Set read timeout
	err = sc.port.SetReadTimeout(serial.NoTimeout)
	if err != nil {
		sc.port.Close()
		return nil, fmt.Errorf("failed to set read timeout: %v", err)
	}

	return &sc, nil
}

// Close closes the serial connection
func (sc *SerialCamera) Close() error {
	sc.mutex.Lock()
	defer sc.mutex.Unlock()

	err := sc.port.Close()
	if err != nil {
		return err
	}

	return nil
}

// Stream reads images from the camera and sends them to the channel
func (sc *SerialCamera) Stream(
	ctx context.Context,
	ch chan *image.Gray,
) {
	var errChan = make(chan error, 1)
	for {
		select {
		case <-ctx.Done():
			return
		case img := <-sc.read(errChan):
			if img == nil {
				continue
			}
			ch <- img
		case err := <-errChan:
			log.Printf("Error reading image: %v", err)
		}
	}
}

func (sc *SerialCamera) read(errChan chan error) <-chan *image.Gray {
	mkdCh := make(chan *image.Gray, 1)
	data, err := sc.readImageData()
	if err != nil {
		errChan <- fmt.Errorf("failed to read image data: %v", err)
		return nil
	}
	img, err := sc.convertRawToImage(data)
	if err != nil {
		errChan <- fmt.Errorf("failed to convert raw data to image: %v", err)
		return nil
	}
	mkdCh <- img
	return mkdCh
}

// readImageData reads image data from the serial port
func (sc *SerialCamera) readImageData() ([]byte, error) {
	// Buffer to store image data
	var buffer bytes.Buffer

	// Read data until timeout or end delimiter is found
	inImageData := false
	startDelimiter := sc.StartDelimiter

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
			}
		}

		// Safety check for buffer size
		if buffer.Len() > sc.ImageWidth*sc.ImageHeight*3 {
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
func (sc *SerialCamera) convertRawToImage(data []byte) (*image.Gray, error) {
	expectedSize := sc.ImageWidth * sc.ImageHeight

	// Check if we have reasonable data size for grayscale or RGB format
	if len(data) != expectedSize && len(data) != expectedSize*3 {
		return nil, fmt.Errorf("unexpected data size: got %d bytes, expected %d (grayscale) or %d (RGB)",
			len(data), expectedSize, expectedSize*3)
	}

	// Create a new RGBA image
	img := image.NewGray(image.Rect(0, 0, sc.ImageWidth, sc.ImageHeight))

	if len(data) == expectedSize {
		// Grayscale format (1 byte per pixel)
		for y := range sc.ImageHeight { // y := 0; y < sc.config.ImageHeight; y++
			for x := range sc.ImageWidth { // x := 0; x < sc.config.ImageWidth; x++
				i := y*sc.ImageWidth + x
				gray := data[i]
				img.Set(x, y, color.RGBA{gray, gray, gray, 255})
			}
		}
	} else {
		// RGB format (3 bytes per pixel)
		for y := range sc.ImageHeight { // y := 0; y < sc.config.ImageHeight; y++
			for x := range sc.ImageWidth { // x := 0; x < sc.config.ImageWidth; x++
				i := (y*sc.ImageWidth + x) * 3
				r, g, b := data[i], data[i+1], data[i+2]
				img.Set(x, y, color.RGBA{r, g, b, 255})
			}
		}
	}

	return img, nil
}
