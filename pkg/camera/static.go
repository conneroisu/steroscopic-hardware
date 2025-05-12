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

// StaticCamera represents a ZedBoard camera.
type StaticCamera struct {
	Path   string
	ctx    context.Context
	cancel context.CancelFunc
}

// NewStaticCamera creates a new ZedBoard camera.
func NewStaticCamera(path string, ch chan *image.Gray) *StaticCamera {
	ctx, cancel := context.WithCancel(context.Background())
	sc := &StaticCamera{
		Path:   path,
		ctx:    ctx,
		cancel: cancel,
	}
	go sc.Stream(sc.ctx, ch)
	return sc
}

var _ Camer = (*StaticCamera)(nil)

// Stream streams the camera images to the given channel.
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
			slog.Debug("read image")
			if img == nil {
				continue
			}
			slog.Debug("sending image")
			outCh <- img
		}
	}
}

// Config returns the current configuration of the camera.
func (z *StaticCamera) Config() *Config { return &Config{} }

func (z *StaticCamera) read(errChan chan error) <-chan *image.Gray {
	slog.Debug("reading image: " + z.Path)
	mkdCh := make(chan *image.Gray, 1)
	img, err := z.getImage()
	if err != nil {
		slog.Error("failed to get image", "err", err)
		errChan <- fmt.Errorf("failed to get image: %v", err)
		return nil
	}
	mkdCh <- img
	return mkdCh
}

func (z *StaticCamera) getImage() (*image.Gray, error) {
	slog.Debug("getting image: " + z.Path)
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

// Close closes the camera.
func (z *StaticCamera) Close() error {
	return nil
}
