package rpc

import (
	"log"
	"net/rpc"
	"reflect"

	"github.com/schidstorm/go-remote/lib"
	"github.com/schidstorm/go-remote/lib/connector"
)

type Client struct {
	Connector connector.Connector
	rpcClient *rpc.Client
}

func NewClient(ctor connector.Connector) *Client {
	client := new(Client)
	client.Connector = ctor
	return client
}

func (client *Client) Connect() error {
	rwc, err := client.Connector.Connect()

	if err != nil {
		return err
	}

	client.rpcClient = rpc.NewClient(rwc)

	return nil
}

func (client *Client) call(serviceMethod string, args interface{}, reply interface{}) error {
	if client.rpcClient == nil {
		log.Fatalln("Client is not initialized (.Initialize())")
	}
	log.Println("Calling remote " + serviceMethod)
	return client.rpcClient.Call(serviceMethod, args, reply)
}

func (client *Client) Call(serviceMethod string, args interface{}, TReturn reflect.Type) []lib.ErrorResult {
	log.Println("Calling " + serviceMethod)
	returnPtr := reflect.New(TReturn)
	err := client.call(serviceMethod, args, returnPtr.Interface())
	return []lib.ErrorResult{lib.ErrorResult{err, returnPtr.Interface()}}
}

func (client *Client) Close() error {
	return client.rpcClient.Close()
}
