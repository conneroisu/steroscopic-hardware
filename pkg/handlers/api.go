package handlers

import (
	"fmt"
	"log/slog"
	"net/http"
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
