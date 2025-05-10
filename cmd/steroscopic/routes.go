package steroscopic

import (
	"embed"
	"net/http"

	"github.com/conneroisu/steroscopic-hardware/cmd/steroscopic/components"
	"github.com/conneroisu/steroscopic-hardware/cmd/steroscopic/handlers"
	"github.com/conneroisu/steroscopic-hardware/pkg/camera"
	"github.com/conneroisu/steroscopic-hardware/pkg/despair"
	"github.com/conneroisu/steroscopic-hardware/pkg/logger"
	"github.com/conneroisu/steroscopic-hardware/pkg/web"
)

//go:embed static/*
var static embed.FS

// AddRoutes adds the routes/handlers to the mux.
func AddRoutes(
	mux *http.ServeMux,
	logger *logger.Logger,
	params *despair.Parameters,
	leftStream, rightStream, outputStream *camera.StreamManager,
) error {
	mux.Handle(
		"GET /",
		http.FileServer(http.FS(static)), // adds `/static/*` to path
	)
	mux.Handle("GET /{$}", handlers.MorphableHandler(
		components.AppFn(web.LivePageTitle),
		components.Live(params.BlockSize, params.MaxDisparity, leftStream, rightStream),
	))
	mux.Handle("GET /manual", handlers.MorphableHandler(
		components.AppFn(web.ManualPageTitle),
		components.Manual(params.BlockSize, params.MaxDisparity),
	))
	mux.HandleFunc(
		"POST /update-params",
		handlers.Make(handlers.ParametersHandler(logger, params)),
	)
	mux.HandleFunc(
		"GET /stream/left", // Left Camera
		handlers.Make(handlers.StreamHandlerFn(leftStream)),
	)
	mux.HandleFunc(
		"GET /stream/right", // Right Camera
		handlers.Make(handlers.StreamHandlerFn(rightStream)),
	)
	mux.HandleFunc(
		"GET /stream/out", // Depth Map
		handlers.Make(handlers.StreamHandlerFn(outputStream)),
	)
	mux.HandleFunc(
		"POST /manual-calc-depth-map",
		handlers.Make(handlers.ManualCalcDepthMapHandler(logger)),
	)
	mux.HandleFunc(
		"POST /left/configure",
		handlers.Make(handlers.ErrorHandler(
			handlers.ConfigureCamera(
				logger,
				params,
				leftStream,
				rightStream,
				outputStream,
				true,
			))),
	)
	mux.HandleFunc(
		"POST /right/configure",
		handlers.Make(handlers.ErrorHandler(
			handlers.ConfigureCamera(
				logger,
				params,
				leftStream,
				rightStream,
				outputStream,
				false,
			))),
	)
	mux.HandleFunc("GET /ports", handlers.Make(handlers.GetPorts(logger)))
	return nil
}
