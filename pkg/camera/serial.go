package camera

import (
	"context"
	"fmt"
	"image"
	"image/color"
	"log"
	"log/slog"
	"sync"
	"time"

	"github.com/conneroisu/steroscopic-hardware/pkg/homedir"
	"go.bug.st/serial"
)

var (
	// DefaultStartSeq is the default start marker for image data.
	DefaultStartSeq = []byte{0xff, 0xd8}
	// DefaultEndSeq is the default end marker for image data.
	DefaultEndSeq = []byte{0xff, 0xd9}
	// DefaultImageWidth is the default expected image width in pixels.
	DefaultImageWidth = 1920
	// DefaultImageHeight is the default expected image height in pixels.
	DefaultImageHeight = 1080

	_ Camer = (*SerialCamera)(nil)
)

type (
	// SerialCamera represents a camera connected via serial port.
	SerialCamera struct {
		mu          sync.Mutex
		ctx         context.Context
		cancel      context.CancelFunc
		port        serial.Port
		StartSeq    []byte // Byte sequence indicating start of image data
		EndSeq      []byte // Byte sequence indicating end of image data
		ImageWidth  int    // Expected image width in pixels
		ImageHeight int    // Expected image height in pixels
		logger      *slog.Logger
		config      Config
		OnClose     func()
		ch          chan *image.Gray
	}
)

// NewSerialCamera creates a new SerialCamera instance.
func NewSerialCamera(
	typ Type,
	portName string,
	baudRate int,
	compression int,
) (*SerialCamera, error) {
	var err error
	switch typ {
	case LeftCameraType:
		err = Left().Close()
		if err != nil {
			return nil, err
		}
	case RightCameraType:
		err = Right().Close()
		if err != nil {
			return nil, err
		}
	}
	ctx, cancel := context.WithCancel(context.Background())
	// Open the port
	sc := SerialCamera{
		ctx:         ctx,
		cancel:      cancel,
		StartSeq:    DefaultStartSeq,
		EndSeq:      DefaultEndSeq,
		ImageWidth:  DefaultImageWidth,
		ImageHeight: DefaultImageHeight,
		port:        nil,
		mu:          sync.Mutex{},
		config: Config{
			Port:        portName,
			BaudRate:    baudRate,
			Compression: compression,
		},
		logger: slog.Default().WithGroup("serial-port-" + string(typ)),
	}
	// Configure serial port
	mode := &serial.Mode{
		BaudRate: baudRate,
		DataBits: 8,
		Parity:   serial.NoParity,
		StopBits: serial.OneStopBit,
	}

	sc.logger.Info("opening serial port", "port", portName)
	sc.port, err = serial.Open(sc.config.Port, mode)
	if err != nil {
		return nil, fmt.Errorf("failed to open serial port %s: %v", portName, err)
	}
	// Set read timeout
	err = sc.port.SetReadTimeout(serial.NoTimeout)
	if err != nil {
		sc.port.Close()
		return nil, fmt.Errorf("failed to set read timeout: %v", err)
	}
	go sc.Stream(sc.ctx, sc.ch)

	return &sc, nil
}

// Config returns the current configuration of the camera.
func (sc *SerialCamera) Config() *Config { return &sc.config }

// Close closes the serial connection.
func (sc *SerialCamera) Close() error {
	sc.mu.Lock()
	defer sc.mu.Unlock()

	if sc.OnClose != nil {
		sc.OnClose()
	}
	sc.cancel()
	err := sc.port.Close()
	if err != nil {
		return err
	}
	return nil
}

// Stream reads images from the camera and sends them to the channel.
func (sc *SerialCamera) Stream(
	ctx context.Context,
	ch chan *image.Gray,
) {
	sc.logger.Debug("SerialCamera.Stream()")
	defer sc.logger.Debug("SerialCamera.Stream() done")

	sc.ch = ch
	var errChan = make(chan error, 1)
	readFn, err := sc.start(ctx, errChan, ch)
	if err != nil {
		sc.logger.Error("failed to read image data", "err", err)
		return
	}

	go readFn()

	for {
		select {
		case <-ctx.Done():
			sc.logger.Debug("context done, stopping read")
			return
		case <-sc.ctx.Done():
			sc.logger.Debug("inner context done, stopping read")
			return
		case err := <-errChan:
			sc.logger.Debug("error reading image", "err", err)
		}
	}
}

