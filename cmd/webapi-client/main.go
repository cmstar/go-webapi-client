package main

import (
	client "github.com/cmstar/go-webapi-client"
	"github.com/cmstar/go-webapi-client/slimauth_client"
)

func main() {
	client.RunClients([]client.Client{
		slimauth_client.NewClient(),
	})
}
