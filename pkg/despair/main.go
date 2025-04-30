package despair

import (
	"image"
	"image/color"
	"math"
	"runtime"
	"sync"
)

// InputChunk represents a portion of the image to process
type InputChunk struct {
	Left, Right *image.Gray
	Region      image.Rectangle
}

// OutputChunk represents the processed disparity data for a region
type OutputChunk struct {
	DisparityData []uint8
	Region        image.Rectangle
}

// SetupConcurrentSAD sets up a concurrent SAD processing pipeline
// It returns an input channel to feed image chunks into and an
// output channel to receive results from.
//
// If the input channel is closed, the processing pipeline will stop.
func SetupConcurrentSAD(
	blockSize, maxDisparity int,
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
				left := chunk.Left
				right := chunk.Right
				region := chunk.Region

				// Create output data for this region
				width := region.Dx()
				height := region.Dy()
				disparityData := make([]uint8, width*height)

				// Process each row in the region
				for y := range height { // y := 0; y < height; y++
					globalY := region.Min.Y + y
					for x := range width { // x := 0; x < width; x++
						globalX := region.Min.X + x

						// Calculate disparity for this pixel
						minSAD := math.MaxInt32
						bestDisparity := 0

						for d := 0; d <= maxDisparity; d++ {
							// Skip if we would go beyond the left edge
							if globalX-d < left.Rect.Min.X {
								continue
							}

							sad := sumAbsoluteDifferences(
								left,
								right,
								globalX,
								globalY,
								globalX-d,
								globalY,
								blockSize,
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
						disparityData[y*width+x] = uint8((bestDisparity * 255) / maxDisparity)
					}
				}

				// Send the processed chunk to the output channel
				outputChan <- OutputChunk{
					DisparityData: disparityData,
					Region:        region,
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
// feeds the images, and assembles the disparity map
func RunSad(
	left, right *image.Gray,
	blockSize, maxDisparity int,
) *image.Gray {
	// Determine number of workers and chunks
	numWorkers := runtime.NumCPU() * 4
	numChunks := numWorkers * 4 // Create more chunks than workers for better load balancing

	// Set up the processing pipeline
	inputChan, outputChan := SetupConcurrentSAD(blockSize, maxDisparity, numWorkers)

	// Split the images into chunks
	chunks := splitImage(left.Rect, numChunks)

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

// RunSingleSad computes the disparity map with optimizations, using the concurrent infrastructure
func RunSingleSad(
	left, right *image.Gray,
	blockSize, maxDisparity int,
) *image.Gray {
	// Determine number of workers
	numWorkers := runtime.NumCPU() * 4

	// Set up the processing pipeline
	inputChan, outputChan := SetupConcurrentSAD(blockSize, maxDisparity, numWorkers)

	// Distribute work in row-based chunks for better cache utilization
	chunkSize := max(1, left.Rect.Dy()/(numWorkers*4))

	// Start a goroutine to feed chunks into the pipeline
	go func() {
		for y := left.Rect.Min.Y; y < left.Rect.Max.Y; y += chunkSize {
			endY := min(y+chunkSize, left.Rect.Max.Y)

			// Create a rectangular region for this chunk of rows
			region := image.Rect(left.Rect.Min.X, y, left.Rect.Max.X, endY)

			inputChan <- InputChunk{
				Left:   left,
				Right:  right,
				Region: region,
			}
		}
		close(inputChan)
	}()
	numChunks := (left.Rect.Dy() + chunkSize - 1) / chunkSize
	// Assemble and return the disparity map
	return AssembleDisparityMap(outputChan, left.Rect, numChunks)
}

// sumAbsoluteDifferences calculates SAD directly on image data
func sumAbsoluteDifferences(
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

	return calculateSAD(
		left,
		right,
		leftMinX,
		leftMaxX,
		leftMinY,
		leftMaxY,
		rightMinX,
		rightMinY,
	)
}

// calculateSAD performs the actual SAD calculation
func calculateSAD(
	left, right *image.Gray,
	leftMinX, leftMaxX, leftMinY, leftMaxY, rightMinX, rightMinY int,
) int {
	var (
		sad            int
		lx, ly, rx, ry int
	)
	for ly = leftMinY; ly < leftMaxY; ly++ {
		ry = rightMinY + (ly - leftMinY)
		if ry >= right.Rect.Max.Y {
			break
		}

		leftRowStart := ly*left.Stride + leftMinX
		rightRowStart := ry*right.Stride + rightMinX

		for lx = leftMinX; lx < leftMaxX; lx++ {
			rx = rightMinX + (lx - leftMinX)
			if rx >= right.Rect.Max.X {
				break
			}

			leftVal := left.Pix[leftRowStart+lx-leftMinX]
			rightVal := right.Pix[rightRowStart+rx-rightMinX]

			// Use abs without floating point
			diff := int(leftVal) - int(rightVal)
			if diff < 0 {
				// positive diff
				diff = -diff
			}
			sad += diff
		}
	}

	return sad
}

// AssembleDisparityMap assembles the disparity map from output chunks
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
		// Copy the disparity data to the appropriate location in the output image
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

// splitImage divides an image into rectangular chunks for parallel processing
func splitImage(
	dimensions image.Rectangle,
	numChunks int,
) []image.Rectangle {
	chunks := make([]image.Rectangle, 0, numChunks)

	// Try to make roughly square chunks
	totalPixels := dimensions.Dx() * dimensions.Dy()
	pixelsPerChunk := totalPixels / numChunks

	// Approximate width and height that gives square chunks
	chunkWidth := int(math.Sqrt(float64(pixelsPerChunk)))

	// Adjust to make sure we don't have too many chunks
	horChunks := max(1, dimensions.Dx()/chunkWidth)
	verChunks := max(1, numChunks/horChunks)

	chunkWidth = dimensions.Dx() / horChunks
	chunkHeight := dimensions.Dy() / verChunks

	// Create the chunks
	for startY := dimensions.Min.Y; startY < dimensions.Max.Y; startY += chunkHeight {
		endY := min(startY+chunkHeight, dimensions.Max.Y)

		for startX := dimensions.Min.X; startX < dimensions.Max.X; startX += chunkWidth {
			endX := min(startX+chunkWidth, dimensions.Max.X)

			chunks = append(chunks, image.Rect(startX, startY, endX, endY))
		}
	}

	return chunks
}
