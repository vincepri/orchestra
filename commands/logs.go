package commands

import "github.com/codegangsta/cli"

var LogsCommand = &cli.Command{
	Name:   "logs",
	Usage:  "Aggregate services logs",
	Action: StartAction,
}

func LogsAction(c *cli.Context) {

}