// readImageData reads image data from the serial port.
func (sc *SerialCamera) start(
	ctx context.Context,
	errChan chan error,
	imgCh chan *image.Gray,
) (func(), error) {
	var tries = 0
	for {
		sc.logger.Info("SerialCamera.read()")
		defer sc.logger.Info("SerialCamera.read() done")

		// Temporary read buffer
		tempBuf := make([]byte, 1024)

		// Send the start sequence
		_, err := sc.port.Write(sc.StartSeq)
		if err != nil {
			return nil, fmt.Errorf("failed to send start sequence: %v", err)
		}
		// After sending the start sequence, we should receive a 1-byte acknowledgement
		length, err := sc.port.Read(tempBuf)
		if err != nil {
			return nil, fmt.Errorf("failed to read acknowledgement: %v", err)
		}
		if length != 1 {
			// If we didn't receive a byte, try again by first sending the end sequence
			tries++
			if tries > 4 {
				return nil, fmt.Errorf("unexpected acknowledgement length: %d", length)
			}
			n, err := sc.port.Write(sc.EndSeq)
			if err != nil {
				return nil, fmt.Errorf("failed to send end sequence: %v", err)
			}
			if n != len(sc.EndSeq) {
				return nil, fmt.Errorf("failed to send end sequence: sent %d bytes, expected %d", n, len(sc.EndSeq))
			}
			continue
		}

		sc.OnClose = func() {
			sc.logger.Debug("sending end sequence")
			_, err := sc.port.Write(sc.EndSeq)
			if err != nil {
				log.Printf("failed to send end sequence: %v", err)
			}
		}

		return func() {
			for {
				select {
				case <-ctx.Done():
					sc.logger.Debug("context done, stopping read")
					return
				default:
					sc.readFn(ctx, errChan, imgCh)
				}
			}
		}, nil
	}
}

func (sc *SerialCamera) readFn(
	ctx context.Context,
	errChan chan error,
	imgCh chan *image.Gray,
) {
	sc.logger.Debug("----------------------------")
	defer sc.logger.Debug("----------------------------")
	var (
		tempBuf = []byte{}
		total   int
	)

	expectedLength := DefaultImageWidth * DefaultImageHeight

	sc.mu.Lock()
	defer sc.mu.Unlock()

	doneCh := make(chan struct{})
	go func() {
		start := time.Now()
		for {
			select {
			case <-time.After(time.Second * 30):
				// note to self that we are still working
				sc.logger.Info("still working", "time", time.Since(start))
			case <-ctx.Done():
				sc.logger.Debug("context done, stopping read")
				return
			case <-doneCh:
				return
			}
		}
	}()

	for {
		var buf = make([]byte, 100)
		length, err := sc.port.Read(buf)
		if err != nil {
			sc.logger.Error("error reading from serial port", "error", err)
			errChan <- fmt.Errorf("error reading from serial port: %v", err)
			return
		}
		tempBuf = append(tempBuf, buf[:length]...)
		total += length
		if total >= sc.ImageWidth*sc.ImageHeight {
			break
		}
	}
	doneCh <- struct{}{}
	sc.logger.Info("read", "total", total, "expected", expectedLength)

	img := image.NewGray(image.Rect(0, 0, sc.ImageWidth, sc.ImageHeight))

	// Grayscale format (1 byte per pixel)
	for y := range sc.ImageHeight { // y := 0; y < sc.config.ImageHeight; y++
		for x := range sc.ImageWidth { // x := 0; x < sc.config.ImageWidth; x++
			i := y*sc.ImageWidth + x
			gray := tempBuf[i]
			img.SetGray(x, y, color.Gray{Y: gray})
		}
	}

	// save image to homedir using pkg.homedir
	err := homedir.SaveImage(img)
	if err != nil {
		sc.logger.Error("failed to save image", "err", err)
	}

	select {
	case <-ctx.Done():
		sc.logger.Debug("context done, stopping read")
		return
	case imgCh <- img:
		sc.logger.Debug("image sent to channel")
	}
}
