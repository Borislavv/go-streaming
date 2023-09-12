package request

import (
	"errors"
	"github.com/Borislavv/video-streaming/internal/domain/errs"
	"github.com/gorilla/mux"
	"net/http"
)

type ParametersExtractor struct {
}

func NewParametersExtractor() *ParametersExtractor {
	return &ParametersExtractor{}
}

// Parameters - will return a map with param. names as keys to values
func (e *ParametersExtractor) Parameters(req *http.Request) map[string]string {
	params := make(map[string]string)

	for name, param := range req.URL.Query() {
		params[name] = param[0]
	}
	for name, param := range mux.Vars(req) {
		params[name] = param
	}

	return params
}

// HasParameter - checking the param. is existing in request
func (e *ParametersExtractor) HasParameter(param string, req *http.Request) bool {
	if param == "" {
		return false
	}
	if _, ok := e.Parameters(req)[param]; ok {
		return true
	}
	return false
}

// GetParameter - checking if the param. is existing in request, it will be returned
func (e *ParametersExtractor) GetParameter(param string, req *http.Request) (string, error) {
	if param == "" {
		return "", errs.NewFieldCannotBeEmptyError(param)
	}
	if v, ok := e.Parameters(req)[param]; ok {
		return v, nil
	}
	return "", errors.New("parameter '" + param + "' not found into query or path")
}
