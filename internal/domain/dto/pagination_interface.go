package dto

type PaginatedRequest interface {
	GetPage() int
	GetLimit() int
}
