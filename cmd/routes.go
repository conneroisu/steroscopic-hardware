package cmd

import (
	"context"
	"embed"
	"log"
	"net/http"

	"github.com/conneroisu/steroscopic-hardware/cmd/components"
	"github.com/conneroisu/steroscopic-hardware/cmd/handlers"
	"github.com/conneroisu/steroscopic-hardware/pkg/camera"
	"github.com/conneroisu/steroscopic-hardware/pkg/logger"
	"github.com/conneroisu/steroscopic-hardware/pkg/web"
)

// static contains embedded static web assets served by the HTTP server.
// This includes CSS, JavaScript, and icon files needed by the web UI.
//
//go:embed static/*
var static embed.FS

// AddRoutes configures all HTTP routes and handlers for the application.
//
// This function registers the following endpoints:
//   - GET /: Static file server for web assets
//   - GET /exit: Endpoint to request server shutdown (returns logs)
//   - GET /{$}: Main UI page with live camera streams and controls
//   - POST /update-params: Update stereoscopic algorithm parameters
//   - GET /stream/{left,right,out}: Live MJPEG streams for both cameras and depth map
//   - POST /{left,right}/configure: Configure camera devices and parameters
//   - GET /ports: List available serial ports for camera connections
//
// Parameters:
//   - mux: The HTTP server mux to register routes on
//   - logger: Application logger for recording events
//   - params: Stereoscopic algorithm parameters
//   - leftStream: Stream manager for the left camera
//   - rightStream: Stream manager for the right camera
//   - outputStream: Stream manager for the generated depth map
//   - cancel: CancelFunc to trigger application shutdown
//
// Returns any error encountered during route configuration.
func AddRoutes(
	ctx context.Context,
	mux *http.ServeMux,
	logger *logger.Logger,
	cancel context.CancelFunc,
) error {
	mux.HandleFunc("GET /checkhealth", func(_ http.ResponseWriter, _ *http.Request) {})
	mux.Handle("GET /", http.FileServer(http.FS(static)))
	mux.HandleFunc("GET /exit", func(w http.ResponseWriter, _ *http.Request) {
		_, err := w.Write(logger.Bytes())
		if err != nil {
			log.Fatal("failed to write log", "err", err)
		}
		cancel()
	})
	mux.Handle("GET /{$}", handlers.MorphableHandler(
		components.AppFn(web.LivePageTitle),
		components.Live(),
	))
	mux.HandleFunc(
		"POST /update-params",
		handlers.Make(handlers.ParametersHandler()),
	)
	mux.HandleFunc(
		"GET /stream/left", // Left Camera
		handlers.Make(handlers.HandleLeftStream),
	)
	mux.HandleFunc(
		"GET /stream/right", // Right Camera
		handlers.Make(handlers.HandleRightStream),
	)
	mux.HandleFunc(
		"GET /stream/out", // Depth Map
		handlers.Make(handlers.HandleOutputStream),
	)
	mux.HandleFunc(
		"POST /left/configure",
		handlers.Make(
			handlers.ErrorHandler(
				handlers.ConfigureMiddleware(
					handlers.ConfigureCamera(
						ctx,
						camera.LeftCameraType,
					)))))
	mux.HandleFunc(
		"POST /right/configure",
		handlers.Make(
			handlers.ErrorHandler(
				handlers.ConfigureMiddleware(
					handlers.ConfigureCamera(
						ctx,
						camera.RightCameraType,
					)))))
	mux.HandleFunc("GET /ports", handlers.Make(handlers.GetPorts(logger)))
	return nil
}
