package commands

import (
	"github.com/cihub/seelog"
	"github.com/codegangsta/cli"
	"github.com/vinceprignano/orchestra/services"
)

var PsCommand = &cli.Command{
	Name:   "ps",
	Usage:  "Outputs the status of all services",
	Action: PsAction,
}

func PsAction(c *cli.Context) {
	for name, service := range services.Registry {
		if service.Process != nil {
			seelog.Infof("Service %s is RUNNING", name)
		} else {
			seelog.Infof("Service %s is not running", name)
		}
	}
}
