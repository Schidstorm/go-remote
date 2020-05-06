
package shell

import (
	"errors"
	"reflect"

	"github.com/schidstorm/go-remote/lib"
)

type receiverType = ShellRemote
type inputType = *ShellOptions
type resultType = *ShellResult

type ShellError struct {
	Data resultType
	Error error
}

type errorResultType = ShellError

type ShellRemote struct {}


func (s *receiverType) Run(client lib.Callable, opts inputType) []errorResultType {
	interfaceResults := client.Call("Shell.Run", opts, reflect.TypeOf((resultType)(nil)).Elem())

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
			