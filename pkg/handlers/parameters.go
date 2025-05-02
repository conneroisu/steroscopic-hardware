package handlers

import (
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/conneroisu/steroscopic-hardware/pkg/despair"
)

// ParametersHandler handles client requests to change the parameters of the
// desparity map generator.
func ParametersHandler(params *despair.Parameters) APIFn {
	return func(_ http.ResponseWriter, r *http.Request) error {
		// Parse form data
		if err := r.ParseForm(); err != nil {
			return fmt.Errorf("failed to parse form data: %w", err)
		}

		// For application/x-www-form-urlencoded or multipart/form-data
		blockSizeStr := r.FormValue("blockSize")
		maxDisparityStr := r.FormValue("maxDisparity")

		// Convert string values to integers
		blockSize, err := strconv.Atoi(blockSizeStr)
		if err != nil {
			return fmt.Errorf("invalid block size value: %w", err)
		}

		maxDisparity, err := strconv.Atoi(maxDisparityStr)
		if err != nil {
			return fmt.Errorf("invalid max disparity value: %w", err)
		}
		params.BlockSize = blockSize
		params.MaxDisparity = maxDisparity
		slog.Info(
			"received parameters:", "blocksize", params.BlockSize, "maxdisparity", params.MaxDisparity)
		return nil
	}
}
