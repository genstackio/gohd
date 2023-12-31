package gohd

import (
	"encoding/json"
	"github.com/genstackio/goerror"
	"github.com/genstackio/goerror/errors"
	"net/http"
)

//goland:noinspection GoUnusedExportedFunction
func DeleteAndReturn[T interface{}](w http.ResponseWriter, req *http.Request, worker func(*http.Request) (T, error)) {
	result, err := worker(req)
	if err != nil {
		goerror.WriteError(w, err)
		return
	}
	body, err := json.Marshal(result)
	if nil != err {
		goerror.WriteError(w, errors.MarshallError{Err: err})
		return
	}
	statusCode := http.StatusOK
	contentType := "application/json;charset=utf-8"
	if z, ok := any(result).(ResponseCompatible); ok {
		forcedStatusCode := z.GetStatusCode()
		if forcedStatusCode != 0 {
			statusCode = forcedStatusCode
		}
		forcedContentType := z.GetContentType()
		if "" != forcedContentType {
			contentType = forcedContentType
		}
	}
	w.Header().Set("Content-Type", contentType)
	w.WriteHeader(statusCode)
	_, err = w.Write(body)
	if nil != err {
		goerror.WriteError(w, errors.WriteError{Err: err})
	}
}
