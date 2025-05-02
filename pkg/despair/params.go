package despair

// Parameters is a struct that holds the parameters for the stereoscopic
// image processing.
type Parameters struct {
	BlockSize    int `json:"blockSize"`
	MaxDisparity int `json:"maxDisparity"`
}
