package steroscopic

import (
	"embed"
	"net/http"

	"github.com/conneroisu/steroscopic-hardware/cmd/steroscopic/components"
	"github.com/conneroisu/steroscopic-hardware/pkg/camera"
	"github.com/conneroisu/steroscopic-hardware/pkg/despair"
	"github.com/conneroisu/steroscopic-hardware/pkg/handlers"
	"github.com/conneroisu/steroscopic-hardware/pkg/logger"
)

//go:embed static/*
var static embed.FS

// AddRoutes adds the routes/handlers to the mux.
func AddRoutes(
	mux *http.ServeMux,
	logger *logger.Logger,
	params *despair.Parameters,
	leftStream, rightStream, outputStream *camera.StreamManager,
	leftCam, rightCam camera.Camer,
) error {
	mux.Handle(
		"GET /",
		http.FileServer(http.FS(static)), // adds `/static/*` to path
	)

	mux.Handle("GET /{$}", handlers.MorphableHandler(
		components.AppFn(handlers.LivePageTitle),
		components.Live(params.BlockSize, params.MaxDisparity),
	))
	mux.Handle("GET /manual", handlers.MorphableHandler(
		components.AppFn(handlers.ManualPageTitle),
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
		"GET /logs",
		handlers.Make(handlers.LogHandler(logger)))

	mux.HandleFunc("GET /configure/left", handlers.Make(
		handlers.ConfigureCameraHandler(logger, leftCam),
	))
	mux.HandleFunc("GET /configure/right", handlers.Make(
		handlers.ConfigureCameraHandler(logger, rightCam),
	))
	mux.HandleFunc("GET /ports", handlers.Make(handlers.GetPorts(logger)))
	return nil
}
