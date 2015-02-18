package main

import "github.com/codegangsta/cli"

var StopCommand = &cli.Command{
	Name:   "stop",
	Usage:  "Stops all the services",
	Action: StopAction,
}

func StopAction(c *cli.Context) {

}
