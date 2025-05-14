package handlers

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/conneroisu/steroscopic-hardware/pkg/logger"
	"go.bug.st/serial/enumerator"
)

// GetPorts handles client requests to configure the camera.
func GetPorts(
	logger *logger.Logger,
) APIFn {
	return func(w http.ResponseWriter, _ *http.Request) error {
		var (
			strBuilder strings.Builder
			tries      int
			port       *enumerator.PortDetails
			ports      []*enumerator.PortDetails
			err        error
		)
		for {
			ports, err = enumerator.GetDetailedPortsList()
			if err != nil || len(ports) == 0 {
				if tries > 10 {
					return errors.New("no serial ports found")
				}
				logger.Error(
					"no serial ports found",
					"tries",
					tries,
				)
				time.Sleep(time.Second)
				tries++

				continue
			}
			logger.Info(
				"Found serial ports",
				"# of ports",
				len(ports),
			)
			for _, port = range ports {
				strBuilder.WriteString("<option value=\"")
				strBuilder.WriteString(port.Name)
				strBuilder.WriteString("\">")
				strBuilder.WriteString(port.Name)
				strBuilder.WriteString("</option>\n")
			}
			_, err = w.Write([]byte(strBuilder.String()))
			if err != nil {
				return err
			}

			return nil
		}
	}
}
