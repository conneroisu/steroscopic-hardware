package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/conneroisu/steroscopic-hardware/pkg/camera"
	"github.com/conneroisu/steroscopic-hardware/pkg/despair"
	"github.com/conneroisu/steroscopic-hardware/pkg/logger"
)

// ConfigureCamera handles client requests to configure all camera parameters at once.
func ConfigureCamera(
	logger *logger.Logger,
	params *despair.Parameters,
	leftStream, rightStream, outputStream *camera.StreamManager,
	isLeft bool,
) APIFn {
	return func(_ http.ResponseWriter, r *http.Request) error {
		var (
			compression     int
			baudRate        int
			portStr         string
			baudStr         string
			compressionStr  string
			configureStream *camera.StreamManager
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
		config := camera.DefaultCameraConfig()

		// Get all camera parameters
		portStr = r.FormValue("port")
		baudStr = r.FormValue("baudrate")
		compressionStr = r.FormValue("compression")

		// Configure port if provided
		if portStr == "" {
			return fmt.Errorf("port not provided")
		}
		config.Port = portStr
		logger.Info("configured camera port", "port", config.Port)

		// Configure baud rate if provided
		if baudStr == "" {
			return fmt.Errorf("baud rate not provided")
		}
		baudRate, err = strconv.Atoi(baudStr)
		if err != nil {
			return fmt.Errorf("invalid baud value: %w", err)
		}
		config.BaudRate = baudRate
		logger.Info("configured camera baud rate", "baud", config.BaudRate)

		// Configure compression if provided
		if compressionStr == "" {
			return fmt.Errorf("compression not provided")
		}
		compression, err = strconv.Atoi(compressionStr)
		if err != nil {
			return fmt.Errorf("invalid compression value: %w", err)
		}
		config.Compression = compression
		logger.Info("configured camera compression", "compression", config.Compression)

		// After configuration, attempt to connect
		err = configureStream.Configure(config)
		if err != nil {
			return fmt.Errorf("failed to configure camera: %w", err)
		}

		// After Connection, reconfigure the output camera
		outputStream.Stop()

		logger.Info("stopped output camera, creating new output camera")
		outputCamera := camera.NewOutputCamera(
			logger,
			params,
			leftStream,
			rightStream,
		)
		logger.Info("reconfigured output camera setting pointer")
		outputStream = camera.NewStreamManager(outputCamera, logger)

		return nil
	}
}
