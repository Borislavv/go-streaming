package response

import (
	"encoding/json"
	"log"
)

type ErrorResponse struct {
	Error error `json:"error"`
}

func NewErrorResponse(err error) *ErrorResponse {
	return &ErrorResponse{Error: err}
}

func (r *ErrorResponse) Wrap() []byte {
	bytes, err := json.Marshal(ErrorResponse{Error: r.Error})
	if err != nil {
		log.Fatalln(err)
	}
	return bytes
}
