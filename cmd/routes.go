package cmd

import (
	"embed"
	"html/template"
	"net/http"
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
	mux.HandleFunc("/{$}", func(w http.ResponseWriter, _ *http.Request) {
		if err := tmpl.ExecuteTemplate(w, "index", nil); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.FS(static))))
	return nil
}
