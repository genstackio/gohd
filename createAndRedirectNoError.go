package gohd

import (
	baseErrors "errors"
	"net/http"
)

//goland:noinspection GoUnusedExportedFunction
func CreateAndRedirectNoError(w http.ResponseWriter, req *http.Request, worker func(*http.Request) (string, int, error), errorUrlFactory func(code int, err error) string) {
	CreateAndRedirect(w, req, func(request *http.Request) (string, int, error) {
		url, ttl, err := worker(request)
		if err != nil {
			url := errorUrlFactory(10104, err)
			return url, 0, err
		}
		if len(url) == 0 {
			url = errorUrlFactory(10106, baseErrors.New("empty payment url"))
			return url, 0, err
		}
		return url, ttl, nil
	})
}
