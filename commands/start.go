package commands

import (
	"os"
	"os/exec"
	"strconv"
	"strings"

	log "github.com/cihub/seelog"
	"github.com/codegangsta/cli"
	"github.com/vinceprignano/orchestra/services"
	"github.com/wsxiaoys/terminal"
)

var StartCommand = &cli.Command{
	Name:   "start",
	Usage:  "Starts all the services",
	Action: StartAction,
}

func StartAction(c *cli.Context) {
	for _, service := range FilterServices(c) {
		spacing := strings.Repeat(" ", services.MaxServiceNameLength+2-len(service.Name))
		err := startService(service)
		if err != nil {
			log.Error(err)
		} else {
			terminal.Stdout.Colorf("%s%s| @{g} started\n", service.Name, spacing)
		}
	}
}

func startService(service *services.Service) error {
	cmd := exec.Command(service.Name)
	outputFile, err := os.Create(service.LogFilePath)
	if err != nil && os.IsNotExist(err) {
		return err
	}
	defer outputFile.Close()
	pidFile, err := os.Create(service.PidFilePath)
	if err != nil && os.IsNotExist(err) {
		return err
	}
	defer pidFile.Close()
	cmd.Stdout = outputFile
	cmd.Stderr = outputFile
	cmd.Env = services.OrchestraConfig.Environment
	if err := cmd.Start(); err != nil {
		return err
	}
	pidFile.WriteString(strconv.Itoa(cmd.Process.Pid))
	return nil
}
