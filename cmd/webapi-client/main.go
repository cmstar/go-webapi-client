package main

import (
	"flag"

	client "github.com/cmstar/go-webapi-client"
	"github.com/cmstar/go-webapi-client/slimauth_client"
)

var fConfigPath = flag.String("c", "", "specify the directory of config files")

func main() {
	flag.Parse()

	op := &client.MainWindowOption{
		ConfigPath: *fConfigPath,
		Clients: []client.Client{
			slimauth_client.NewClient(),
		},
	}
	win := client.NewMainWindow(op)
	win.ShowAndRun()
}
