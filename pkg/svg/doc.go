package svg

import (
	_ "embed"

	"github.com/a-h/templ"
)

//go:embed circle-question.svg
var circleQuestion string

// CircleQuestion is a template for the SVG circle-question icon.
var CircleQuestion = templ.Raw(circleQuestion)
