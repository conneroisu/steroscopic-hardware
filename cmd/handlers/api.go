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
		defer func() {
			if err := recover(); err != nil {
				// redirect to error page
				http.Redirect(w, r, "/error", http.StatusFound)
			}
		}()
		err := fn(w, r)
		if err != nil {
			slog.Error(
				"api error",
				"err",
				err,
				"url",
				r.URL,
				"method",
				r.Method,
			)
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

// ErrorHandler returns a handler that returns an error response.
func ErrorHandler(
	fn APIFn,
) APIFn {
	return func(w http.ResponseWriter, r *http.Request) error {
		err := fn(w, r)
		if err == nil {
			_, Werr := w.Write([]byte(`<span class="text-sm text-green-500">Success!</span>`))
			if Werr != nil {
				return fmt.Errorf("failed to write success response: %w", Werr)
			}
			return nil
		}
		// Return error response
		_, Werr := w.Write([]byte(`<span class="text-sm text-red-500">Failure: ` + err.Error() + `</span>`))
		return Werr
	}
}
