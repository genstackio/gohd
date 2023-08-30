package gohd

import (
	baseErrors "errors"
	"github.com/genstackio/goerror"
	"github.com/genstackio/goerror/errors"
	"net/http"
	"strconv"
)

//goland:noinspection GoUnusedExportedFunction
func CreateAndRedirect(w http.ResponseWriter, req *http.Request, worker func(*http.Request) (string, int, error)) {
	url, ttl, err := worker(req)
	if err != nil {
		goerror.WriteError(w, errors.DocumentCreationError{
			Type: "redirect",
			Err:  baseErrors.New("unable to create redirect (" + err.Error() + ")"),
		})
		return
	}
	if 0 >= len(url) {
		goerror.WriteError(w, errors.UnknownDocumentError{
			Type: "redirect",
			Err:  baseErrors.New("unknown redirect (empty url)"),
		})
		return
	}
	w.Header().Set("Location", url)
	if ttl > 0 {
		w.Header().Set("cache-control", strconv.Itoa(ttl))
	}
	w.WriteHeader(http.StatusFound)
}
