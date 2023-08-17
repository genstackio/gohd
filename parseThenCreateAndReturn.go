package gohd

import (
	"encoding/json"
	"github.com/genstackio/goerror"
	"github.com/genstackio/goerror/errors"
	"net/http"
)

//goland:noinspection GoUnusedExportedFunction
func ParseThenCreateAndReturn(w http.ResponseWriter, req *http.Request, init func(*http.Request) (interface{}, error), worker func(interface{}, *http.Request) (interface{}, error)) {
	data, err := init(req)

	if err != nil {
		JSONError(w, err, http.StatusBadRequest)
		return
	}

	err = json.NewDecoder(req.Body).Decode(data)
	if err != nil {
		goerror.WriteError(w, errors.MalformedPayloadError{Err: err})
		return
	}

	CreateAndReturn(w, req, func(r *http.Request) (interface{}, error) {
		return worker(data, req)
	})
}
