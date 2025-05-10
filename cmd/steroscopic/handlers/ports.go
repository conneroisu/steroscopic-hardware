package handlers

import (
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/conneroisu/steroscopic-hardware/pkg/logger"
	"go.bug.st/serial/enumerator"
)

// GetPorts handles client requests to configure the camera.
func GetPorts(
	logger *logger.Logger,
) APIFn {
	return func(w http.ResponseWriter, r *http.Request) error {
		var (
			strBuilder strings.Builder
			tries      int
		)
		for {
			ports, err := enumerator.GetDetailedPortsList()
			if err != nil || len(ports) == 0 {
				if tries > 10 {
					return fmt.Errorf("no serial ports found")
				}
				slog.ErrorContext(r.Context(), "no serial ports found", "tries", tries)
				tries++
				continue
			}
			logger.Info("Found serial ports", "ports", ports, "len", len(ports))
			for _, port := range ports {
				strBuilder.WriteString(
					fmt.Sprintf("<option value=\"%s\">%s</option>\n",
						port.Name,
						port.Name,
					))
			}
			_, err = w.Write([]byte(strBuilder.String()))
			if err != nil {
				return err
			}
			return nil
		}
	}
}
