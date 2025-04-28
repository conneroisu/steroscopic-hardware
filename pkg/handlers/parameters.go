package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/conneroisu/steroscopic-hardware/pkg/despair"
	"github.com/conneroisu/steroscopic-hardware/pkg/routing"
)

// ParametersHandler handles client requests to change the parameters of the
// desparity map generator.
func ParametersHandler(params *despair.Parameters) routing.APIFn {
	return func(_ http.ResponseWriter, r *http.Request) error {
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(params)
		if err != nil {
			return err
		}
		return nil
	}
}
