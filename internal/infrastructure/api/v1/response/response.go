package response

import (
	"errors"
	"github.com/Borislavv/video-streaming/internal/domain/errs"
)

const InternalServerError = "Internal server error. Please, contact with service administrator."

type Response interface {
	Wrap() ([]byte, error)
}

// TODO нужно реализовать парсинг статус кода из ошибки
func Respond(data any, err error) ([]byte, error) {
	if err != nil {
		public, ok := err.(errs.PublicError)
		if ok {
			return NewErrorResponse(public).Wrap()
		} else {
			return NewErrorResponse(errors.New(InternalServerError)).Wrap()
		}
	}
	return NewDataResponse(data).Wrap()
}
