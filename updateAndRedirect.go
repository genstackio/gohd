package gohd

import (
	"net/http"
)

//goland:noinspection GoUnusedExportedFunction
func UpdateAndRedirect(w http.ResponseWriter, req *http.Request, worker func(*http.Request) (string, error), createErrorUrl CreateErrorUrlFn) {
	url, err := worker(req)
	if nil != err {
		locale := req.URL.Query().Get("locale")
		if 0 == len(locale) {
			locale = "FR"
		}
		url, _ = createErrorUrl(10210, err, locale)
	}
	w.Header().Set("Location", url)
	w.WriteHeader(http.StatusFound)
}
