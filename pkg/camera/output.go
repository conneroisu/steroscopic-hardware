package camera

import (
	"context"
	"image"
	"log/slog"
	"time"

	"github.com/conneroisu/steroscopic-hardware/pkg/despair"
)

var _ Camer = (*OutputCamera)(nil)

const defaultNumWorkers = 32

// OutputCamera represents a the camera output of a sad output.
type OutputCamera struct {
	Params   *despair.Parameters
	InputCh  chan<- despair.InputChunk  // InputCh sends input chunks to sad.
	OutputCh <-chan despair.OutputChunk // OutputCh receives output chunks from sad algo.
	logger   *slog.Logger
}

// NewOutputCamera creates a new OutputCamera.
func NewOutputCamera(params *despair.Parameters) *OutputCamera {
	oC := &OutputCamera{
		logger: slog.Default().WithGroup("output-camera"),
	}
	oC.InputCh, oC.OutputCh = despair.SetupConcurrentSAD(
		params, defaultNumWorkers,
	)
	return oC
}

// :GoImpl o *OutputCamera camera.Camer

// Stream streams the output "camera", the sad output.
func (o *OutputCamera) Stream(
	ctx context.Context,
	outCh chan *image.Gray,
) {
	for {
		select {
		case <-ctx.Done():
			slog.Debug("stopping stream")
			return
		case img := <-o.read():
			if img == nil {
				continue
			}
			outCh <- img
		}
	}
}

// Config returns the current configuration of the output camera.
//
// It is not used, but is required by the Camer interface.
func (o *OutputCamera) Config() *Config { return &Config{} }

func (o *OutputCamera) read() <-chan *image.Gray {
	mkdCh := make(chan *image.Gray, 1)
	leftImg := <-LeftOutputCh()
	rightImg := <-RightOutputCh()

	chunkSize := max(1, leftImg.Rect.Dy()/(defaultNumWorkers*4))
	numChunks := (leftImg.Rect.Dy() + chunkSize - 1) / chunkSize

	start := time.Now()
	for y := leftImg.Rect.Min.Y; y < leftImg.Rect.Max.Y; y += chunkSize {
		o.InputCh <- despair.InputChunk{
			Left:  leftImg,
			Right: rightImg,
			Region: image.Rect(
				leftImg.Rect.Min.X,
				y,
				leftImg.Rect.Max.X,
				min(y+chunkSize, leftImg.Rect.Max.Y),
			),
		}
	}
	got := despair.AssembleDisparityMap(o.OutputCh, leftImg.Rect, numChunks)
	end := time.Now()
	o.logger.Debug("Elapsed time", "took", end.Sub(start))
	mkdCh <- got
	return mkdCh
}

// Close closes the output camera.
func (o *OutputCamera) Close() error {
	close(o.InputCh)
	return nil
}
