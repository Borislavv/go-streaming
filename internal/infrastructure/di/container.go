package di

import (
	"fmt"
	di_interface "github.com/Borislavv/video-streaming/internal/infrastructure/di/interface"
	"reflect"
)

type ServiceContainer struct {
	container map[reflect.Type]reflect.Value
}

func NewServiceContainer(services ...any) *ServiceContainer {
	s := &ServiceContainer{
		container: make(map[reflect.Type]reflect.Value),
	}

	if len(services) > 0 {
		for service := range services {
			s.Set(service, nil)
		}
	}

	return s
}

func (s *ServiceContainer) Set(service any, alias reflect.Type) (self di_interface.Container) {
	if alias == nil || alias == reflect.TypeOf(nil) {
		alias = reflect.TypeOf(service)
	}
	s.container[alias] = reflect.ValueOf(service)
	return s
}

func (s *ServiceContainer) Has(key reflect.Type) (has bool) {
	_, has = s.container[key]
	return has
}

func (s *ServiceContainer) Get(key reflect.Type) (service reflect.Value, notFoundErr error) {
	service, found := s.container[key]
	if !found {
		return reflect.Value{}, fmt.Errorf("service not found by key '%s'", key)
	}
	return service, nil
}
