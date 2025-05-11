// Package web contains SVG templates and dom targets for the web UI.
//
// SVG templates are used to render SVG icons and text in the web UI.
// Templates are embedded into the package using the go:embed directive.
//
// Targets are used to identify DOM elements in the web UI and are defined as
// constants with a unique ID and a CSS selector.
// Targets are used to update the DOM with new content or to insert new elements.
package web

//go:generate gomarkdoc -o README.md -e .
