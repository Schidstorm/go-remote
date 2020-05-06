package echo

import "fmt"

type Echo struct {
}

func (echo *Echo) Echo(text, stdout *string) error {
	_, err := fmt.Print(text)
	return err
}
