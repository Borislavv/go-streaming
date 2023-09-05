package response

import (
	"errors"
	"github.com/Borislavv/video-streaming/internal/domain/errs"
	"github.com/Borislavv/video-streaming/internal/domain/service"
	"net/http"
)

const InternalServerErrorMessage = "Internal server error. Please, contact with service administrator."

type Response struct {
	logger service.Logger
}

func NewResponse(logger service.Logger) *Response {
	return &Response{logger: logger}
}

func (r *Response) RespondData(w http.ResponseWriter, data any) {
	if _, err := w.Write(NewDataResponse(data).Wrap()); err != nil {
		r.logger.Emergency(err)
		return
	}
}

func (r *Response) RespondError(w http.ResponseWriter, err error) {
	public, ok := err.(errs.PublicError)
	if ok {
		w.WriteHeader(public.Status())
		if _, err = w.Write(NewErrorResponse(public).Wrap()); err != nil {
			r.logger.Emergency(err)
			return
		}
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		if _, err = w.Write(NewErrorResponse(errors.New(InternalServerErrorMessage)).Wrap()); err != nil {
			r.logger.Emergency(err)
			return
		}
	}
}
