package handlers

import (
	"net/http"

	"github.com/a-h/templ"
)

// MorphableHandler returns a handler that checks for the presence of the
// hx-trigger header and serves either the full or morphed view.
func MorphableHandler(
	full templ.Component,
	morph templ.Component,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var header = r.Header.Get("HX-Request")
		if header == "" {
			templ.Handler(full).ServeHTTPStreamed(w, r)
		} else {
			templ.Handler(morph).ServeHTTPStreamed(w, r)
		}
	}
}
