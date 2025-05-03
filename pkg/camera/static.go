package camera

import (
	"context"
	"fmt"
	"image"
	"image/color"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/conneroisu/steroscopic-hardware/pkg/despair"
)

// StaticCamera represents a ZedBoard camera
type StaticCamera struct {
	Path string
}

// NewStaticCamera creates a new ZedBoard camera
func NewStaticCamera(path string) *StaticCamera {
	return &StaticCamera{
		Path: path,
	}
}

var _ Camer = (*StaticCamera)(nil)

// ID returns the camera ID
func (z *StaticCamera) ID() string {
	return filepath.Clean(z.Path)
}

// Stream streams the camera
func (z *StaticCamera) Stream(ctx context.Context, outCh chan *image.Gray) {
	var errChan = make(chan error, 1)
	for {
		select {
		case <-ctx.Done():
			slog.Debug("stopping stream")
			return
		case err := <-errChan:
			slog.Error("Error reading image", "err", err)
			return
		case img := <-z.read(errChan):
			if img == nil {
				continue
			}
			outCh <- img
		}
	}
}

func (z *StaticCamera) read(errChan chan error) <-chan *image.Gray {
	mkdCh := make(chan *image.Gray, 1)
	img, err := z.getImage()
	if err != nil {
		errChan <- fmt.Errorf("failed to get image: %v", err)
		return nil
	}
	mkdCh <- img
	return mkdCh
}

func (z *StaticCamera) getImage() (*image.Gray, error) {
	// Open the image file
	file, err := os.Open(z.Path)
	if err != nil {
		return nil, fmt.Errorf("error opening image file: %w", err)
	}
	defer file.Close()

	var grayImg *image.Gray
	ext := filepath.Ext(z.Path)
	switch ext {
	case ".png":
		grayImg, err = despair.LoadPNG(z.Path)
		if err != nil {
			return nil, err
		}
	case ".jpg", ".jpeg":
		// Decode the image
		img, _, err := image.Decode(file)
		if err != nil {
			return nil, fmt.Errorf("error decoding image (%s): %w", z.Path, err)
		}
		// Create a new grayscale image with the same dimensions
		bounds := img.Bounds()
		grayImg = image.NewGray(bounds)

		// Convert to grayscale
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			for x := bounds.Min.X; x < bounds.Max.X; x++ {
				grayImg.Set(x, y, color.GrayModel.Convert(img.At(x, y)))
			}
		}
	default:
		return nil, fmt.Errorf("unsupported image format: %s", ext)
	}

	return grayImg, nil
}

// Close closes the camera
func (z *StaticCamera) Close() error {
	panic("not implemented") // TODO: Implement
}
