package main

import (
	"fmt"
	"os"

	"github.com/schidstorm/go-remote/lib/action/shell"
	"github.com/schidstorm/go-remote/lib/connector"
	"github.com/schidstorm/go-remote/lib/listener"
	"github.com/schidstorm/go-remote/lib/rpc"
)

func main() {
	if rpc.IsServerMode() {
		remote()
	} else {
		local()
	}
}

func remote() {
	server := rpc.NewServer()
	server.Listen(&listener.Io{})
}

func local() {
	connector1, err := connector.SshWithPassphrase(os.Getenv("SERVER_SSH_ENDPOINT"), os.Getenv("SERVER_SSH_PASSWORD"), os.Getenv("SERVER_SSH_USER"))
	if err != nil {
		panic(err)
	}
	client1 := rpc.NewClient(connector1)

	err = client1.Connect()
	if err != nil {
		panic(err)
	}

	s := shell.Shell{}
	results := s.Remote.Run(client1, &shell.ShellOptions{
		Program:   "cat",
		Arguments: []string{"/etc/passwd"},
	})

	for _, res := range results {
		fmt.Println(string(res.Data.Stdout))
	}
}
