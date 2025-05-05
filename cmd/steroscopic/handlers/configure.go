package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/conneroisu/steroscopic-hardware/pkg/camera"
	"github.com/conneroisu/steroscopic-hardware/pkg/despair"
	"github.com/conneroisu/steroscopic-hardware/pkg/logger"
)

// ParametersHandler handles client requests to change the parameters of the
// desparity map generator.
func ParametersHandler(logger *logger.Logger, params *despair.Parameters) APIFn {
	return func(_ http.ResponseWriter, r *http.Request) error {
		params.Lock()
		defer params.Unlock()
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
		logger.Info(
			"received parameters:", "blocksize", params.BlockSize, "maxdisparity", params.MaxDisparity)
		return nil
	}
}

// CameraConfig represents all configurable camera parameters
type CameraConfig struct {
	Port        string
	BaudRate    int
	Compression int
}

// DefaultCameraConfig returns default camera configuration
func DefaultCameraConfig() CameraConfig {
	return CameraConfig{
		Port:        "",
		BaudRate:    115200,
		Compression: 0,
	}
}

// ConfigureCamera handles all camera configuration in a single handler
func ConfigureCamera(
	logger *logger.Logger,
	stream *camera.StreamManager,
) APIFn {
	return func(w http.ResponseWriter, r *http.Request) error {
		// Parse form/query data
		if err := r.ParseForm(); err != nil {
			return fmt.Errorf("failed to parse form data: %w", err)
		}

		config := DefaultCameraConfig()

		// Get port value (check both form values and query parameters)
		if port := r.FormValue("port"); port != "" {
			config.Port = port
		} else if port := r.URL.Query().Get("port"); port != "" {
			config.Port = port
		}

		// Get baud rate value
		if baudStr := r.FormValue("baudrate"); baudStr != "" {
			baud, err := strconv.Atoi(baudStr)
			if err != nil {
				return fmt.Errorf("invalid baud rate value: %w", err)
			}
			config.BaudRate = baud
		} else if baudStr := r.URL.Query().Get("baudrate"); baudStr != "" {
			baud, err := strconv.Atoi(baudStr)
			if err != nil {
				return fmt.Errorf("invalid baud rate value: %w", err)
			}
			config.BaudRate = baud
		}

		// Get compression value
		if compStr := r.FormValue("compression"); compStr != "" {
			comp, err := strconv.Atoi(compStr)
			if err != nil {
				return fmt.Errorf("invalid compression value: %w", err)
			}
			config.Compression = comp
		} else if compStr := r.URL.Query().Get("compression"); compStr != "" {
			comp, err := strconv.Atoi(compStr)
			if err != nil {
				return fmt.Errorf("invalid compression value: %w", err)
			}
			config.Compression = comp
		}

		// Configure the camera with all parameters
		logger.Info("configuring camera",
			"port", config.Port,
			"baudrate", config.BaudRate,
			"compression", config.Compression)

		// First configure port if provided
		if config.Port != "" {
			if err := stream.Configure(config.Port); err != nil {
				return fmt.Errorf("failed to configure camera port: %w", err)
			}
		}

		// Then configure baud rate
		if err := stream.Configure(config.BaudRate); err != nil {
			return fmt.Errorf("failed to configure camera baud rate: %w", err)
		}

		// Finally configure compression
		if err := stream.Configure(config.Compression); err != nil {
			return fmt.Errorf("failed to configure camera compression: %w", err)
		}

		// Write a status response
		w.Header().Set("Content-Type", "text/html")
		connected := (config.Port != "")
		var statusHTML string
		if connected {
			statusHTML = `<span class="inline-block w-3 h-3 bg-green-500 rounded-full"></span>
<span class="text-sm">Connected</span>`
		} else {
			statusHTML = `<span class="inline-block w-3 h-3 bg-red-500 rounded-full"></span>
<span class="text-sm">Disconnected</span>`
		}
		_, err := w.Write([]byte(statusHTML))
		return err
	}
}
