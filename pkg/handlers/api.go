package handlers

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/a-h/templ"
)

// APIFn is a function that handles an API request.
type APIFn func(w http.ResponseWriter, r *http.Request) error

// Make returns a function that can be used as an http.HandlerFunc.
func Make(fn APIFn) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := fn(w, r)
		if err != nil {
			slog.Error(fmt.Sprintf("API error: %s", err.Error()))
			http.Error(
				w,
				fmt.Sprintf(
					`{"success": false, "error": "%s"}`,
					err.Error(),
				),
				http.StatusInternalServerError,
			)
		}
	}
}

// MorphableHandler returns a handler that checks for the presence of the
// hx-trigger header and serves either the full or morphed view.
func MorphableHandler(
	wrapper func(templ.Component) templ.Component,
	morph templ.Component,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var header = r.Header.Get("HX-Request")
		if header == "" {
			templ.Handler(wrapper(morph)).ServeHTTPStreamed(w, r)
		} else {
			templ.Handler(morph).ServeHTTPStreamed(w, r)
		}
	}
}
