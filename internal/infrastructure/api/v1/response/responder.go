package response

import "net/http"

type Responder interface {
	Respond(w http.ResponseWriter, dataOrErr any)
}
