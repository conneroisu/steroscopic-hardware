package handlers

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/conneroisu/steroscopic-hardware/pkg/camera"
)

// UploadHandler handles the upload of a file.
func UploadHandler(typ camera.Type) APIFn {
	return func(_ http.ResponseWriter, r *http.Request) error {
		file, header, err := r.FormFile("file")
		if err != nil {
			return err
		}
		// save the file to tmp
		defer file.Close()
		body, err := io.ReadAll(file)
		if err != nil {
			return err
		}

		dir := os.TempDir()

		path := dir + "/" + header.Filename
		err = os.WriteFile(path, body, 0644)
		if err != nil {
			return err
		}

		switch typ {
		case camera.LeftCameraType:
			camera.SetLeftCamera(r.Context(), camera.NewStaticCamera(path, camera.LeftCh()))
		case camera.RightCameraType:
			camera.SetRightCamera(r.Context(), camera.NewStaticCamera(path, camera.RightCh()))
		default:
			return fmt.Errorf("unsupported camera type: %v", typ)
		}

		return nil
	}
}
