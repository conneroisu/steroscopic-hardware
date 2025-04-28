package despair

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
	"runtime"
	"sync"
)

// Parameters is a struct that holds the parameters for the stereoscopic
// image processing.
type Parameters struct {
	BlockSize    int `json:"blockSize"`
	MaxDisparity int `json:"maxDisparity"`
}

// processingChunk represents a chunk of work for parallel processing
type processingChunk struct {
	startY, endY int
}

// sumAbsoluteDifferencesOptimized calculates SAD directly on image data
func sumAbsoluteDifferencesOptimized(
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

	return calculateSAD(left, right, leftMinX, leftMaxX, leftMinY, leftMaxY, rightMinX, rightMinY)
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
				diff = -diff
			}
			sad += diff
		}
	}

	return sad
}

// processChunk processes a chunk of rows for disparity calculation
func processChunk(
	chunk processingChunk,
	left, right *image.Gray,
	bounds image.Rectangle,
	disparityMap *image.Gray,
	blockSize, maxDisparity int,
) {
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
		disparity = findBestDisparity(
			x,
			y,
			left,
			right,
			bounds,
			blockSize,
			maxDisparity,
		)
		disparityMap.SetGray(
			x,
			y,
			color.Gray{Y: uint8((disparity * 255) / maxDisparity)},
		)
	}
}

// findBestDisparity finds the best disparity value for a given point
func findBestDisparity(
	x, y int,
	left, right *image.Gray,
	bounds image.Rectangle,
	blockSize, maxDisparity int,
) int {
	minSAD := math.MaxInt32
	bestDisparity := 0

	for d := 0; d <= maxDisparity; d++ {
		// Skip if we would go beyond the left edge
		if x-d < bounds.Min.X {
			continue
		}

		sad := sumAbsoluteDifferencesOptimized(
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
			bestDisparity = d

			// Early termination for perfect matches
			if sad == 0 {
				break
			}
		}
	}

	return bestDisparity
}

// calculateDisparityMapOptimized computes the disparity map with optimizations
func calculateDisparityMapOptimized(left, right *image.Gray, blockSize, maxDisparity int) *image.Gray {
	disparityMap := image.NewGray(left.Rect)

	// Create worker pool
	numWorkers := runtime.NumCPU()
	chunksChan := make(chan processingChunk, numWorkers*2)
	var wg sync.WaitGroup

	// Each worker processes chunks of rows
	for range numWorkers { // i := 0; i < numWorkers; i++
		wg.Add(1)
		go func() {
			defer wg.Done()
			for chunk := range chunksChan {
				processChunk(
					chunk,
					left,
					right,
					left.Rect,
					disparityMap,
					blockSize,
					maxDisparity,
				)
			}
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

// loadPNGAsGrayOptimized loads a PNG image and converts it to grayscale with optimizations
func loadPNGAsGrayOptimized(filename string) (*image.Gray, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	img, err := png.Decode(file)
	if err != nil {
		return nil, err
	}

	bounds := img.Bounds()
	grayImg := image.NewGray(bounds)

	// Direct access to pixel data
	grayPix := grayImg.Pix
	stride := grayImg.Stride

	// Optimize by checking image type
	switch img := img.(type) {
	case *image.Gray:
		convertGrayToGray(img, grayPix)
	case *image.RGBA:
		convertRGBAToGray(img, grayPix, stride, bounds)
	default:
		convertGenericToGray(img, grayPix, stride, bounds)
	}

	return grayImg, nil
}

func savePNG(filename string, img image.Image) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// Use best compression for speed
	encoder := png.Encoder{CompressionLevel: png.BestSpeed}
	return encoder.Encode(file, img)
}

// RunSad runs the optimized SAD algorithm on the given images
func RunSad(left, right string, params *Parameters) error {
	// Load left and right images
	leftImg, err := loadPNGAsGrayOptimized(left)
	if err != nil {
		return fmt.Errorf("error loading left image: %v", err)
	}

	rightImg, err := loadPNGAsGrayOptimized(right)
	if err != nil {
		return fmt.Errorf("error loading right image: %v", err)
	}

	// Parameters for block matching
	blockSize := params.BlockSize
	maxDisparity := params.MaxDisparity

	// Calculate disparity map
	fmt.Println("Calculating disparity map...")
	fmt.Printf("Using %d CPU cores for parallel processing\n", runtime.NumCPU())
	disparityMap := calculateDisparityMapOptimized(leftImg, rightImg, blockSize, maxDisparity)

	// Save disparity map
	err = savePNG("disparity_map.png", disparityMap)
	if err != nil {
		return fmt.Errorf("error saving disparity map: %v", err)
	}

	fmt.Println("Disparity map saved successfully!")
	return nil
}
