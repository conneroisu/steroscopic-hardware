package web

import (
	_ "embed"

	"github.com/a-h/templ"
)

//go:embed circle-question.svg
var circleQuestion string

//go:embed settings-gear.svg
var settingsGear string

//go:embed refresh-cw.svg
var refreshCw string

//go:embed green-up.svg
var greenUp string

//go:embed red-down.svg
var redDown string

//go:embed circle-x.svg
var circleX string

//go:embed file-icon.svg
var fileIcon string

// CircleQuestion is a template for the SVG circle-question icon.
var CircleQuestion = templ.Raw(circleQuestion)

// SettingsGear is a template for the SVG settings-geat icon.
var SettingsGear = templ.Raw(settingsGear)

// RefreshCw is a template for the SVG refresh-cw icon.
var RefreshCw = templ.Raw(refreshCw)

// GreenUp is a template for the SVG green-up icon.
var GreenUp = templ.Raw(greenUp)

// RedDown is a template for the SVG red-down icon.
var RedDown = templ.Raw(redDown)

// CircleX is a template for the SVG circle-x icon.
var CircleX = templ.Raw(circleX)

// FileIcon is a template for the SVG file-icon icon.
var FileIcon = templ.Raw(fileIcon)
