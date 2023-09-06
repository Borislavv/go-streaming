package response

import (
	"encoding/json"
	"github.com/Borislavv/video-streaming/internal/domain/errs"
	"github.com/Borislavv/video-streaming/internal/domain/service"
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

type Response struct {
	logger service.Logger
}

func NewResponseService(logger service.Logger) *Response {
	return &Response{logger: logger}
}

func (r *Response) Respond(w http.ResponseWriter, dataOrErr any) {
	err, isErr := dataOrErr.(error)
	if isErr {
		publicErr, isPublicErr := dataOrErr.(errs.PublicError)
		if isPublicErr {
			r.logger.Info(publicErr.Error())
			w.WriteHeader(publicErr.Status())
			if _, err = w.Write(
				r.toBytes(
					NewErrorResponse(publicErr),
				),
			); err != nil {
				r.logger.Emergency(err)
			}
		} else {
			r.logger.Critical(err)
			w.WriteHeader(errs.DefaultErrorStatusCode)
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
