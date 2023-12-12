package di_interface

import "reflect"

type Container interface {
	Set(service any, alias reflect.Type) (self Container)
	Has(key reflect.Type) (has bool)
	Get(key reflect.Type) (service reflect.Value, notFoundErr error)
}
