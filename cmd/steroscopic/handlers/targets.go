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

// Target is a struct representing a dom target.
type Target struct {
	ID  string `json:"id"`
	Sel string `json:"sel"`
}

var (
	// TargetLogContainer is a target for the log container.
	// It is used to insert log entries into the DOM.
	TargetLogContainer = Target{
		ID:  "log-container",
		Sel: "#log-container",
	}

	// TargetStatusContent is a target for the status content.
	// It is used to update the status content in the DOM.
	TargetStatusContent = Target{
		ID:  "status-content",
		Sel: "#status-content",
	}
)

var (
	// ManualPageTitle is the title of the manual page.
	ManualPageTitle = "Manual Depth Map Generator"
	// LivePageTitle is the title of the live page.
	LivePageTitle = "Live Camera System"
)
