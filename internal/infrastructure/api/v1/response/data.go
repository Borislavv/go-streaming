package response

import "encoding/json"

type DataResponse struct {
	Data any `json:"data"`
}

func NewDataResponse(data any) *DataResponse {
	return &DataResponse{Data: data}
}

func (r *DataResponse) Wrap() ([]byte, error) {
	bytes, err := json.Marshal(DataResponse{Data: r.Data})
	if err != nil {
		return []byte{}, err
	}
	return bytes, nil
}
