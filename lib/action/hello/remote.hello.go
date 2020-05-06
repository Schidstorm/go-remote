
package hello

import (
	"errors"
	"reflect"

	"github.com/schidstorm/go-remote/lib"
)

type receiverType = HelloRemote
type inputType = *string
type resultType = *string

type HelloError struct {
	Data resultType
	Error error
}

type errorResultType = HelloError

type HelloRemote struct {}


func (s *receiverType) HelloLocal(client lib.Callable, opts inputType) []errorResultType {
	interfaceResults := client.Call("Hello.HelloLocal", opts, reflect.TypeOf((resultType)(nil)).Elem())

	results := []errorResultType{}
	for _, interfaceResult := range interfaceResults {
		switch v := interfaceResult.Data.(type) {
		case resultType:
			results = append(results, errorResultType{v, interfaceResult.Error})
		default:
			err := interfaceResult.Error
			if err == nil {
				err = errors.New("Unknown result type")
			}
			results = append(results, errorResultType{nil, err})
		}
	}

	return results
}
			