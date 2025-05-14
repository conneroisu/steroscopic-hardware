package despair

import (
	"image"
	"image/color"
	"math"
	"runtime"
	"sync"
)

// InputChunk represents a portion of the image to process.
type InputChunk struct {
	Left, Right *image.Gray
	Region      image.Rectangle
}

// OutputChunk represents the processed disparity data for a region.
type OutputChunk struct {
	DisparityData []uint8
	Region        image.Rectangle
}

// SetupConcurrentSAD sets up a concurrent SAD processing pipeline.
//
// It returns an input channel to feed image chunks into and an
// output channel to receive results from.
//
// If the input channel is closed, the processing pipeline will stop.
func SetupConcurrentSAD(
	numWorkers int, // Allow configurable worker count
) (chan<- InputChunk, <-chan OutputChunk) {
	if numWorkers <= 0 {
		numWorkers = runtime.NumCPU() * 4
	}

	inputChan := make(chan InputChunk, numWorkers*2)
	outputChan := make(chan OutputChunk, numWorkers*2)

	var wg sync.WaitGroup

	for range numWorkers { // i := 0; i < numWorkers; i++
		wg.Add(1)
		go func() {
			defer wg.Done()

			// Process chunks until the input channel is closed
			for chunk := range inputChan {
				data := make([]uint8,
					chunk.Region.Dx()*chunk.Region.Dy(),
				)
				defaultParamsMu.Lock()
				params := *defaultParams.Load()
				defaultParamsMu.Unlock()
				// Process each row in the region
				for y := range chunk.Region.Dy() { // y := 0; y < height; y++
					globalY := chunk.Region.Min.Y + y
					for x := range chunk.Region.Dx() { // x := 0; x < width; x++
						// Calculate disparity for this pixel
						minSAD := math.MaxInt32
						var bestDisparity int

						for d := 0; d <= params.MaxDisparity; d++ {
							// Skip if we would go beyond the left edge
							if (chunk.Region.Min.X+x)-d <
								chunk.Left.Rect.Min.X {
								continue
							}

							sad := SumAbsoluteDifferences(
								chunk.Left,
								chunk.Right,
								chunk.Region.Min.X+x,
								globalY,
								chunk.Region.Min.X+x-d,
								globalY,
								params.BlockSize,
							)

							if sad < minSAD {
								minSAD = sad
								bestDisparity = d

								// Early termination for perfect matches
								if sad == 0 {
									break
								}
							}
						}

						// Store the disparity value
						data[y*chunk.Region.Dx()+x] = uint8(
							(bestDisparity * 255) / params.MaxDisparity,
						)
					}
				}

				// Send the processed chunk to the output channel
				outputChan <- OutputChunk{
					DisparityData: data,
					Region:        chunk.Region,
				}
			}
		}()
	}

	// Goroutine to close output channel when all workers are done
	go func() {
		wg.Wait()
		close(outputChan)
	}()

	return inputChan, outputChan
}

// RunSad is a convenience function that sets up the pipeline,
// feeds the images, and assembles the disparity map.
//
// This is not used in the web UI, but is useful for testing.
func RunSad(
	left, right *image.Gray,
	blockSize, maxDisparity int,
) *image.Gray {
	SetDefaultParams(Parameters{
		BlockSize:    blockSize,
		MaxDisparity: maxDisparity,
	})
	// Determine number of workers and chunks
	numWorkers := runtime.NumCPU() * 4
	numChunks := numWorkers * 4

	// Set up the processing pipeline
	inputChan, outputChan := SetupConcurrentSAD(numWorkers)

	// Split the images into chunks
	var dims = left.Rect
	chunks := make([]image.Rectangle, 0, numChunks)

	chunkWidth := int(
		math.Sqrt(
			float64((dims.Dx() * dims.Dy()) / numChunks),
		),
	)
	horChunks := max(1, dims.Dx()/chunkWidth)
	verChunks := max(1, numChunks/horChunks)
	chunkWidth = dims.Dx() / horChunks
	chunkHeight := dims.Dy() / verChunks
	for startY := dims.Min.Y; startY < dims.Max.Y; startY += chunkHeight {
		endY := min(startY+chunkHeight, dims.Max.Y)
		for startX := dims.Min.X; startX < dims.Max.X; startX += chunkWidth {
			endX := min(startX+chunkWidth, dims.Max.X)
			chunks = append(chunks, image.Rect(startX, startY, endX, endY))
		}
	}

	// Start a goroutine to feed chunks into the pipeline
	go func() {
		for _, chunk := range chunks {
			inputChan <- InputChunk{
				Left:   left,
				Right:  right,
				Region: chunk,
			}
		}
		close(inputChan)
	}()

	// Assemble and return the disparity map
	return AssembleDisparityMap(outputChan, left.Rect, len(chunks))
}

// AssembleDisparityMap assembles the disparity map from output chunks.
func AssembleDisparityMap(
	outputChan <-chan OutputChunk,
	dimensions image.Rectangle,
	chunks int,
) *image.Gray {
	disparityMap := image.NewGray(dimensions)

	var i int
	for chunk := range outputChan {
		i++
		if i >= chunks {
			break
		}
		width := chunk.Region.Dx()
		for y := range chunk.Region.Dy() { // y := 0; y < chunk.Region.Dy(); y++
			globalY := chunk.Region.Min.Y + y
			for x := range width { // x := 0; x < width; x++
				globalX := chunk.Region.Min.X + x

				disparityValue := chunk.DisparityData[y*width+x]
				disparityMap.SetGray(
					globalX,
					globalY,
					color.Gray{Y: disparityValue},
				)
			}
		}
	}

	return disparityMap
}

// SumAbsoluteDifferences calculates SAD directly on image data.
func SumAbsoluteDifferences(
	left, right *image.Gray,
	leftX, leftY, rightX, rightY, blockSize int,
) int {
	halfSize := blockSize / 2

	// Optimize bounds checking by doing it once
	leftMinY := max(leftY-halfSize, 0)
	leftMaxY := min(leftY+halfSize+1, left.Rect.Max.Y)
	leftMinX := max(leftX-halfSize, 0)
	leftMaxX := min(leftX+halfSize+1, left.Rect.Max.X)

	rightMinY := max(rightY-halfSize, 0)
	rightMinX := max(rightX-halfSize, 0)

	// calculateSAD
	var (
		sad, lx int
	)
	for ly := leftMinY; ly < leftMaxY; ly++ {
		if rightMinY+(ly-leftMinY) >= right.Rect.Max.Y {
			break
		}
		leftRowStart := ly*left.Stride + leftMinX
		rightRowStart := (rightMinY+(ly-leftMinY))*right.Stride + rightMinX
		for lx = leftMinX; lx < leftMaxX; lx++ {
			if rightMinX+(lx-leftMinX) >= right.Rect.Max.X {
				break
			}
			diff := int(left.Pix[leftRowStart+lx-leftMinX]) -
				int(right.Pix[rightRowStart+(rightMinX+(lx-leftMinX))-rightMinX])
			if diff < 0 {
				diff = -diff
			}
			sad += diff
		}
	}

	return sad
}
