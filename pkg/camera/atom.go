package camera

import (
	"image"
	"sync/atomic"
)

var (
	defaultLeftCh   = atomic.Pointer[chan *image.Gray]{}
	defaultRightCh  = atomic.Pointer[chan *image.Gray]{}
	defaultOutputCh = atomic.Pointer[chan *image.Gray]{}
)

func init() {
	defaultLeftCh.Store(new(chan *image.Gray))
	defaultRightCh.Store(new(chan *image.Gray))
	defaultOutputCh.Store(new(chan *image.Gray))
}

// Left returns the left channel.
func Left() chan *image.Gray {
	return *defaultLeftCh.Load()
}

// Right returns the right channel.
func Right() chan *image.Gray {
	return *defaultRightCh.Load()
}

// Output returns the output channel.
func Output() chan *image.Gray {
	return *defaultOutputCh.Load()
}
