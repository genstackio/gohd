package gohd

import (
	"encoding/json"
	"github.com/genstackio/goerror"
	"github.com/genstackio/goerror/errors"
	"net/http"
)

//goland:noinspection GoUnusedExportedFunction
func GetAndReturn[T interface{}](w http.ResponseWriter, req *http.Request, worker func(*http.Request) (T, error)) {
	result, err := worker(req)
	if nil != err {
		goerror.WriteError(w, err)
		return
	}
	body, err := json.Marshal(result)
	if nil != err {
		goerror.WriteError(w, errors.MarshallError{Err: err})
		return
	}
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(body)
	if nil != err {
		goerror.WriteError(w, errors.WriteError{Err: err})
		return
	}
}
