package camera

import (
	"context"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"github.com/conneroisu/steroscopic-hardware/pkg/homedir"
)

// StaticCamera represents a camera that loads images from files.
type StaticCamera struct {
	BaseCamera
	path   string
	logger *slog.Logger
}

// NewStaticCamera creates a new static camera that reads from the specified file path.
func NewStaticCamera(ctx context.Context, path string, typ Type) *StaticCamera {
	return &StaticCamera{
		BaseCamera: NewBaseCamera(ctx, typ),
		path:       path,
		logger:     slog.Default().WithGroup(fmt.Sprintf("static-camera-%s", typ)),
	}
}

// Stream continuously reads the static image and sends it to the output channel.
func (sc *StaticCamera) Stream(ctx context.Context, outCh ImageChannel) {
	sc.logger.Info("starting static camera stream", "path", sc.path)
	defer sc.logger.Info("static camera stream stopped")

	// Create error channel for internal communication
	errChan := make(chan error, 1)

	// Set up ticker for a reasonable frame rate
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	// Backoff parameters
	initialBackoff := 10 * time.Millisecond
	maxBackoff := 1 * time.Second
	backoff := initialBackoff
	consecutiveFailures := 0

	for {
		select {
		case <-ctx.Done():
			return
		case <-sc.Context().Done():
			return
		case err := <-errChan:
			sc.logger.Error("error in static camera", "err", err)
			time.Sleep(1 * time.Second) // Delay before retry
		case <-ticker.C:
			if sc.IsPaused() {
				continue
			}

			// Load image file
			img, err := sc.loadImage()
			if err != nil {
				select {
				case errChan <- err:
				default:
					sc.logger.Error("error loading image", "err", err)
				}
				continue
			}

			// Send image to output channel
			select {
			case outCh <- img:
				sc.logger.Debug("image sent to channel")
				// Reset backoff on success
				backoff = initialBackoff
				consecutiveFailures = 0
			case <-ctx.Done():
				return
			case <-sc.Context().Done():
				return
			default:
				// Channel is full, apply backoff
				consecutiveFailures++

				// Only log at certain thresholds to prevent log spam
				if consecutiveFailures == 1 || consecutiveFailures%10 == 0 {
					sc.logger.Debug("output channel full, applying backoff",
						"consecutiveFailures", consecutiveFailures,
						"currentBackoff", backoff)
				}

				// Apply backoff delay
				time.Sleep(backoff)

				// Exponential backoff with a maximum cap
				backoff = min(time.Duration(float64(backoff)*1.5), maxBackoff)
			}
		}
	}
}

// loadImage reads and processes the image file.
func (sc *StaticCamera) loadImage() (*image.Gray, error) {
	// Check if file exists
	_, err := os.Stat(sc.path)
	if os.IsNotExist(err) {
		return nil, fmt.Errorf("image file not found: %s", sc.path)
	}

	// Determine image format based on extension
	ext := filepath.Ext(sc.path)
	var (
		grayImg *image.Gray
		file    *os.File
		img     image.Image
	)

	switch ext {
	case ".png":
		file, err = os.Open(sc.path)
		if err != nil {
			return nil, err
		}
		defer file.Close()

		img, err = png.Decode(file)
		if err != nil {
			return nil, err
		}
		bounds := img.Bounds()
		grayImg = image.NewGray(bounds)

		// Direct access to pixel data
		grayPix := grayImg.Pix
		stride := grayImg.Stride

		// Optimize by checking image type
		switch img := img.(type) {
		case *image.Gray:
			return img, nil
		case *image.RGBA:
			for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
				rowStart := (y - bounds.Min.Y) * stride
				for x := bounds.Min.X; x < bounds.Max.X; x++ {
					i := img.PixOffset(x, y)
					r := img.Pix[i]
					g := img.Pix[i+1]
					b := img.Pix[i+2]

					// Use integer arithmetic
					grayPix[rowStart+x-bounds.Min.X] = uint8((19595*uint32(r) +
						38470*uint32(g) +
						7471*uint32(b) + 1<<15) >> 24)
				}
			}
		default:
			for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
				rowStart := (y - bounds.Min.Y) * stride
				for x := bounds.Min.X; x < bounds.Max.X; x++ {
					r, g, b, _ := img.At(x, y).RGBA()
					grayPix[rowStart+x-bounds.Min.X] = uint8((19595*r +
						38470*g +
						7471*b + 1<<15) >> 24)
				}
			}
		}
	case ".jpg", ".jpeg":
		// Open the file
		file, err = os.Open(sc.path)
		if err != nil {
			return nil, fmt.Errorf("error opening image file: %w", err)
		}
		defer file.Close()

		// Decode the image
		img, _, err = image.Decode(file)
		if err != nil {
			return nil, fmt.Errorf("error decoding image: %w", err)
		}

		// Convert to grayscale
		bounds := img.Bounds()
		grayImg = image.NewGray(bounds)
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			for x := bounds.Min.X; x < bounds.Max.X; x++ {
				grayImg.Set(x, y, color.GrayModel.Convert(img.At(x, y)))
			}
		}
	default:
		return nil, fmt.Errorf("unsupported image format: %s", ext)
	}

	// Save a copy for debugging
	err = homedir.SaveImage(grayImg)
	if err != nil {
		sc.logger.Error("failed to save debug image", "err", err)
	}

	return grayImg, nil
}

// Close releases all resources.
func (sc *StaticCamera) Close() error {
	sc.logger.Info("closing static camera")
	sc.Cancel()
	return nil
}
