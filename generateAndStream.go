package gohd

import (
	"github.com/genstackio/goerror"
	"net/http"
)

//goland:noinspection GoUnusedExportedFunction
func GenerateAndStream(w http.ResponseWriter, req *http.Request, worker func(http.ResponseWriter, *http.Request) error) {
	w.WriteHeader(http.StatusOK)
	err := worker(w, req)
	if nil != err {
		goerror.WriteError(w, err)
		return
	}
}
