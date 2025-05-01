package handlers

import (
	"net/http"
)

// ManualCalcDepthMapHandler is a handler for the manual depth map calculation endpoint.
func ManualCalcDepthMapHandler() APIFn {
	return func(w http.ResponseWriter, r *http.Request) error {
		return nil
	}
}
