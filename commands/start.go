package commands

import (
	"os"
	"os/exec"
	"strconv"
	"strings"

	log "github.com/cihub/seelog"
	"github.com/codegangsta/cli"
	"github.com/vinceprignano/orchestra/config"
	"github.com/vinceprignano/orchestra/services"
	"github.com/wsxiaoys/terminal"
)

var StartCommand = &cli.Command{
	Name:         "start",
	Usage:        "Starts all the services",
	Action:       StartAction,
	BashComplete: ServicesBashComplete,
}

// StartAction starts all the services (or the specified ones)
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

// startService takes a Service struct as input, creates a new log file in .orchestra,
// redirects the command stdout and stderr to the log file, configures the environment
// variables for the command and starts it. If cmd.Start() doesn't return any
// error, it will write a service.pid file in .orchestra
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
	cmd.Env = config.GetEnvironmentVars(service)
	if err := cmd.Start(); err != nil {
		return err
	}
	pidFile.WriteString(strconv.Itoa(cmd.Process.Pid))
	return nil
}
