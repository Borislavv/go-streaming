package model

type StreamByIdData struct {
	ID    string `json:"id"`
	Token string `json:"token"`
}

type StreamByIdWithOffsetData struct {
	ID       string  `json:"id"`
	Token    string  `json:"token"`
	From     float64 `json:"from"`
	Duration float64 `json:"duration"`
}
