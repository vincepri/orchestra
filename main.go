package main

import (
	"os"

	"github.com/codegangsta/cli"
	"github.com/vinceprignano/orchestra/commands"
)

var app *cli.App

func main() {
	app = cli.NewApp()
	app.Name = "Orchestra"
	app.Usage = "Orchestrate Go Services"
	app.Author = "Vincenzo Prignano"
	app.Email = ""
	app.Commands = []cli.Command{
		*commands.StartCommand,
		*commands.StopCommand,
	}
	app.Version = "0.1"
	app.Run(os.Args)
}
