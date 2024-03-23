package queryinterface

type Pagination interface {
	GetPage() int
	GetLimit() int
}
