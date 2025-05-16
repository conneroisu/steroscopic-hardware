package despair

import "image"

var rgbToGrayLUT [256]uint16

func init() {
	// Precompute RGB to grayscale conversion factors
	for i := range 256 {
		rgbToGrayLUT[i] = uint16(i) * 255
	}
}

// convertGrayToGray directly copies gray image data.
func convertGrayToGray(src *image.Gray, grayPix []uint8) {
	copy(grayPix, src.Pix)
}

// convertRGBAToGray converts RGBA image to grayscale.
func convertRGBAToGray(
	src *image.RGBA,
	grayPix []uint8,
	stride int,
	bounds image.Rectangle,
) {
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		rowStart := (y - bounds.Min.Y) * stride
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			i := src.PixOffset(x, y)
			r := src.Pix[i]
			g := src.Pix[i+1]
			b := src.Pix[i+2]

			// Use integer arithmetic
			grayPix[rowStart+x-bounds.Min.X] = uint8((19595*uint32(r) +
				38470*uint32(g) +
				7471*uint32(b) + 1<<15) >> 24)
		}
	}
}

// convertGenericToGray converts any image to grayscale.
func convertGenericToGray(
	src image.Image,
	grayPix []uint8,
	stride int,
	bounds image.Rectangle,
) {
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		rowStart := (y - bounds.Min.Y) * stride
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := src.At(x, y).RGBA()
			grayPix[rowStart+x-bounds.Min.X] = uint8((19595*r +
				38470*g +
				7471*b + 1<<15) >> 24)
		}
	}
}
