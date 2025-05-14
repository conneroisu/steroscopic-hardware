package despair

import (
	"sync"
	"sync/atomic"
)

var (
	defaultParams   = atomic.Pointer[Parameters]{}
	defaultParamsMu sync.Mutex
)

func init() {
	SetDefaultParams(Parameters{
		BlockSize:    16,
		MaxDisparity: 64,
	})
}

// SetDefaultParams sets the default stereoscopic algorithm parameters.
func SetDefaultParams(params Parameters) {
	defaultParamsMu.Lock()
	defer defaultParamsMu.Unlock()
	defaultParams.Store(&params)
}

// DefaultParams returns the default stereoscopic algorithm parameters.
func DefaultParams() *Parameters {
	return defaultParams.Load()
}

// Parameters is a struct that holds the parameters for the stereoscopic
// image processing.
type Parameters struct {
	BlockSize    int `json:"blockSize"`
	MaxDisparity int `json:"maxDisparity"`
}
