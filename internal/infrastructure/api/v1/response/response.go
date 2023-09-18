package response

import (
	"encoding/json"
	"github.com/Borislavv/video-streaming/internal/domain/enum"
	"github.com/Borislavv/video-streaming/internal/domain/errs"
	"github.com/Borislavv/video-streaming/internal/domain/logger"
	"io"
	"net/http"
	"time"
)

const LogType = "response"

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
		publicErr, isPublicErr := err.(errs.PublicError)
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
						errs.NewInternalServerError(),
					),
				),
			); err != nil {
				r.logger.Emergency(err)
			}
		}
		return
	}

	resp := NewDataResponse(dataOrErr)
	if _, err = w.Write(r.toBytes(resp)); err != nil {
		r.logger.Emergency(err)
	}
	r.logResponse(resp)
}

func (r *Response) logResponse(resp DataResponse) {
	requestID := ""
	if uniqReqID := r.logger.Context().Value(enum.UniqueRequestIdKey); uniqReqID != nil {
		if strUniqReqID, ok := uniqReqID.(string); ok {
			requestID = strUniqReqID
		}
	}

	respData := struct {
		Date      time.Time `json:"date"`
		RequestId string    `json:"requestID,omitempty"`
		Type      string    `json:"type"`
		DataResponse
	}{
		Date:         time.Now(),
		RequestId:    requestID,
		Type:         LogType,
		DataResponse: resp,
	}

	r.logger.LogData(respData)
}

func (r *Response) toBytes(resp any) []byte {
	bytes, err := json.Marshal(resp)
	if err != nil {
		r.logger.Emergency(err)
	}
	return bytes
}
