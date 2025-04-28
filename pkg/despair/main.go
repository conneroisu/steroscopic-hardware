package despair

import (
	"fmt"
	"image"
	"image/color"
	"math"
	"runtime"
	"sync"
)

// processingChunk represents a chunk of work for parallel processing
type processingChunk struct {
	startY, endY int
}

// RunSad computes the disparity map with optimizations
func RunSad(
	left, right *image.Gray,
	blockSize, maxDisparity int,
) *image.Gray {
	disparityMap := image.NewGray(left.Rect)

	// Create worker pool
	numWorkers := runtime.NumCPU() * 4
	chunksChan := make(chan processingChunk, numWorkers*2)
	var wg sync.WaitGroup

	// Each worker processes chunks of rows
	for range numWorkers { // i := 0; i < numWorkers; i++
		wg.Add(1)
		go func() {
			for chunk := range chunksChan {
				var bounds = left.Rect
				for y := chunk.startY; y < chunk.endY; y++ {
					processRow(
						y,
						left,
						right,
						bounds,
						disparityMap,
						blockSize,
						maxDisparity,
					)
				}
			}
			wg.Done()
		}()
	}

	// Distribute work in larger chunks for better cache utilization
	chunkSize := max(1, left.Rect.Dy()/(numWorkers*4))
	for y := left.Rect.Min.Y; y < left.Rect.Max.Y; y += chunkSize {
		endY := min(y+chunkSize, left.Rect.Max.Y)
		chunksChan <- processingChunk{startY: y, endY: endY}
	}

	close(chunksChan)
	wg.Wait()

	return disparityMap
}

// RunSadPaths runs the optimized SAD algorithm on the given images
func RunSadPaths(
	left, right string,
	blockSize, maxDisparity int,
) error {
	// Load left and right images
	leftImg, err := LoadPNG(left)
	if err != nil {
		return err
	}

	rightImg, err := LoadPNG(right)
	if err != nil {
		return err
	}

	disparityMap := RunSad(
		leftImg,
		rightImg,
		blockSize,
		maxDisparity,
	)

	// Save disparity map
	err = SavePNG("disparity_map.png", disparityMap)
	if err != nil {
		return fmt.Errorf("error saving disparity map: %v", err)
	}

	fmt.Println("Disparity map saved successfully!")
	return nil
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

// processRow processes a single row for disparity calculation
func processRow(
	y int,
	left, right *image.Gray,
	bounds image.Rectangle,
	disparityMap *image.Gray,
	blockSize, maxDisparity int,
) {
	var x, disparity int
	for x = bounds.Min.X; x < bounds.Max.X; x++ {
		minSAD := math.MaxInt32

		for d := 0; d <= maxDisparity; d++ {
			// Skip if we would go beyond the left edge
			if x-d < bounds.Min.X {
				continue
			}

			sad := sumAbsoluteDifferences(
				left,
				right,
				x,
				y,
				x-d,
				y,
				blockSize,
			)

			if sad < minSAD {
				minSAD = sad
				disparity = d

				// Early termination for perfect matches
				if sad == 0 {
					break
				}
			}
		}
		disparityMap.SetGray(
			x,
			y,
			color.Gray{
				Y: uint8((disparity * 255) / maxDisparity),
			},
		)
	}
}
