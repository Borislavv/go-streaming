package query

type Pagination interface {
	GetPage() int
	GetLimit() int
}
