package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/a-h/templ"
	"github.com/conneroisu/steroscopic-hardware/cmd/steroscopic/components"
	"github.com/conneroisu/steroscopic-hardware/pkg/camera"
	"github.com/conneroisu/steroscopic-hardware/pkg/web"
)

// PreviewDelimiterHandler handles requests to preview delimiters in different formats
func PreviewDelimiterHandler(w http.ResponseWriter, r *http.Request) {
	// Parse form data
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Failed to parse form data", http.StatusBadRequest)
		return
	}

	// Get the input values
	startDelimiter := r.FormValue("startdelimiter")
	endDelimiter := r.FormValue("enddelimiter")
	mode := r.FormValue("mode")

	// If no mode is specified, default to hex
	if mode == "" {
		mode = "hex"
	}

	// Parse the delimiters based on the selected mode
	startBytes := parseDelimiter(startDelimiter, mode)
	endBytes := parseDelimiter(endDelimiter, mode)

	// Create the preview HTML
	startPreviewHTML := web.FormatBytesForPreview(startBytes)
	endPreviewHTML := web.FormatBytesForPreview(endBytes)

	// Get camera type from the form data (if not provided, use a default)
	cameraType := strings.TrimSuffix(r.FormValue("id"), "-startDelimiter")
	cameraType = strings.TrimSuffix(cameraType, "-endDelimiter")
	if cameraType == "" {
		cameraType = "camera"
	}

	templ.Handler(components.DelimPreviewContainer(camera.Type(cameraType), startPreviewHTML, endPreviewHTML)).ServeHTTP(w, r)
}

// parseDelimiter parses a delimiter string based on the specified mode
func parseDelimiter(input string, mode string) []byte {
	var bytes []byte

	switch mode {
	case "hex":
		// Split by spaces and convert each hex value to a byte
		hexValues := strings.Fields(input)
		for _, hex := range hexValues {
			if hex = strings.TrimSpace(hex); hex != "" {
				// Remove "0x" prefix if present
				hex = strings.TrimPrefix(hex, "0x")
				value, err := strconv.ParseUint(hex, 16, 8)
				if err == nil && value <= 255 {
					bytes = append(bytes, byte(value))
				}
			}
		}
	case "decimal":
		// Split by spaces and convert each decimal value to a byte
		decValues := strings.Fields(input)
		for _, dec := range decValues {
			if dec = strings.TrimSpace(dec); dec != "" {
				value, err := strconv.ParseUint(dec, 10, 8)
				if err == nil && value <= 255 {
					bytes = append(bytes, byte(value))
				}
			}
		}
	case "text":
		// Convert each character to its character code
		bytes = []byte(input)
	}

	return bytes
}
