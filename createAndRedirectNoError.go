package gohd

import (
	baseErrors "errors"
	"net/http"
)

//goland:noinspection GoUnusedExportedFunction
func CreateAndRedirectNoError(w http.ResponseWriter, req *http.Request, worker func(*http.Request) (string, int, error, string), errorUrlFactory func(code int, err error, lang string) (string, error)) {
	CreateAndRedirect(w, req, func(request *http.Request) (string, int, error) {
		url, ttl, err, lang := worker(request)
		if err != nil {
			url, err2 := errorUrlFactory(10104, err, lang)
			if err2 != nil {
				return "", 0, err2
			}
			return url, 0, nil
		}
		if len(url) == 0 {
			url, err = errorUrlFactory(10106, baseErrors.New("empty payment url"), lang)
			if err != nil {
				return "", 0, err
			}
			return url, 0, nil
		}
		return url, ttl, nil
	})
}
