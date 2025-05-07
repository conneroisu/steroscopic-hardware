package handlers

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"
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
	startPreviewHTML := formatBytesForPreview(startBytes)
	endPreviewHTML := formatBytesForPreview(endBytes)

	// Get camera type from the form data (if not provided, use a default)
	cameraType := strings.TrimSuffix(r.FormValue("id"), "-startDelimiter")
	cameraType = strings.TrimSuffix(cameraType, "-endDelimiter")
	if cameraType == "" {
		// Try to extract from the referer or other parts of the request
		// For now, use a default
		cameraType = "camera"
	}

	// Render the preview template
	tmpl := `
	<div class="flex flex-col">
		<span class="text-sm text-gray-300 mb-1">Start preview:</span>
		<div class="bg-gray-700 text-gray-200 rounded px-3 py-2 text-sm border border-gray-600 min-h-8 font-mono" id="{{ .CameraType }}-startPreview">
			{{ if .StartPreview }}{{ .StartPreview | unescapeHTML }}{{ else }}Preview will appear here{{ end }}
		</div>
	</div>
	<div class="flex flex-col mt-2">
		<span class="text-sm text-gray-300 mb-1">End preview:</span>
		<div class="bg-gray-700 text-gray-200 rounded px-3 py-2 text-sm border border-gray-600 min-h-8 font-mono" id="{{ .CameraType }}-endPreview">
			{{ if .EndPreview }}{{ .EndPreview | unescapeHTML }}{{ else }}Preview will appear here{{ end }}
		</div>
	</div>
	`

	// Create template with a function to unescape HTML
	t, err := template.New("preview").Funcs(template.FuncMap{
		"unescapeHTML": func(s string) template.HTML {
			return template.HTML(s)
		},
	}).Parse(tmpl)

	if err != nil {
		http.Error(w, "Failed to parse template", http.StatusInternalServerError)
		return
	}

	// Execute the template
	data := struct {
		CameraType   string
		StartPreview string
		EndPreview   string
	}{
		CameraType:   cameraType,
		StartPreview: startPreviewHTML,
		EndPreview:   endPreviewHTML,
	}

	w.Header().Set("Content-Type", "text/html")
	err = t.Execute(w, data)
	if err != nil {
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
	}
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

// formatBytesForPreview converts a byte slice to an HTML preview
func formatBytesForPreview(bytes []byte) string {
	if len(bytes) == 0 {
		return ""
	}

	var sb strings.Builder
	for _, b := range bytes {
		var visual string
		if b >= 32 && b <= 126 {
			// Printable ASCII
			visual = string(b)
			// Escape HTML special characters
			visual = strings.ReplaceAll(visual, "&", "&amp;")
			visual = strings.ReplaceAll(visual, "<", "&lt;")
			visual = strings.ReplaceAll(visual, ">", "&gt;")
			visual = strings.ReplaceAll(visual, "\"", "&quot;")
		} else {
			// Non-printable, show as dot
			visual = "·"
		}

		sb.WriteString(fmt.Sprintf(`<span title="Decimal: %d">0x%02x (%s)</span> `, b, b, visual))
	}

	return sb.String()
}
