package _interface

type PaginatedRequest interface {
	GetPage() int
	GetLimit() int
}
