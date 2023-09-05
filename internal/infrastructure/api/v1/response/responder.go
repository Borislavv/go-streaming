package response

import "net/http"

type Responder interface {
	RespondData(w http.ResponseWriter, data any)
	RespondError(w http.ResponseWriter, err error)
}
