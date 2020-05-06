package hello

//go:generate go run github.com/schidstorm/go-remote/generator

type Hello struct {
}

func (hw *Hello) HelloLocal(name *string, result *string) error {
	(*result) = "Hello " + *name + "!"
	return nil
}
