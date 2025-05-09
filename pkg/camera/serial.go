package camera

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/color"
	"log"
	"log/slog"
	"sync"

	"go.bug.st/serial"
)

var (
	// DefaultStartSeq is the default start marker for image data
	DefaultStartSeq = []byte{0xff, 0xd8}
	// DefaultEndSeq is the default end marker for image data
	DefaultEndSeq = []byte{0xff, 0xd9}
	// DefaultImageWidth is the default expected image width in pixels
	DefaultImageWidth = 1920
	// DefaultImageHeight is the default expected image height in pixels
	DefaultImageHeight       = 1080
	_                  Camer = (*SerialCamera)(nil)
)

type (
	// SerialCamera represents a camera connected via serial port
	SerialCamera struct {
		mu             sync.Mutex
		ctx            context.Context
		cancel         context.CancelFunc
		port           serial.Port
		portID         string
		StartSeq       []byte // Byte sequence indicating start of image data
		EndSeq         []byte // Byte sequence indicating end of image data
		ImageWidth     int    // Expected image width in pixels
		ImageHeight    int    // Expected image height in pixels
		baudRate       int
		useCompression bool
		OnClose        func()
		ch             chan *image.Gray
	}

	// SerialCameraOption is a function that configures a SerialCamera.
	SerialCameraOption func(*SerialCamera)
)

// NewSerialCamera creates a new SerialCamera instance
func NewSerialCamera(
	portName string,
	baudRate int,
	useCompression bool,
	opts ...SerialCameraOption,
) (*SerialCamera, error) {
	ctx, cancel := context.WithCancel(context.Background())
	sc := SerialCamera{
		ctx:            ctx,
		cancel:         cancel,
		StartSeq:       DefaultStartSeq,
		EndSeq:         DefaultEndSeq,
		ImageWidth:     DefaultImageWidth,
		ImageHeight:    DefaultImageHeight,
		port:           nil,
		mu:             sync.Mutex{},
		portID:         portName,
		baudRate:       baudRate,
		useCompression: useCompression,
	}
	// for _, opt := range opts {
	// 	opt(&sc)
	// }

	// Configure serial port
	mode := &serial.Mode{
		BaudRate: baudRate,
		DataBits: 8,
		Parity:   serial.NoParity,
		StopBits: serial.OneStopBit,
	}

	// Open the port
	var err error
	sc.port, err = serial.Open(sc.portID, mode)
	if err != nil {
		return nil, fmt.Errorf("failed to open serial port %s: %v", portName, err)
	}

	// Set read timeout
	err = sc.port.SetReadTimeout(serial.NoTimeout)
	if err != nil {
		sc.port.Close()
		return nil, fmt.Errorf("failed to set read timeout: %v", err)
	}

	return &sc, nil
}

// WithStartSeq sets the start sequence for the serial camera.
func WithStartSeq(startSeq []byte) SerialCameraOption {
	return func(sc *SerialCamera) { sc.StartSeq = startSeq }
}

// WithEndSeq sets the end sequence for the serial camera.
func WithEndSeq(endSeq []byte) SerialCameraOption {
	return func(sc *SerialCamera) { sc.EndSeq = endSeq }
}

// Close closes the serial connection
func (sc *SerialCamera) Close() error {
	sc.mu.Lock()
	defer sc.mu.Unlock()

	if sc.OnClose != nil {
		sc.OnClose()
	}
	err := sc.port.Close()
	if err != nil {
		return err
	}

	sc.cancel()

	return nil
}

// Stream reads images from the camera and sends them to the channel
func (sc *SerialCamera) Stream(
	ctx context.Context,
	ch chan *image.Gray,
) {
	sc.mu.Lock()
	defer sc.mu.Unlock()
	sc.ch = ch
	var errChan = make(chan error, 1)
	readFn, err := sc.read(ctx, errChan, ch)
	if err != nil {
		slog.Error("failed to read image data", "err", err)
		return
	}
	go readFn()
	for {
		select {
		case <-ctx.Done():
			return
		case <-sc.ctx.Done():
			return
		case img := <-ch:
			if img == nil {
				continue
			}
			ch <- img
		case err := <-errChan:
			log.Printf("Error reading image: %v", err)
		}
	}
}

// readImageData reads image data from the serial port
func (sc *SerialCamera) read(
	ctx context.Context,
	errChan chan error,
	imgCh chan *image.Gray,
) (func(), error) {
	// Buffer to store image data
	var buffer bytes.Buffer

	// Temporary read buffer
	tempBuf := make([]byte, 1024)

	// Send the start sequence
	_, err := sc.port.Write(sc.StartSeq)
	if err != nil {
		return nil, fmt.Errorf("failed to send start sequence: %v", err)
	}
	slog.Debug("sent start sequence", "seq", sc.StartSeq)
	// After sending the start sequence, we should receive a 1-byte acknowledgement
	bit, err := sc.port.Read(tempBuf)
	if err != nil {
		return nil, fmt.Errorf("failed to read acknowledgement: %v", err)
	}
	if bit != 1 {
		return nil, fmt.Errorf("unexpected acknowledgement byte: %d", bit)
	}
	slog.Debug("received acknowledgement", "byte", tempBuf[0])
	sc.OnClose = func() {
		_, err := sc.port.Write(sc.EndSeq)
		if err != nil {
			log.Printf("failed to send end sequence: %v", err)
		}
	}

	return func() {
		for {
			sc.mu.Lock()
			slog.Debug("reading image data")

			_, err := sc.port.Read(tempBuf)
			if err != nil {
				errChan <- fmt.Errorf("error reading from serial port: %v", err)
			}

			// Safety check for buffer size
			if buffer.Len() > sc.ImageWidth*sc.ImageHeight {
				errChan <- fmt.Errorf("received data exceeds expected image size")
			}

			img, err := sc.convertRawToImage(tempBuf)
			if err != nil {
				errChan <- fmt.Errorf("failed to convert raw data to image: %v", err)
			}
			select {
			case <-ctx.Done():
				return
			case imgCh <- img:
			}
			slog.Debug("image data read successfully", "size", buffer.Len())
			sc.mu.Unlock()
		}
	}, nil
}

// convertRawToImage converts raw pixel data to an image.Image
func (sc *SerialCamera) convertRawToImage(
	data []byte,
) (*image.Gray, error) {
	expectedSize := sc.ImageWidth * sc.ImageHeight

	// Check if we have reasonable data size for grayscale or RGB format
	if len(data) != expectedSize && len(data) != expectedSize*3 {
		return nil, fmt.Errorf(
			"unexpected data size: got %d bytes, expected %d (grayscale) or %d (RGB)",
			len(data),
			expectedSize,
			expectedSize*3,
		)
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
		return nil, fmt.Errorf(
			"unexpected data size: got %d bytes, expected %d (grayscale)",
			len(data),
			expectedSize,
		)
	}

	return img, nil
}
