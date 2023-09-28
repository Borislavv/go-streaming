package model

type StreamByIdData struct {
	ID string
}

type StreamByIdWithOffsetData struct {
	ID       string
	From     float64
	Duration float64
}
