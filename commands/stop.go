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
	for _, service := range services.Registry {
		killService(service)
	}
}

func killService(service *services.Service) {
	if service.Process != nil {
		err := service.Process.Kill()
		defer os.Remove(service.PidFilePath)
		if err != nil {
			seelog.Error(err.Error())
			return
		}
		seelog.Infof("Stopped %s", service.Name)
	}
}
