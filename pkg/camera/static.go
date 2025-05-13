package camera

import (
	"context"
	"fmt"
	"image"
	"image/color"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"github.com/conneroisu/steroscopic-hardware/pkg/despair"
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
			case <-ctx.Done():
				return
			case <-sc.Context().Done():
				return
			default:
				// If channel is full, we'll try again next tick
				sc.logger.Debug("output channel full, skipping frame")
				time.Sleep(10 * time.Millisecond) // Prevent busy-loop with a brief delay
			}
		}
	}
}

// loadImage reads and processes the image file.
func (sc *StaticCamera) loadImage() (*image.Gray, error) {
	// Check if file exists
	if _, err := os.Stat(sc.path); os.IsNotExist(err) {
		return nil, fmt.Errorf("image file not found: %s", sc.path)
	}

	// Determine image format based on extension
	ext := filepath.Ext(sc.path)
	var grayImg *image.Gray
	var err error

	switch ext {
	case ".png":
		grayImg, err = despair.LoadPNG(sc.path)
		if err != nil {
			return nil, fmt.Errorf("error loading PNG: %w", err)
		}
	case ".jpg", ".jpeg":
		// Open the file
		file, err := os.Open(sc.path)
		if err != nil {
			return nil, fmt.Errorf("error opening image file: %w", err)
		}
		defer file.Close()

		// Decode the image
		img, _, err := image.Decode(file)
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
