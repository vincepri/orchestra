package commands

import (
	"os"

	"github.com/cihub/seelog"
	"github.com/codegangsta/cli"
	"github.com/vinceprignano/orchestra/services"
)

var StopCommand = &cli.Command{
	Name:   "stop",
	Usage:  "Stops all the services",
	Action: StopAction,
}

func StopAction(c *cli.Context) {
	for name, service := range services.Registry {
		if service.Process != nil {
			err := service.Process.Kill()
			if err != nil {
				seelog.Error(err.Error())
				continue
			} else {
				defer os.Remove(service.PidFilePath)
			}
			seelog.Infof("Stopped %s", name)
		}
	}
}
