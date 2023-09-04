package response

import "encoding/json"

type ErrorResponse struct {
	Error error `json:"error"`
}

func NewErrorResponse(err error) *ErrorResponse {
	return &ErrorResponse{Error: err}
}

func (r *ErrorResponse) Wrap() ([]byte, error) {
	bytes, err := json.Marshal(ErrorResponse{Error: r.Error})
	if err != nil {
		return []byte{}, err
	}
	return bytes, nil
}
