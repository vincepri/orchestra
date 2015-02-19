package commands

import (
	"os"
	"os/exec"
	"strconv"

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
		outputFile, err := os.Create(service.LogFilePath)
		if err != nil && os.IsNotExist(err) {
			log.Error(err)
			continue
		}
		defer outputFile.Close()
		pidFile, err := os.Create(service.PidFilePath)
		if err != nil && os.IsNotExist(err) {
			log.Error(err)
			continue
		}
		defer pidFile.Close()
		cmd.Stdout = outputFile
		cmd.Stderr = outputFile
		if err := cmd.Start(); err != nil {
			log.Error(err.Error())
			continue
		}
		pidFile.WriteString(strconv.Itoa(cmd.Process.Pid))
		log.Infof("Started %s", name)
	}
}
