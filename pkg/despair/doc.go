// Package despair provides a Go implementation of a stereoscopic depth mapping algorithm,
// designed for efficient generation of depth/disparity maps from stereo image pairs.
//
// # Core Functionality
//
// The package implements the Sum of Absolute Differences (SAD) algorithm, a common technique
// in stereoscopic vision that:
//
//  1. Takes left and right grayscale images from slightly different viewpoints
//  2. Compares blocks of pixels to find matching points between images
//  3. Calculates disparity (horizontal displacement) between matching points
//  4. Generates a grayscale disparity map where pixel brightness represents depth
//
// # Data Structures
//
//	InputChunk: Represents a portion of the image pair to process
//	OutputChunk: Contains processed disparity data for a specific region
//	Parameters: Configuration settings for the algorithm including:
//	`BlockSize`: Size of pixel blocks for comparison
//	`MaxDisparity`: Maximum pixel displacement to check
//
// # Processing Pipeline
//
//  1. `SetupConcurrentSAD`: Creates a pipeline with configurable worker count, returning input/output channels
//
//  2. `RunSad`: Convenience function that orchestrates the entire process:
//     - Divides images into chunks
//     - Distributes processing across workers
//     - Assembles final disparity map
//
//  3. `AssembleDisparityMap`: Combines processed chunks into a complete disparity map
//
//  4. `sumAbsoluteDifferences`: Low-level function that calculates block matching scores
//
// # Image Handling
//
// The package includes efficient image handling utilities:
//
//   - PNG Loading/Saving: Optimized functions for loading and saving grayscale PNG images
//   - Type-Specific Conversions: Specialized routines for different image formats (Gray, RGBA, generic)
//   - Error Handling: Both standard error-returning functions and "Must" variants that panic on failure
//
// # Performance Optimizations
//
//   - Concurrent Processing: Utilizes Go's concurrency with multiple worker go-routines
//   - Chunked Processing: Splits images into smaller regions for parallel processing
//   - Direct Pixel Access: Works with underlying pixel arrays rather than the higher-level interface
//   - Type-Specific Optimizations: Different code paths for different image types
//   - Early Termination: Breaks comparison loops when perfect matches are found
//   - Optimized Bounds Checking: Reduces redundant checks in inner loops
//   - Precomputed Lookup Tables: Uses LUTs for common conversions
//
// Example:
//
//	```go
//	// Load stereo image pair
//	left := despair.MustLoadPNG("left.png")
//	right := despair.MustLoadPNG("right.png")
//
//	// Generate disparity map with block size 9 and max disparity 64
//	disparityMap := despair.RunSad(left, right, 9, 64)
//
//	// Save the result
//	despair.MustSavePNG("depth_map.png", disparityMap)
//	```
package despair

//go:generate gomarkdoc -o README.md -e .
