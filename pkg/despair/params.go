package despair

import (
	"sync"
)

// Parameters is a struct that holds the parameters for the stereoscopic
// image processing.
type Parameters struct {
	mu           sync.Mutex
	BlockSize    int `json:"blockSize"`
	MaxDisparity int `json:"maxDisparity"`
}

// Lock locks the mutex.
func (p *Parameters) Lock() { p.mu.Lock() }

// Unlock unlocks the mutex.
func (p *Parameters) Unlock() { p.mu.Unlock() }
