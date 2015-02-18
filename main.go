package main

import (
	"os"

	"github.com/codegangsta/cli"
)

var app *cli.App

func main() {
	app = cli.NewApp()
	app.Name = "Orchestra"
	app.Usage = "Orchestrate Go Services"
	app.Commands = []cli.Command{
		*StartCommand,
		*StopCommand,
	}
	app.Version = "0.1"
	app.Run(os.Args)
}
