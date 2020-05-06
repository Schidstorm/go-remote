package rpc

import (
	"log"
	"reflect"

	"github.com/schidstorm/go-remote/lib"
)

type ClientPool struct {
	Clients []*Client
}

func NewClientPool(clients []*Client) *ClientPool {
	return &ClientPool{
		Clients: clients,
	}
}

func (pool *ClientPool) Connect() error {
	for _, client := range pool.Clients {
		err := client.Connect()
		if err != nil {
			return err
		}
	}

	return nil
}

func (pool *ClientPool) Call(serviceMethod string, args interface{}, TReturn reflect.Type) []lib.ErrorResult {
	log.Println("Calling " + serviceMethod)
	results := []lib.ErrorResult{}
	for _, client := range pool.Clients {
		results = append(results, client.Call(serviceMethod, args, TReturn)...)
	}

	return results
}

func (pool *ClientPool) Close() error {
	for _, client := range pool.Clients {
		err := client.Close()
		if err != nil {
			return err
		}
	}

	return nil
}
