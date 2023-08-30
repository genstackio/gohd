package gohd

import (
	baseErrors "errors"
	"net/http"
)

//goland:noinspection GoUnusedExportedFunction
func CreateAndRedirectNoError(w http.ResponseWriter, req *http.Request, worker func(*http.Request) (string, int, error, string), errorUrlFactory func(code int, err error, lang string) string) {
	CreateAndRedirect(w, req, func(request *http.Request) (string, int, error) {
		url, ttl, err, lang := worker(request)
		if err != nil {
			url := errorUrlFactory(10104, err, lang)
			return url, 0, err
		}
		if len(url) == 0 {
			url = errorUrlFactory(10106, baseErrors.New("empty payment url"), lang)
			return url, 0, err
		}
		return url, ttl, nil
	})
}
