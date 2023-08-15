package gohd

import (
	"encoding/json"
	"log"
	"net/http"
)

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
	w.WriteHeader(302)
}
func CreateAndReturn(w http.ResponseWriter, req *http.Request, worker func(*http.Request) (interface{}, error)) {
	result, err := worker(req)
	if nil != err {
		w.WriteHeader(500)
		_, err = w.Write([]byte(err.Error()))
		if nil != err {
			log.Println(err.Error())
		}
		return
	}
	body, err := json.Marshal(result)
	if nil != err {
		w.WriteHeader(500)
		_, err = w.Write([]byte(err.Error()))
		if nil != err {
			log.Println(err.Error())
		}
		return
	}
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	w.WriteHeader(201)
	_, err = w.Write(body)
	if nil != err {
		log.Println(err.Error())
	}
}
