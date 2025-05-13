package camera

import (
	"context"
	"image"
	"log/slog"
	"time"

	"github.com/conneroisu/steroscopic-hardware/pkg/despair"
)

// DefaultNumWorkers is the default number of worker goroutines for disparity calculations.
const DefaultNumWorkers = 32

// OutputCamera processes left and right camera images to generate a depth map.
type OutputCamera struct {
	BaseCamera
	inputCh  chan<- despair.InputChunk
	outputCh <-chan despair.OutputChunk
	logger   *slog.Logger
}

// NewOutputCamera creates a new output camera for disparity mapping.
func NewOutputCamera(ctx context.Context) *OutputCamera {
	oc := &OutputCamera{
		BaseCamera: NewBaseCamera(ctx, OutputCameraType),
		logger:     slog.Default().WithGroup("output-camera"),
	}

	// Initialize the disparity processing pipeline
	oc.inputCh, oc.outputCh = despair.SetupConcurrentSAD(DefaultNumWorkers)

	return oc
}

// Stream processes input images and generates depth maps.
func (oc *OutputCamera) Stream(ctx context.Context, outCh ImageChannel) {
	oc.logger.Info("starting output camera stream")
	defer oc.logger.Info("output camera stream stopped")

	// Backoff parameters
	initialBackoff := 10 * time.Millisecond
	maxBackoff := 1 * time.Second
	backoff := initialBackoff
	consecutiveFailures := 0

	for {
		select {
		case <-ctx.Done():
			oc.logger.Debug("context canceled, stopping stream")
			return
		case <-oc.Context().Done():
			oc.logger.Debug("camera context canceled, stopping stream")
			return
		default:
			if oc.IsPaused() {
				time.Sleep(100 * time.Millisecond)
				continue
			}

			// Generate depth map from input channels
			img, err := oc.processDepthMap(DefaultManager().GetOutputChannel(LeftCameraType), DefaultManager().GetOutputChannel(RightCameraType))
			if err != nil {
				oc.logger.Error("error processing depth map", "err", err)
				time.Sleep(100 * time.Millisecond)
				continue
			}

			if img == nil {
				time.Sleep(100 * time.Millisecond)
				continue
			}

			// Send processed image to output channel
			select {
			case outCh <- img:
				oc.logger.Debug("depth map sent to channel")
				// Reset backoff on success
				backoff = initialBackoff
				consecutiveFailures = 0
			case <-ctx.Done():
				return
			case <-oc.Context().Done():
				return
			default:
				// Channel is full, apply backoff
				consecutiveFailures++

				// Only log at certain thresholds to prevent log spam
				if consecutiveFailures == 1 || consecutiveFailures%10 == 0 {
					oc.logger.Debug("output channel full, applying backoff",
						"consecutiveFailures", consecutiveFailures,
						"currentBackoff", backoff)
				}

				// Apply backoff delay
				time.Sleep(backoff)

				// Exponential backoff with a maximum cap
				backoff = time.Duration(float64(backoff) * 1.5)
				if backoff > maxBackoff {
					backoff = maxBackoff
				}
			}
		}
	}
}

// processDepthMap generates a depth map from left and right camera images.
func (oc *OutputCamera) processDepthMap(leftCh, rightCh ImageChannel) (*image.Gray, error) {
	// Try to receive images from both channels
	var leftImg, rightImg *image.Gray

	// Create a timeout context for receiving images
	timeoutCtx, cancel := context.WithTimeout(oc.Context(), 1*time.Second)
	defer cancel()

	// Try to receive left image
	select {
	case img := <-leftCh:
		leftImg = img
	case <-timeoutCtx.Done():
		return nil, nil // No image available yet
	}

	// Try to receive right image
	select {
	case img := <-rightCh:
		rightImg = img
	case <-timeoutCtx.Done():
		return nil, nil // No image available yet
	}

	// Process images if both are available
	if leftImg != nil && rightImg != nil {
		startTime := time.Now()

		// Get current parameters
		params := despair.DefaultParams()

		// Divide image into chunks for parallel processing
		chunkSize := max(1, leftImg.Rect.Dy()/(DefaultNumWorkers*4))
		numChunks := (leftImg.Rect.Dy() + chunkSize - 1) / chunkSize

		// Send chunks to processing pipeline
		for y := leftImg.Rect.Min.Y; y < leftImg.Rect.Max.Y; y += chunkSize {
			oc.inputCh <- despair.InputChunk{
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

		// Assemble the resulting disparity map
		disparityMap := despair.AssembleDisparityMap(oc.outputCh, leftImg.Rect, numChunks)

		elapsedTime := time.Since(startTime)
		oc.logger.Info("depth map generated",
			"elapsed", elapsedTime,
			"blockSize", params.BlockSize,
			"maxDisparity", params.MaxDisparity)

		return disparityMap, nil
	}

	return nil, nil
}

// Close releases all resources.
func (oc *OutputCamera) Close() error {
	oc.logger.Info("closing output camera")
	oc.Cancel()

	// Close the input channel to stop workers
	close(oc.inputCh)

	return nil
}
