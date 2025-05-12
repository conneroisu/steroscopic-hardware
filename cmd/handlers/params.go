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
func ParametersHandler() APIFn {
	logger := slog.Default().WithGroup("params-handler")
	return func(_ http.ResponseWriter, r *http.Request) error {
		var (
			blockSizeStr    string
			maxDisparityStr string
		)
		// Parse form data
		if err := r.ParseForm(); err != nil {
			return fmt.Errorf("failed to parse form data: %w", err)
		}

		// For application/x-www-form-urlencoded or multipart/form-data
		blockSizeStr = r.FormValue("blockSize")
		maxDisparityStr = r.FormValue("maxDisparity")

		// Convert string values to integers
		blockSize, err := strconv.Atoi(blockSizeStr)
		if err != nil {
			return fmt.Errorf(
				"invalid block size value: %w",
				err,
			)
		}

		maxDisparity, err := strconv.Atoi(maxDisparityStr)
		if err != nil {
			return fmt.Errorf(
				"invalid max disparity value: %w",
				err,
			)
		}
		despair.SetDefaultParams(despair.Parameters{
			BlockSize:    blockSize,
			MaxDisparity: maxDisparity,
		})
		logger.Info(
			"received parameters:",
			"blocksize",
			blockSize,
			"maxdisparity",
			maxDisparity,
		)
		return nil
	}
}
