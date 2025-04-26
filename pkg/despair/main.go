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

// processingChunk represents a chunk of work for parallel processing
type processingChunk struct {
	startY, endY int
}

// Precomputed lookup table for RGB to grayscale conversion
var rgbToGrayLUT [256]uint16

func init() {
	// Precompute RGB to grayscale conversion factors
	for i := range 256 { // i := 0; i < 256; i++
		rgbToGrayLUT[i] = uint16(i) * 255
	}
}

// sumAbsoluteDifferencesOptimized calculates SAD directly on image data
func sumAbsoluteDifferencesOptimized(left, right *image.Gray, leftX, leftY, rightX, rightY, blockSize int) int {
	var sad int
	halfSize := blockSize / 2

	// Direct access to the underlying pixel data
	leftStride := left.Stride
	rightStride := right.Stride
	leftPix := left.Pix
	rightPix := right.Pix

	// Optimize bounds checking by doing it once
	leftMinY := max(leftY-halfSize, 0)
	leftMaxY := min(leftY+halfSize+1, left.Rect.Max.Y)
	leftMinX := max(leftX-halfSize, 0)
	leftMaxX := min(leftX+halfSize+1, left.Rect.Max.X)

	rightMinY := max(rightY-halfSize, 0)
	rightMinX := max(rightX-halfSize, 0)

	// Process only the valid window
	for ly := leftMinY; ly < leftMaxY; ly++ {
		ry := rightMinY + (ly - leftMinY)
		if ry >= right.Rect.Max.Y {
			break
		}

		leftRowStart := ly*leftStride + leftMinX
		rightRowStart := ry*rightStride + rightMinX

		for lx := leftMinX; lx < leftMaxX; lx++ {
			rx := rightMinX + (lx - leftMinX)
			if rx >= right.Rect.Max.X {
				break
			}

			leftVal := leftPix[leftRowStart+lx-leftMinX]
			rightVal := rightPix[rightRowStart+rx-rightMinX]

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

// calculateDisparityMapOptimized computes the disparity map with optimizations
func calculateDisparityMapOptimized(left, right *image.Gray, blockSize, maxDisparity int) *image.Gray {
	bounds := left.Bounds()
	disparityMap := image.NewGray(bounds)

	// Create worker pool
	numWorkers := runtime.NumCPU()
	chunksChan := make(chan processingChunk, numWorkers*2)
	var wg sync.WaitGroup

	// Each worker processes chunks of rows
	for range numWorkers { // i := 0; i < numWorkers; i++
		wg.Add(1)
		go func() {
			defer wg.Done()

			// Reuse arrays to reduce allocations
			for chunk := range chunksChan {
				for y := chunk.startY; y < chunk.endY; y++ {
					for x := bounds.Min.X; x < bounds.Max.X; x++ {
						minSAD := math.MaxInt32
						bestDisparity := 0

						// Search for best match within disparity range
						for d := 0; d <= maxDisparity; d++ {
							// Skip if we would go beyond the left edge
							if x-d < bounds.Min.X {
								continue
							}

							sad := sumAbsoluteDifferencesOptimized(left, right, x, y, x-d, y, blockSize)

							if sad < minSAD {
								minSAD = sad
								bestDisparity = d

								// Early termination for perfect matches
								if sad == 0 {
									break
								}
							}
						}

						// Normalize disparity value for better visualization
						disparityValue := uint8((bestDisparity * 255) / maxDisparity)
						disparityMap.SetGray(x, y, color.Gray{Y: disparityValue})
					}
				}
			}
		}()
	}

	// Distribute work in larger chunks for better cache utilization
	chunkSize := max(1, bounds.Dy()/(numWorkers*4))
	for y := bounds.Min.Y; y < bounds.Max.Y; y += chunkSize {
		endY := min(y+chunkSize, bounds.Max.Y)
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
		copy(grayPix, img.Pix)
	case *image.RGBA:
		// Direct access to RGBA pixel data
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			rowStart := (y - bounds.Min.Y) * stride
			for x := bounds.Min.X; x < bounds.Max.X; x++ {
				i := img.PixOffset(x, y)
				r := img.Pix[i]
				g := img.Pix[i+1]
				b := img.Pix[i+2]

				// Use integer arithmetic
				grayPix[rowStart+x-bounds.Min.X] = uint8((19595*uint32(r) + 38470*uint32(g) + 7471*uint32(b) + 1<<15) >> 24)
			}
		}
	default:
		// Fallback for other image types
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			rowStart := (y - bounds.Min.Y) * stride
			for x := bounds.Min.X; x < bounds.Max.X; x++ {
				r, g, b, _ := img.At(x, y).RGBA()
				grayPix[rowStart+x-bounds.Min.X] = uint8((19595*r + 38470*g + 7471*b + 1<<15) >> 24)
			}
		}
	}

	return grayImg, nil
}

// savePNG saves an image as PNG
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
func RunSad(left, right string) error {
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
	blockSize := 15
	maxDisparity := 64

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
