package di

import (
	"fmt"
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
			s.Set(service, "")
		}
	}

	return s
}

func (s *ServiceContainer) Set(service any, alias reflect.Type) (self Container) {
	key := ""
	if alias != "" {
		key = alias
	} else {
		key = reflect.TypeOf(service).String()
	}
	s.container[key] = service
	return s
}

func (s *ServiceContainer) Has(key string) (has bool) {
	_, has = s.container[key]
	return has
}

func (s *ServiceContainer) Get(key string) (service any, notFoundErr error) {
	service, found := s.container[key]
	if !found {
		return nil, fmt.Errorf("service not found by key '%s'", key)
	}
	return service, nil
}
