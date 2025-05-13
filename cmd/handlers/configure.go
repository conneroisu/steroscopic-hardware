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

// CtxKey is a type alias for context keys used to store camera configuration.
type CtxKey string

const (
	ctxKeyConfig CtxKey = "config"
)

// ConfigureMiddleware parses camera configuration from form data.
func ConfigureMiddleware(apiFn APIFn) APIFn {
	return func(w http.ResponseWriter, r *http.Request) error {
		// Parse form data
		if err := r.ParseForm(); err != nil {
			return fmt.Errorf("failed to parse form data: %w", err)
		}

		// Get form values
		portStr := r.FormValue("port")
		baudStr := r.FormValue("baudrate")
		compressionStr := r.FormValue("compression")

		// Validate port
		if portStr == "" {
			return errors.New("port not provided")
		}

		// Validate and convert baud rate
		if baudStr == "" {
			return errors.New("baud rate not provided")
		}
		baudRate, err := strconv.Atoi(baudStr)
		if err != nil {
			return fmt.Errorf("invalid baud rate value: %w", err)
		}

		// Validate and convert compression
		if compressionStr == "" {
			return errors.New("compression not provided")
		}
		compression, err := strconv.Atoi(compressionStr)
		if err != nil {
			return fmt.Errorf("invalid compression value: %w", err)
		}

		// Create config
		config := camera.Config{
			Port:        portStr,
			BaudRate:    baudRate,
			Compression: compression,
		}

		// Add to request context
		ctx := context.WithValue(r.Context(), ctxKeyConfig, config)

		// Call the next handler
		return apiFn(w, r.WithContext(ctx))
	}
}

// ConfigureCamera handles client requests to configure camera parameters.
func ConfigureCamera(ctx context.Context, typ camera.Type) APIFn {
	logger := slog.Default().WithGroup("configure-camera")

	return func(w http.ResponseWriter, r *http.Request) error {
		// Get config from context
		config, ok := r.Context().Value(ctxKeyConfig).(camera.Config)
		if !ok {
			return errors.New("camera configuration not found in request context")
		}

		// Log configuration
		logger.Info(
			"configuring camera",
			"type", string(typ),
			"port", config.Port,
			"baud", config.BaudRate,
			"compression", config.Compression,
		)

		// Create and configure the camera - using the application context instead of request context
		cam, err := camera.NewSerialCamera(ctx, typ, config.Port, config.BaudRate, config.Compression)
		if err != nil {
			return fmt.Errorf("failed to create serial camera: %w", err)
		}

		// Set the camera in the manager
		err = camera.SetCamera(ctx, typ, cam)
		if err != nil {
			return fmt.Errorf("failed to set camera: %w", err)
		}

		// Return success HTML
		successHTML := `
		<span class="inline-block w-3 h-3 bg-green-500 rounded-full"></span>
		<span class="text-sm">Connected successfully!</span>
		`
		_, err = w.Write([]byte(successHTML))
		if err != nil {
			return fmt.Errorf("failed to write success HTML: %w", err)
		}

		return nil
	}
}
