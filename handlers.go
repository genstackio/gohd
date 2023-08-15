package gohd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
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
	w.WriteHeader(302)
}

//goland:noinspection GoUnusedExportedFunction
func GetAndReturn(w http.ResponseWriter, req *http.Request, worker func(*http.Request) (interface{}, error)) {
	result, err := worker(req)
	if nil != err {
		JSONError(w, err, 500)
		return
	}
	body, err := json.Marshal(result)
	if nil != err {
		JSONError(w, err, 500)
		return
	}
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	w.WriteHeader(200)
	_, err = w.Write(body)
	if nil != err {
		log.Println(err.Error())
	}
}

//goland:noinspection GoUnusedExportedFunction
func CreateAndReturn(w http.ResponseWriter, req *http.Request, worker func(*http.Request) (interface{}, error)) {
	result, err := worker(req)
	if nil != err {
		JSONError(w, err, 500)
		return
	}
	body, err := json.Marshal(result)
	if nil != err {
		JSONError(w, err, 500)
		return
	}
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	w.WriteHeader(201)
	_, err = w.Write(body)
	if nil != err {
		log.Println(err.Error())
	}
}

//goland:noinspection GoUnusedExportedFunction
func PushEventToBackend(w http.ResponseWriter, req *http.Request, uriPrefix string, allowedHeaders map[string]bool, noBackendsStatusCode int, backends []string) {
	rawBody, err := io.ReadAll(req.Body)
	if err != nil {
		log.Printf("Error reading request body : %s", err.Error())
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(412)
		return
	}
	for i := 0; i < len(backends); i++ {
		body := io.NopCloser(bytes.NewBuffer(rawBody))
		uri := req.RequestURI
		if len(uriPrefix) > 0 {
			uri = uriPrefix + uri
		}
		req2, err := http.NewRequest("POST", fmt.Sprintf("%s%s", backends[i], uri), body)
		if err != nil {
			log.Println(err.Error())
			continue
		}
		for name, values := range req.Header {
			for _, value := range values {
				allowed, foundheader := allowedHeaders[name]
				if (true == foundheader) && (true == allowed) {
					req2.Header.Add(name, value)
				} else {
					// ignored request header
				}
			}
		}
		resp, err := http.DefaultClient.Do(req2)
		if err != nil {
			log.Println(err.Error())
			continue
		}
		if (resp.StatusCode >= 200) && (resp.StatusCode < 300) {
			defer resp.Body.Close()
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				log.Printf("Error reading response body : %s", err.Error())
				w.WriteHeader(500)
				break
			}
			for key, values := range resp.Header {
				for _, value := range values {
					w.Header().Add(key, value)
				}
			}
			w.Header().Set("X-Backend", fmt.Sprintf("Hit backend #%d", i))
			w.WriteHeader(resp.StatusCode)
			w.Write(body)
			return
		}
	}
	log.Println(fmt.Sprintf("No backend found for body [%s]", string(rawBody)))
	w.Header().Set("X-Error", "No backend found")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(noBackendsStatusCode)
}

//goland:noinspection GoUnusedExportedFunction
func JSONError(w http.ResponseWriter, err interface{}, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(err)
}
