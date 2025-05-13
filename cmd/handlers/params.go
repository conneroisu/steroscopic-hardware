package handlers

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/conneroisu/steroscopic-hardware/pkg/despair"
)

// ParametersHandler handles client requests to update disparity algorithm parameters.
func ParametersHandler() APIFn {
	logger := slog.Default().WithGroup("params-handler")

	return func(_ http.ResponseWriter, r *http.Request) error {
		// Parse form data
		if err := r.ParseForm(); err != nil {
			return fmt.Errorf("failed to parse form data: %w", err)
		}

		// Get form values
		blockSizeStr := r.FormValue("blockSize")
		maxDisparityStr := r.FormValue("maxDisparity")

		// Validate block size
		if blockSizeStr == "" {
			return errors.New("block size not provided")
		}
		blockSize, err := strconv.Atoi(blockSizeStr)
		if err != nil {
			return fmt.Errorf("invalid block size value: %w", err)
		}

		// Block size must be odd and within range
		if blockSize < 3 || blockSize > 31 || blockSize%2 == 0 {
			return errors.New("block size must be odd and between 3 and 31")
		}

		// Validate max disparity
		if maxDisparityStr == "" {
			return errors.New("max disparity not provided")
		}
		maxDisparity, err := strconv.Atoi(maxDisparityStr)
		if err != nil {
			return fmt.Errorf("invalid max disparity value: %w", err)
		}

		// Max disparity must be within range and divisible by 16
		if maxDisparity < 16 || maxDisparity > 256 || maxDisparity%16 != 0 {
			return errors.New("max disparity must be between 16 and 256 and divisible by 16")
		}

		// Update parameters
		despair.SetDefaultParams(despair.Parameters{
			BlockSize:    blockSize,
			MaxDisparity: maxDisparity,
		})

		logger.Info(
			"parameters updated",
			"blockSize", blockSize,
			"maxDisparity", maxDisparity,
		)

		return nil
	}
}
