package cacher

type Cacher interface {
	Get(key string, fn func() interface{}) error
	Delete(key string) error
}
