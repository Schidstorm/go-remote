package rpc

import (
	"log"
	nrpc "net/rpc"
	"os"

	"github.com/schidstorm/go-remote/lib/action"
	"github.com/schidstorm/go-remote/lib/listener"
)

type RpcServer struct {
	server   *nrpc.Server
	Registry action.Registry
}

func IsServerMode() bool {
	return len(os.Args) > 1 && os.Args[1] == "--listen"
}

func NewServer() *RpcServer {
	log.Println("Creating RPC listener")
	server := new(RpcServer)
	server.server = nrpc.NewServer()
	server.Registry.Children = append(server.Registry.Children, new(action.Builtins))
	return server
}

func (server *RpcServer) Listen(listener listener.Listener) error {
	log.Println("Registering actions.")
	server.Registry.HandleRegistration(server.server)
	log.Println("Listening on stdin")
	rwc, err := listener.Listen()
	if err != nil {
		return err
	}

	server.server.ServeConn(rwc)
	return nil
}
