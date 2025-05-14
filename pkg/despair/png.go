package despair

import (
	"image"
	"image/png"
	"os"
)

// LoadPNG loads a PNG image and converts it to grayscale with optimizations.
func LoadPNG(filename string) (*image.Gray, error) {
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

// MustLoadPNG loads a PNG image and converts it to grayscale with
// optimizations and panics if an error occurs.
func MustLoadPNG(filename string) *image.Gray {
	img, err := LoadPNG(filename)
	if err != nil {
		panic(err)
	}

	return img
}

// SavePNG saves a PNG image with optimizations to the given filename
// and returns an error if one occurs.
func SavePNG(filename string, img image.Image) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// Use best compression for speed
	encoder := png.Encoder{CompressionLevel: png.BestSpeed}

	return encoder.Encode(file, img)
}

// MustSavePNG saves a PNG image with optimizations to the given filename
// and panics if an error occurs.
func MustSavePNG(filename string, img image.Image) {
	err := SavePNG(filename, img)
	if err != nil {
		panic(err)
	}
}
