package response

import (
	"encoding/json"
	"log"
)

type DataResponse struct {
	Data any `json:"data"`
}

func NewDataResponse(data any) *DataResponse {
	return &DataResponse{Data: data}
}

func (r *DataResponse) Wrap() []byte {
	bytes, err := json.Marshal(DataResponse{Data: r.Data})
	if err != nil {
		log.Fatalln(err)
	}
	return bytes
}
