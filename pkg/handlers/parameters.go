package handlers

import (
	"encoding/json"
	"net/http"
)

// Parameters is a struct that holds the parameters for the stereoscopic
// image processing.
type Parameters struct {
	BlockSize    int `json:"blockSize"`
	MaxDisparity int `json:"maxDisparity"`
}

// ParametersHandler handles client requests to change the parameters of the
// desparity map generator.
func ParametersHandler(params *Parameters) APIFn {
	return func(_ http.ResponseWriter, r *http.Request) error {
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(params)
		if err != nil {
			return err
		}
		return nil
	}
}
