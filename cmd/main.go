package main

import (
	client "github.com/cmstar/go-webapi-client"
	"github.com/cmstar/go-webapi-client/slimauth_client"
)

func main() {
	op := &client.RunOption{
		Clients: []client.Client{
			slimauth_client.NewClient(),
		},
	}
	client.Run(op)
}
