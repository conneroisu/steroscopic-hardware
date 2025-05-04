// Package svg contains SVG templates for the web UI.
package svg

import (
	_ "embed"

	"github.com/a-h/templ"
)

//go:embed circle-question.svg
var circleQuestion string

//go:embed settings-gear.svg
var settingsGear string

// CircleQuestion is a template for the SVG circle-question icon.
var CircleQuestion = templ.Raw(circleQuestion)

// SettingsGear is a template for the SVG settings-geat icon.
var SettingsGear = templ.Raw(settingsGear)
