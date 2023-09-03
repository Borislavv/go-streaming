package dto

type PaginationRequestDto struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
}
