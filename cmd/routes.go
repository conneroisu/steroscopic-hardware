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
// This function registers endpoints for camera control, streaming, and UI components.
func AddRoutes(
	ctx context.Context,
	mux *http.ServeMux,
	logger *logger.Logger,
	cancel context.CancelFunc,
) error {
	// Health check endpoint
	mux.HandleFunc("GET /checkhealth", func(_ http.ResponseWriter, _ *http.Request) {})

	// Static file server
	mux.Handle("GET /", http.FileServer(http.FS(static)))

	// Exit endpoint
	mux.HandleFunc("GET /exit", func(w http.ResponseWriter, _ *http.Request) {
		_, err := w.Write(logger.Bytes())
		if err != nil {
			log.Fatal("failed to write log", "err", err)
		}
		cancel()
	})

	// Main UI page
	mux.Handle("GET /{$}", handlers.MorphableHandler(
		components.AppFn(web.LivePageTitle),
		components.Live(),
	))

	// Parameter update endpoint
	mux.HandleFunc(
		"POST /update-params",
		handlers.Make(handlers.ParametersHandler()),
	)

	// Camera stream endpoints
	mux.HandleFunc(
		"GET /stream/left",
		handlers.Make(handlers.HandleLeftStream),
	)
	mux.HandleFunc(
		"GET /stream/right",
		handlers.Make(handlers.HandleRightStream),
	)
	mux.HandleFunc(
		"GET /stream/out",
		handlers.Make(handlers.HandleOutputStream),
	)

	// Left camera configuration and upload endpoints
	mux.HandleFunc(
		"POST /left/configure",
		handlers.Make(
			handlers.ErrorHandler(
				handlers.ConfigureMiddleware(
					handlers.ConfigureCamera(
						ctx,
						camera.LeftCameraType,
					)))),
	)
	mux.HandleFunc(
		"POST /left/upload",
		handlers.Make(
			handlers.ErrorHandler(
				handlers.UploadHandler(camera.LeftCameraType))),
	)

	// Right camera configuration and upload endpoints
	mux.HandleFunc(
		"POST /right/configure",
		handlers.Make(
			handlers.ErrorHandler(
				handlers.ConfigureMiddleware(
					handlers.ConfigureCamera(
						ctx,
						camera.RightCameraType,
					)))),
	)
	mux.HandleFunc(
		"POST /right/upload",
		handlers.Make(
			handlers.ErrorHandler(
				handlers.UploadHandler(camera.RightCameraType))),
	)

	// Available ports endpoint
	mux.HandleFunc("GET /ports", handlers.Make(handlers.GetPorts(logger)))

	return nil
}
