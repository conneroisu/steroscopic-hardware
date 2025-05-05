package handlers

import (
	"net/http"

	"github.com/conneroisu/steroscopic-hardware/pkg/logger"
)

// ManualCalcDepthMapHandler is a handler for the manual depth map calculation endpoint.
func ManualCalcDepthMapHandler(
	logger *logger.Logger,
) APIFn {
	return func(w http.ResponseWriter, r *http.Request) error {
		logger.Info("manual depth map calculation requested")
		return nil
	}
}
