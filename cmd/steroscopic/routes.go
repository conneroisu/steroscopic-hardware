package steroscopic

import (
	"embed"
	"fmt"
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
	template.New("index").Funcs(template.FuncMap{
		"dict": func(values ...any) (map[string]any, error) {
			if len(values)%2 != 0 {
				return nil, fmt.Errorf("invalid dict call")
			}
			dict := make(map[string]any, len(values)/2)
			for i := 0; i < len(values); i += 2 {
				key, ok := values[i].(string)
				if !ok {
					return nil, fmt.Errorf("dict keys must be strings")
				}
				dict[key] = values[i+1]
			}
			return dict, nil
		},
	}).Parse(htmlTmpl),
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

	// Index page route
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}

		data := map[string]any{
			"CurrentPage": "index",
		}

		err := tmpl.ExecuteTemplate(w, "index-with-nav", data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	// Manual page route
	mux.HandleFunc("/manual", func(w http.ResponseWriter, _ *http.Request) {
		data := map[string]any{
			"CurrentPage": "manual-page",
		}

		err := tmpl.ExecuteTemplate(w, "manual-page-with-nav", data)
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
	mux.HandleFunc(
		"/manual-calc-depth-map",
		handlers.Make(handlers.ManualCalcDepthMapHandler(cameraSystem)),
	)
	return nil
}
