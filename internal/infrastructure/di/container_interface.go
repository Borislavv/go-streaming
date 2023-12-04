package di

type Container interface {
	Set(service any, alias string) (self Container)
	Has(key string) (has bool)
	Get(key string) (service any, notFoundErr error)
}
