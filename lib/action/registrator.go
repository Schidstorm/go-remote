package action

type Registrator interface {
	RegisterName(name string, object interface{}) error
}
