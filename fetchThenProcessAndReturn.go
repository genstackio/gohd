package gohd

import (
	"encoding/json"
	"github.com/genstackio/goerror"
	"github.com/genstackio/goerror/errors"
	"net/http"
)

//goland:noinspection GoUnusedExportedFunction
func FetchThenProcessAndReturn[T interface{}, U interface{}](w http.ResponseWriter, req *http.Request, fetcher func(req *http.Request) (*T, error), worker func(data *T, req *http.Request) (*U, error)) {
	data, err := fetcher(req)
	if nil != err {
		goerror.WriteError(w, err)
		return
	}
	result, err2 := worker(data, req)
	if nil != err2 {
		goerror.WriteError(w, err2)
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
	}
}
