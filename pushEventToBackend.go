package gohd

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
)

//goland:noinspection GoUnusedExportedFunction
func PushEventToBackend(w http.ResponseWriter, req *http.Request, uriPrefix string, allowedHeaders map[string]bool, noBackendsStatusCode int, backends []string) {
	rawBody, err := io.ReadAll(req.Body)
	if err != nil {
		log.Printf("Error reading request body : %s", err.Error())
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusPreconditionFailed)
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
				w.WriteHeader(http.StatusInternalServerError)
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
