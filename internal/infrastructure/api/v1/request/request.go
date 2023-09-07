package request

import (
	"github.com/Borislavv/video-streaming/internal/domain/errs"
	"github.com/gorilla/mux"
	"net/http"
)

type Extractor interface {
	Parameters(req *http.Request) map[string]string
	HasParameter(param string, req *http.Request) bool
	GetParameter(param string, req *http.Request) (string, error)
}

type ParametersExtractor struct {
}

func NewParametersExtractor() *ParametersExtractor {
	return &ParametersExtractor{}
}

// Parameters - will return a map with param. names as keys to values
func (e *ParametersExtractor) Parameters(req *http.Request) map[string]string {
	return mux.Vars(req)
}

// HasParameter - checking the param. is existing in request
func (e *ParametersExtractor) HasParameter(param string, req *http.Request) bool {
	if _, ok := e.Parameters(req)[param]; ok {
		return true
	}
	return false
}

// GetParameter - checking if the param. is existing in request, it will be returned
func (e *ParametersExtractor) GetParameter(param string, req *http.Request) (string, error) {
	if v, ok := e.Parameters(req)[param]; ok {
		return v, nil
	}
	return "", errs.NewFieldCannotBeEmptyError(param)
}
