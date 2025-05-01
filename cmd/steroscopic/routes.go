package steroscopic

import (
	"context"
	"embed"
	"image"
	"net/http"

	"github.com/conneroisu/steroscopic-hardware/cmd/steroscopic/components"
	"github.com/conneroisu/steroscopic-hardware/pkg/camera"
	"github.com/conneroisu/steroscopic-hardware/pkg/handlers"
)

//go:embed static/*
var static embed.FS

// AddRoutes adds the routes/handlers to the mux.
func AddRoutes(
	ctx context.Context,
	mux *http.ServeMux,
	leftCamera camera.Camer,
	rightCamera camera.Camer,
) error {
	var (
		params     handlers.Parameters
		leftImgCh  = make(chan *image.Gray)
		rightImgCh = make(chan *image.Gray)
	)
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
		handlers.Make(handlers.ParametersHandler(&params)),
	)
	mux.HandleFunc(
		"/wsl", // Left Camera WebSocket
		handlers.Make(handlers.StreamHandlerFn(ctx, leftCamera, leftImgCh)),
	)
	mux.HandleFunc(
		"/wsr", // Right Camera WebSocket
		handlers.Make(handlers.StreamHandlerFn(ctx, rightCamera, rightImgCh)),
	)
	mux.HandleFunc(
		"/wso", // Depth Map WebSocket
		handlers.Make(handlers.GetMapHandler(leftImgCh, rightImgCh)),
	)
	mux.HandleFunc(
		"POST /manual-calc-depth-map",
		handlers.Make(handlers.ManualCalcDepthMapHandler()),
	)
	return nil
}
