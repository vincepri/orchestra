package commands

import (
	"fmt"
	"os"
	"os/exec"

	log "github.com/cihub/seelog"
	"github.com/codegangsta/cli"
	"github.com/vinceprignano/orchestra/services"
)

var StartCommand = &cli.Command{
	Name:   "start",
	Usage:  "Starts all the services",
	Action: StartAction,
}

func StartAction(c *cli.Context) {
	for name, service := range services.Registry {
		cmd := exec.Command(name)
		outputFile, err := os.Create(fmt.Sprintf("%s/%s.log", service.OrchestraPath, name))
		if err != nil && os.IsNotExist(err) {
			log.Error(err)
			continue
		}
		cmd.Stdout = outputFile
		cmd.Stderr = outputFile
		if err := cmd.Start(); err != nil {
			log.Error(err.Error())
			continue
		}
		log.Infof("Started %s", name)
	}
}
