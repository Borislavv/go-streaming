package request

import (
	"github.com/Borislavv/video-streaming/internal/domain/errs"
	"github.com/gorilla/mux"
	"net/http"
)

type Request struct {
}

func NewRequestService() *Request {
	return &Request{}
}

// Parameters - will return a map with param. names as keys to values
func (r *Request) Parameters(req *http.Request) map[string]string {
	return mux.Vars(req)
}

// HasParameter - checking the param. is existing in request
func (r *Request) HasParameter(param string, req *http.Request) bool {
	if _, ok := r.Parameters(req)[param]; ok {
		return true
	}
	return false
}

// GetParameter - checking if the param. is existing in request, it will be returned
func (r *Request) GetParameter(param string, req *http.Request) (string, error) {
	if v, ok := r.Parameters(req)[param]; ok {
		return v, nil
	}
	return "", errs.NewFieldCannotBeEmptyError(param)
}
