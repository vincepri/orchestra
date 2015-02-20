package main

import (
	"fmt"
	"os"
	"path"
	"path/filepath"

	log "github.com/cihub/seelog"
	"github.com/codegangsta/cli"
	"github.com/vinceprignano/orchestra/commands"
	"github.com/vinceprignano/orchestra/config"
	"github.com/vinceprignano/orchestra/services"
)

var app *cli.App

func main() {
	defer log.Flush()
	app = cli.NewApp()
	app.Name = "Orchestra"
	app.Usage = "Orchestrate Go Services"
	app.Author = "Vincenzo Prignano"
	app.Email = ""
	app.EnableBashCompletion = true
	app.Commands = []cli.Command{
		*commands.StartCommand,
		*commands.StopCommand,
		*commands.LogsCommand,
		*commands.RestartCommand,
		*commands.PsCommand,
		*commands.TestCommand,
	}
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "config, c",
			Value:  "orchestra.yml",
			Usage:  "Specify a different config file to use",
			EnvVar: "ORCHESTRA_CONFIG",
		},
	}
	// init checks for an existing orchestra.yml in the current working directory
	// and creates a new .orchestra directory (if doesn't exist)
	app.Before = func(c *cli.Context) error {
		config.ConfigPath, _ = filepath.Abs(c.GlobalString("config"))
		if _, err := os.Stat(config.ConfigPath); os.IsNotExist(err) {
			fmt.Printf("No %s found. Have you specified the right directory?\n", c.GlobalString("config"))
			os.Exit(1)
		}
		services.ProjectPath, _ = path.Split(config.ConfigPath)
		services.OrchestraServicePath = services.ProjectPath + "/.orchestra"

		if err := os.Mkdir(services.OrchestraServicePath, 0766); err != nil && os.IsNotExist(err) {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		config.ParseGlobalConfig()
		services.Init()
		return nil
	}
	app.Version = "0.1"
	app.Run(os.Args)
	if commands.HasErrors() {
		os.Exit(1)
	}
}
