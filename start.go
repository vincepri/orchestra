package main

import "github.com/codegangsta/cli"

var StartCommand = &cli.Command{
	Name:   "start",
	Usage:  "Starts all the services",
	Action: StartAction,
}

func StartAction(c *cli.Context) {

}
