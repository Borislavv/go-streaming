package query_interface

type Pagination interface {
	GetPage() int
	GetLimit() int
}
