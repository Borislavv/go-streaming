package service

import "net/http"

type Request interface {
	Parameters(req *http.Request) map[string]string
	HasParameter(param string, req *http.Request) bool
	GetParameter(param string, req *http.Request) (string, error)
}
