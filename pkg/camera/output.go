package camera

import (
	"context"
	"image"
	"log/slog"
	"time"

	"github.com/conneroisu/steroscopic-hardware/pkg/despair"
	"github.com/conneroisu/steroscopic-hardware/pkg/logger"
)

// OutputCamera represents a the camera output of a sad output.
type OutputCamera struct {
	Params        *despair.Parameters
	LeftClientCh  chan *image.Gray
	RightClientCh chan *image.Gray
	InputCh       chan<- despair.InputChunk  // InputCh sends input chunks to sad.
	OutputCh      <-chan despair.OutputChunk // OutputCh receives output chunks from sad algo.
	Left          *StreamManager
	Right         *StreamManager
	logger        *logger.Logger
}

var _ Camer = (*OutputCamera)(nil)

const defaultNumWorkers = 32

// NewOutputCamera creates a new OutputCamera
func NewOutputCamera(
	logger *logger.Logger,
	params *despair.Parameters,
	left, right *StreamManager,
) *OutputCamera {
	oC := &OutputCamera{
		Left:   left,
		Right:  right,
		Params: params,
		logger: logger,
	}
	// ensure both streams are started
	left.Start()
	right.Start()
	oC.LeftClientCh = make(chan *image.Gray, 10)  // Buffer a few frames
	oC.RightClientCh = make(chan *image.Gray, 10) // Buffer a few frames

	// Register this client
	left.Register <- oC.LeftClientCh
	right.Register <- oC.RightClientCh
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
	var errChan = make(chan error, 1)
	for {
		select {
		case <-ctx.Done():
			slog.Debug("stopping stream")
			return
		case err := <-errChan:
			slog.Error("Error reading image", "err", err)
			return
		case img := <-o.read(errChan):
			if img == nil {
				continue
			}
			outCh <- img
		}
	}
}

func (o *OutputCamera) read(_ chan error) <-chan *image.Gray {
	o.Params.Lock()
	defer o.Params.Unlock()
	mkdCh := make(chan *image.Gray, 1)
	leftImg := <-o.LeftClientCh
	rightImg := <-o.RightClientCh

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
	o.logger.Debug("Elapsed time", "elapsed", end.Sub(start))
	mkdCh <- got
	return mkdCh
}

// Close closes the output camera.
func (o *OutputCamera) Close() error {
	close(o.InputCh)
	o.Left.Stop()
	o.Right.Stop()
	return nil
}
