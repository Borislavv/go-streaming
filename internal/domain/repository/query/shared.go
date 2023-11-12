package query

type Paginated struct {
	page  int
	limit int
}

func NewPagination(page int, limit int) *Paginated {
	return &Paginated{
		page:  page,
		limit: limit,
	}
}
func (q *Paginated) GetPage() int {
	return q.page
}
func (q *Paginated) GetLimit() int {
	return q.limit
}
