package request

import (
	"net/http"
	"time"
)

const LogType = "request"

type LoggableData struct {
	Date       time.Time         `json:"date"`
	ReqID      string            `json:"requestID"`
	Type       string            `json:"type"`
	Method     string            `json:"method"`
	URL        string            `json:"URL"`
	Header     http.Header       `json:"header"`
	RemoteAddr string            `json:"remoteAddr"`
	Params     map[string]string `json:"params"`
}

func (r *LoggableData) RequestID() string {
	return r.ReqID
}

func (r *LoggableData) SetRequestID(id string) {
	r.ReqID = id
}
