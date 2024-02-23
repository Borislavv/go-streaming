package dtointerface

type PaginatedRequest interface {
	GetPage() int
	GetLimit() int
}
