package gohd

import (
	baseErrors "errors"
	"net/http"
)

type workerFn func(*http.Request) (string, int, error, string)
type errorUrlFactoryFn func(code int, err error, lang string) (string, error)

func process(request *http.Request, worker workerFn, errorUrlFactory errorUrlFactoryFn) (string, int, error, error, string) {
	url, ttl, err, lang := worker(request)
	if err != nil {
		url, err2 := errorUrlFactory(10104, err, lang)
		if err2 != nil {
			return "", 0, err2, err2, lang
		}
		return url, 0, nil, err, lang
	}
	if len(url) == 0 {
		err0 := baseErrors.New("empty payment url")
		url, err = errorUrlFactory(10106, err0, lang)
		if err != nil {
			return "", 0, err, err, lang
		}
		return url, 0, nil, err0, lang
	}
	return url, ttl, nil, nil, lang

}

//goland:noinspection GoUnusedExportedFunction
func CreateAndRedirectNoError(w http.ResponseWriter, req *http.Request, worker workerFn, errorUrlFactory errorUrlFactoryFn) {
	if req.URL.Query().Has("noredirect") {
		CreateAndReturn(w, req, func(request *http.Request) (interface{}, error) {
			url, ttl, err, originalErr, lang := process(request, worker, errorUrlFactory)
			return struct {
				Url           string `json:"url,omitempty"`
				Ttl           int    `json:"ttl,omitempty"`
				Error         error  `json:"error,omitempty"`
				OriginalError error  `json:"originalError,omitempty"`
				Lang          string `json:"lang,omitempty"`
			}{Url: url, Ttl: ttl, Error: err, OriginalError: originalErr, Lang: lang}, err
		})
	}
	CreateAndRedirect(w, req, func(request *http.Request) (string, int, error) {
		url, ttl, err, _, _ := process(request, worker, errorUrlFactory)
		if err != nil {
			return "", 0, err
		}
		return url, ttl, nil
	})
}
