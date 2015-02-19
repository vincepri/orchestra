package commands

import "github.com/codegangsta/cli"

var VendorsCommand = &cli.Command{
	Name:   "vendors",
	Usage:  "Starts all the vendors dependecies using crane",
	Action: VendorsAction,
}

func VendorsAction(c *cli.Context) {

}
