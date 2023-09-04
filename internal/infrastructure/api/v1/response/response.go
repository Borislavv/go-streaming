package response

import (
	"errors"
	"github.com/Borislavv/video-streaming/internal/domain/errs"
	"log"
	"net/http"
)

const InternalServerErrorMessage = "Internal server error. Please, contact with service administrator."

type Response interface {
	Wrap() ([]byte, error)
}

func RespondData(w http.ResponseWriter, data any) {
	if _, err := w.Write(NewDataResponse(data).Wrap()); err != nil {
		log.Fatalln(err)
	}
}

func RespondError(w http.ResponseWriter, err error) {
	public, ok := err.(errs.PublicError)
	if ok {
		w.WriteHeader(public.Status())
		if _, err = w.Write(NewErrorResponse(public).Wrap()); err != nil {
			log.Fatalln(err)
		}
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		if _, err = w.Write(NewErrorResponse(errors.New(InternalServerErrorMessage)).Wrap()); err != nil {
			log.Fatalln(err)
		}
	}
}
