package dto_interface

type PaginatedRequest interface {
	GetPage() int
	GetLimit() int
}
