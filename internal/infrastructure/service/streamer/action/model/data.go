package model

type StreamByIdData struct {
	ID    string
	Token string
}

type StreamByIdWithOffsetData struct {
	ID       string
	Token    string
	From     float64
	Duration float64
}
