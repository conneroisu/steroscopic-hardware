package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/conneroisu/steroscopic-hardware/pkg/camera"
	"github.com/conneroisu/steroscopic-hardware/pkg/despair"
	"github.com/conneroisu/steroscopic-hardware/pkg/logger"
)

// ConfigureCamera handles client requests to configure all camera parameters
// at once.
func ConfigureCamera(
	logger *logger.Logger,
	params *despair.Parameters,
	leftStream, rightStream, outputStream *camera.Stream,
	isLeft bool,
) APIFn {
	return func(_ http.ResponseWriter, r *http.Request) error {
		var (
			compression     int
			baudRate        int
			portStr         string
			baudStr         string
			compressionStr  string
			configureStream *camera.Stream
			presetConfig    = camera.DefaultCameraConfig()
		)
		if isLeft {
			configureStream = leftStream
		} else {
			configureStream = rightStream
		}

		// Parse form data
		err := r.ParseForm()
		if err != nil {
			return fmt.Errorf("failed to parse form data: %w", err)
		}

		// Get all camera parameters
		portStr = r.FormValue("port")
		baudStr = r.FormValue("baudrate")
		compressionStr = r.FormValue("compression")

		// CONFIGURE port if provided
		if portStr == "" {
			return errors.New("port not provided")
		}
		presetConfig.Port = portStr

		// CONFIGURE baud rate if provided
		if baudStr == "" {
			return errors.New("baud rate not provided")
		}
		baudRate, err = strconv.Atoi(baudStr)
		if err != nil {
			return fmt.Errorf("invalid baud value: %w", err)
		}
		presetConfig.BaudRate = baudRate

		// CONFIGURE compression if provided
		if compressionStr == "" {
			return errors.New("compression not provided")
		}
		compression, err = strconv.Atoi(compressionStr)
		if err != nil {
			return fmt.Errorf("invalid compression value: %w", err)
		}
		presetConfig.Compression = compression

		// Log Configuration
		logger.Info(
			"setting",
			"stream",
			func() string { // inlined
				if isLeft {
					return "left"
				}
				return "right"
			}(),
			"port",
			presetConfig.Port,
			"baud",
			presetConfig.BaudRate,
			"compression",
			presetConfig.Compression,
		)
		// After configuration, attempt to connect
		err = configureStream.Configure(presetConfig)
		if err != nil {
			return fmt.Errorf("failed to configure stream: %w", err)
		}

		outputStream = camera.NewStreamManager(
			nil,
			logger,
			camera.WithReplace(
				outputStream,
				params,
				leftStream,
				rightStream,
			),
		)

		return nil
	}
}
