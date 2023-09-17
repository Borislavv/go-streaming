package dto

type PaginationRequestDTO struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
}

func (p PaginationRequestDTO) GetPage() int {
	return p.Page
}

func (p PaginationRequestDTO) GetLimit() int {
	return p.Limit
}
