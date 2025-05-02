package steroscopic

import (
	"context"
	"embed"
	"net/http"

	"github.com/conneroisu/steroscopic-hardware/cmd/steroscopic/components"
	"github.com/conneroisu/steroscopic-hardware/pkg/camera"
	"github.com/conneroisu/steroscopic-hardware/pkg/despair"
	"github.com/conneroisu/steroscopic-hardware/pkg/handlers"
)

//go:embed static/*
var static embed.FS

// AddRoutes adds the routes/handlers to the mux.
func AddRoutes(
	ctx context.Context,
	mux *http.ServeMux,
	params *despair.Parameters,
	leftStream, rightStream, outputStream *camera.StreamManager,
) error {
	// Add static files
	mux.Handle("/static/", http.StripPrefix(
		"/static/",
		http.FileServer(http.FS(static)),
	))

	mux.Handle("/", handlers.MorphableHandler(
		components.AppFn(handlers.LivePageTitle),
		components.Live(params.BlockSize, params.MaxDisparity),
	))
	mux.Handle("/manual", handlers.MorphableHandler(
		components.AppFn(handlers.ManualPageTitle),
		components.Manual(params.BlockSize, params.MaxDisparity),
	))

	mux.HandleFunc(
		"/ws",
		handlers.WSHandler,
	)
	mux.HandleFunc(
		"/api/parameters",
		handlers.Make(handlers.ParametersHandler(params)),
	)
	mux.HandleFunc(
		"/stream/left", // Left Camera
		handlers.Make(handlers.StreamHandlerFn(leftStream)),
	)
	mux.HandleFunc(
		"/stream/right", // Right Camera
		handlers.Make(handlers.StreamHandlerFn(rightStream)),
	)
	mux.HandleFunc(
		"/stream/out", // Depth Map
		handlers.Make(handlers.StreamHandlerFn(outputStream)),
	)
	mux.HandleFunc(
		"POST /manual-calc-depth-map",
		handlers.Make(handlers.ManualCalcDepthMapHandler()),
	)
	return nil
}
