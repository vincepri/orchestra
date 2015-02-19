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
	for _, service := range services.Registry {
		startService(service)
	}
}

func startService(service *services.Service) {
	cmd := exec.Command(service.Name)
	outputFile, err := os.Create(service.LogFilePath)
	if err != nil && os.IsNotExist(err) {
		log.Error(err)
		return
	}
	defer outputFile.Close()
	pidFile, err := os.Create(service.PidFilePath)
	if err != nil && os.IsNotExist(err) {
		log.Error(err)
		return
	}
	defer pidFile.Close()
	cmd.Stdout = outputFile
	cmd.Stderr = outputFile
	if err := cmd.Start(); err != nil {
		log.Error(err.Error())
		return
	}
	pidFile.WriteString(strconv.Itoa(cmd.Process.Pid))
	log.Infof("Started %s", service.Name)
}
