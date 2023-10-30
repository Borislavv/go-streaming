package response

import (
	"encoding/json"
	"github.com/Borislavv/video-streaming/internal/domain/errors"
	"github.com/Borislavv/video-streaming/internal/domain/logger"
	"io"
	"net/http"
	"time"
)

type DataResponse struct {
	Data any `json:"data"`
}

func NewDataResponse(data any) DataResponse {
	return DataResponse{Data: data}
}

type ErrorResponse struct {
	Error error `json:"error"`
}

func NewErrorResponse(err error) ErrorResponse {
	return ErrorResponse{Error: err}
}

// Responder - response service interface
type Responder interface {
	Respond(w io.Writer, dataOrErr any)
}

// Response - response service
type Response struct {
	logger logger.Logger
}

func NewResponseService(logger logger.Logger) *Response {
	return &Response{logger: logger}
}

func (r *Response) Respond(w io.Writer, dataOrErr any) {
	err, isErr := dataOrErr.(error)
	if isErr {
		r.logger.Log(err)
		publicErr, isPublicErr := err.(errors.PublicError)
		if isPublicErr {
			// handle the case when write is http.ResponseWriter
			if httpWriter, ok := w.(http.ResponseWriter); ok {
				httpWriter.WriteHeader(publicErr.Status())
			}

			if _, err = w.Write(
				r.toBytes(
					NewErrorResponse(publicErr),
				),
			); err != nil {
				r.logger.Emergency(err)
			}
		} else {
			// handle the case when write is http.ResponseWriter
			if httpWriter, ok := w.(http.ResponseWriter); ok {
				httpWriter.WriteHeader(http.StatusInternalServerError)
			}

			if _, err = w.Write(
				r.toBytes(
					NewErrorResponse(
						errors.NewInternalServerError(),
					),
				),
			); err != nil {
				r.logger.Emergency(err)
			}
		}
		return
	}

	// building a new response data
	resp := NewDataResponse(dataOrErr)
	// writing a response data
	if _, err = w.Write(r.toBytes(resp)); err != nil {
		r.logger.Emergency(err)
	}
	// logging a response
	r.logResponse(resp)
}

func (r *Response) logResponse(resp DataResponse) {
	r.logger.LogData(
		&LoggableData{
			Date:         time.Now(),
			Type:         LogType,
			DataResponse: resp,
		},
	)
}

func (r *Response) toBytes(resp any) []byte {
	bytes, err := json.Marshal(resp)
	if err != nil {
		r.logger.Emergency(err)
	}
	return bytes
}
