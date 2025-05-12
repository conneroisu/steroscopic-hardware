package handlers

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/conneroisu/steroscopic-hardware/pkg/camera"
)

// ConfigureCamera handles client requests to configure all camera parameters
// at once.
func ConfigureCamera(
	ctx context.Context,
	typ camera.Type,
) APIFn {
	var logger = slog.Default().WithGroup("configure-camera")
	return func(_ http.ResponseWriter, r *http.Request) error {
		var (
			compression    int
			baudRate       int
			portStr        string
			baudStr        string
			compressionStr string
			config         camera.Config
		)

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
		config.Port = portStr

		// CONFIGURE baud rate if provided
		if baudStr == "" {
			return errors.New("baud rate not provided")
		}
		baudRate, err = strconv.Atoi(baudStr)
		if err != nil {
			return fmt.Errorf("invalid baud value: %w", err)
		}
		config.BaudRate = baudRate

		// CONFIGURE compression if provided
		if compressionStr == "" {
			return errors.New("compression not provided")
		}
		compression, err = strconv.Atoi(compressionStr)
		if err != nil {
			return fmt.Errorf("invalid compression value: %w", err)
		}
		config.Compression = compression

		// Log Configuration
		logger.Info(
			"setting",
			"stream",
			string(typ),
			"port",
			config.Port,
			"baud",
			config.BaudRate,
			"compression",
			config.Compression,
		)

		switch typ {
		case camera.LeftCameraType:
			cam, err := camera.NewSerialCamera(typ, portStr, baudRate, compression)
			if err != nil {
				return err
			}
			camera.SetLeftCamera(ctx, cam)
		case camera.RightCameraType:
			cam, err := camera.NewSerialCamera(typ, portStr, baudRate, compression)
			if err != nil {
				return err
			}
			camera.SetRightCamera(ctx, cam)
		default:
			return fmt.Errorf("unsupported camera type: %v", typ)
		}
		return nil
	}
}
