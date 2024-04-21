package di

import (
	"errors"
	diinterface "github.com/Borislavv/video-streaming/internal/infrastructure/di/interface"
	"reflect"
)

var notFoundError = errors.New("service not found in the DI container")

type Container struct {
	container map[reflect.Type]reflect.Value
}

func NewServiceContainer(services ...any) *Container {
	s := &Container{
		container: make(map[reflect.Type]reflect.Value),
	}

	if len(services) > 0 {
		for service := range services {
			s.Set(service, nil)
		}
	}

	return s
}

func (s *Container) Set(service any, alias reflect.Type) (self diinterface.Container) {
	if alias == nil || alias == reflect.TypeOf(nil) {
		alias = reflect.TypeOf(service)
	}
	s.container[alias] = reflect.ValueOf(service)
	return s
}

func (s *Container) Has(key reflect.Type) (has bool) {
	_, has = s.container[key]
	return has
}

func (s *Container) Get(key reflect.Type) (service reflect.Value, notFoundErr error) {
	service, found := s.container[key]
	if !found {
		return reflect.Value{}, notFoundError
	}
	return service, nil
}
