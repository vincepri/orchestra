package main

import (
	"fmt"
	"os"

	log "github.com/cihub/seelog"
	"github.com/codegangsta/cli"
	"github.com/vinceprignano/orchestra/commands"
	"github.com/vinceprignano/orchestra/services"
)

var app *cli.App

// init check for an existing orchestra.yml in the current working directory
// and creates a new .orchestra directory (if doesn't exist)
func init() {
	wd, _ := os.Getwd()
	if _, err := os.Stat(fmt.Sprintf("%s/orchestra.yml", wd)); os.IsNotExist(err) {
		fmt.Println("No orchestra.yml found. Are you in the right directory?")
		os.Exit(1)
	}
	if err := os.Mkdir(fmt.Sprintf("%s/.orchestra", wd), 0766); err != nil && os.IsNotExist(err) {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func main() {
	defer log.Flush()
	services.Init()
	app = cli.NewApp()
	app.Name = "Orchestra"
	app.Usage = "Orchestrate Go Services"
	app.Author = "Vincenzo Prignano"
	app.Email = ""
	app.Commands = []cli.Command{
		*commands.StartCommand,
		*commands.StopCommand,
		*commands.LogsCommand,
		*commands.RestartCommand,
	}
	app.Version = "0.1"
	app.Run(os.Args)
}
