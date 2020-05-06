package action

import (
	"github.com/schidstorm/go-remote/lib/action/hello"
	"github.com/schidstorm/go-remote/lib/action/shell"
)

type iface = interface{}

var builtinsMap = map[string]iface{
	"Hello": new(hello.Hello),
	"Shell": new(shell.Shell),
}

type Builtins struct {
}

func (b *Builtins) HandleRegistration(registrator Registrator) error {
	for name, obj := range builtinsMap {
		err := registrator.RegisterName(name, obj)
		if err != nil {
			return err
		}
	}

	return nil
}
