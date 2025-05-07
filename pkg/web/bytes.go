package web

import (
	"fmt"
	"strings"
)

// BytesToHex converts a byte slice to a hex string.
func BytesToHex(bytes []byte) string {
	var hexString strings.Builder
	for _, b := range bytes {
		hexString.WriteString(fmt.Sprintf("%02x ", b))
	}
	return hexString.String()
}

// FormatBytesForPreview converts a byte slice to an HTML preview.
func FormatBytesForPreview(bytes []byte) string {
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

		sb.WriteString(fmt.Sprintf(`0x%02x (%s)`, b, visual))
	}

	return sb.String()
}
