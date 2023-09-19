package response

import "time"

const LogType = "response"

type LoggableData struct {
	Date  time.Time `json:"date"`
	ReqID string    `json:"requestID,omitempty"`
	Type  string    `json:"type"`
	DataResponse
}

func (r *LoggableData) RequestID() string {
	return r.ReqID
}

func (r *LoggableData) SetRequestID(id string) {
	r.ReqID = id
}
