package gohd

import (
	"github.com/genstackio/goerror"
	"net/http"
)

//goland:noinspection GoUnusedExportedFunction
func FetchThenGenerateAndStream[T interface{}](w http.ResponseWriter, req *http.Request, fetcher func(req *http.Request) (T, error), worker func(data T, w http.ResponseWriter, req *http.Request) error) {
	data, err := fetcher(req)
	if nil != err {
		goerror.WriteError(w, err)
		return
	}
	err = worker(data, w, req)
	if nil != err {
		goerror.WriteError(w, err)
		return
	}
}
