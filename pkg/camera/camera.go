// Package camera provides interfaces and implementations for cameras.
package camera

import (
	"context"
	"image"
)

// Camer is the interface for a camera.
type Camer interface {
	Stream(context.Context, chan *image.Gray)
	Close() error
	ID() string
}
