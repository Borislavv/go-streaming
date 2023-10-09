package request

import "net/http"

type Extractor interface {
	Parameters(req *http.Request) map[string]string
	HasParameter(param string, req *http.Request) bool
	GetParameter(param string, req *http.Request) (string, error)
}
