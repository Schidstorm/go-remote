package shell

//go:generate go run github.com/schidstorm/go-remote/generator

import (
	"log"
	"os/exec"
)

type Shell struct {
	Remote ShellRemote
}

type ShellOptions struct {
	Program   string
	Arguments []string
}

type ShellResult struct {
	Stdout []byte
}

func (s *Shell) Run(opts *ShellOptions, result *ShellResult) error {
	if opts == nil {
		log.Fatalln("Options are nil")
	}

	cmd := exec.Command(opts.Program, opts.Arguments...)
	stdout, err := cmd.Output()
	result.Stdout = stdout
	return err
}
