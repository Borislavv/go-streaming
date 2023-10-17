package extractor

import "net/http"

type RequestParams interface {
	Parameters(r *http.Request) map[string]string
	HasParameter(param string, req *http.Request) bool
	GetParameter(param string, req *http.Request) (string, error)
}
