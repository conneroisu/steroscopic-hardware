package web

import (
	"fmt"
	"strings"
)

func BytesToHex(bytes []byte) string {
	var hexString strings.Builder
	for _, b := range bytes {
		hexString.WriteString(fmt.Sprintf("%02x ", b))
	}
	return hexString.String()
}
