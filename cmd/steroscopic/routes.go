package steroscopic

import (
	"embed"
	"html/template"
	"net/http"

	"github.com/conneroisu/steroscopic-hardware/pkg/handlers"
)

const (
	defLeftPort  = "8080"
	defRightPort = "8081"
)

//go:embed static/*
var static embed.FS

//go:embed html.tmpl
var htmlTmpl string

var tmpl = template.Must(
	template.New("index").Parse(htmlTmpl),
)

// AddRoutes adds the routes/handlers to the mux.
func AddRoutes(
	mux *http.ServeMux,
) error {
	var params handlers.Parameters
	cameraSystem := handlers.NewCameraSystem(
		defLeftPort,
		defRightPort,
		"static/images",
		&params,
	)
	mux.HandleFunc("/{$}", func(w http.ResponseWriter, _ *http.Request) {
		err := tmpl.ExecuteTemplate(w, "index", nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})
	mux.HandleFunc(
		"/ws",
		handlers.WSHandler,
	)
	mux.HandleFunc(
		"/api/parameters",
		handlers.Make(handlers.ParametersHandler(&params)),
	)
	mux.Handle(
		"/static/",
		http.StripPrefix(
			"/static/",
			http.FileServer(http.FS(static)),
		))
	mux.HandleFunc(
		"/api/capture",
		handlers.Make(handlers.CameraHandler(cameraSystem)),
	)
	mux.HandleFunc(
		"/wsl",
		handlers.Make(handlers.GetStreamHandler(cameraSystem, "left")),
	)
	mux.HandleFunc(
		"/wsr",
		handlers.Make(handlers.GetStreamHandler(cameraSystem, "right")),
	)
	mux.HandleFunc(
		"/wso",
		handlers.Make(handlers.GetMapHandler(cameraSystem)),
	)

	return nil
}
