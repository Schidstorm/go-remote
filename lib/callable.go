package lib

import "reflect"

type ErrorResult struct {
	Error error
	Data  interface{}
}

type Callable interface {
	Call(serviceMethod string, args interface{}, TReturn reflect.Type) []ErrorResult
}
