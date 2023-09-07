package gohd

import "github.com/genstackio/goerror/errors"

type CreateErrorUrlFn = func(code int, err error, country string) (string, error)
type ResponseCompatible interface {
	GetStatusCode() int
	GetContentType() string
}
type DebugNoRedirectResponse struct {
	Url             string                    `json:"url,omitempty"`
	Ttl             int                       `json:"ttl,omitempty"`
	Error           *errors.JsonErrorResponse `json:"error,omitempty"`
	UrlFactoryError *errors.JsonErrorResponse `json:"urlFactoryError,omitempty"`
	Lang            string                    `json:"lang,omitempty"`
}

func (r DebugNoRedirectResponse) GetStatusCode() int {
	return r.Error.StatusCode
}

func (r DebugNoRedirectResponse) GetContentType() string {
	return "" // empty string will keep response as json
}
