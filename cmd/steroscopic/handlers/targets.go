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
	return func(w http.ResponseWriter, _ *http.Request) error {
		var strBuilder strings.Builder
		ports, err := enumerator.GetDetailedPortsList()
		if err != nil {
			slog.Error(err.Error())
		}
		if len(ports) == 0 {
			return fmt.Errorf("no serial ports found")
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
