package dto

type PaginationRequestDto struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
}

func (p PaginationRequestDto) GetPage() int {
	return p.Page
}

func (p PaginationRequestDto) GetLimit() int {
	return p.Limit
}
