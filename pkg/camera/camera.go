// Package camera provides interfaces and implementations for cameras.
package camera

import (
	"image"
)

// Camer is the interface for a camera.
type Camer interface {
	Stream(chan *image.Gray)
	Close() error
}
