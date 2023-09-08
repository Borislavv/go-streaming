package response

import (
	"encoding/json"
	"github.com/Borislavv/video-streaming/internal/domain/errs"
	"github.com/Borislavv/video-streaming/internal/domain/service"
	"io"
	"net/http"
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
	logger service.Logger
}

func NewResponseService(logger service.Logger) *Response {
	return &Response{logger: logger}
}

func (r *Response) Respond(w io.Writer, dataOrErr any) {
	err, isErr := dataOrErr.(error)
	if isErr {
		publicErr, isPublicErr := dataOrErr.(errs.PublicError)
		if isPublicErr {
			r.logger.Info(publicErr.Error())

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
			r.logger.Critical(err)

			// handle the case when write is http.ResponseWriter
			if httpWriter, ok := w.(http.ResponseWriter); ok {
				httpWriter.WriteHeader(errs.DefaultErrorStatusCode)
			}

			if _, err = w.Write(
				r.toBytes(
					NewErrorResponse(
						errs.NewError(
							errs.DefaultErrorMessage,
							errs.DefaultErrorType,
							errs.DefaultErrorStatusCode,
						),
					),
				),
			); err != nil {
				r.logger.Emergency(err)
			}
		}
		return
	}
	if _, err = w.Write(r.toBytes(NewDataResponse(dataOrErr))); err != nil {
		r.logger.Emergency(err)
	}
}

func (r *Response) toBytes(resp any) []byte {
	bytes, err := json.Marshal(resp)
	if err != nil {
		r.logger.Emergency(err)
	}
	return bytes
}
