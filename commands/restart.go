package commands

import "github.com/codegangsta/cli"

var RestartCommand = &cli.Command{
	Name:   "restart",
	Usage:  "Restarts all the services",
	Action: StartAction,
}

func RestartAction(c *cli.Context) {

}
