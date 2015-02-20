package commands

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/codegangsta/cli"
	"github.com/vinceprignano/orchestra/config"
	"github.com/vinceprignano/orchestra/services"
	"github.com/wsxiaoys/terminal"
)

var StartCommand = &cli.Command{
	Name:         "start",
	Usage:        "Starts all the services",
	Action:       BeforeAfterWrapper(StartAction),
	BashComplete: ServicesBashComplete,
}

// StartAction starts all the services (or the specified ones)
func StartAction(c *cli.Context) {
	for _, service := range FilterServices(c) {
		spacing := strings.Repeat(" ", services.MaxServiceNameLength+2-len(service.Name))
		if service.Process == nil {
			err := startService(c, service)
			if err != nil {
				appendError(err)
				terminal.Stdout.Colorf("%s%s| @{r} error: @{|}%s\n", service.Name, spacing, err.Error())
			} else {
				terminal.Stdout.Colorf("%s%s| @{g} started\n", service.Name, spacing)
			}
		} else {
			terminal.Stdout.Colorf("%s%s| @{c} already running\n", service.Name, spacing)
		}
	}
}

// startService takes a Service struct as input, creates a new log file in .orchestra,
// redirects the command stdout and stderr to the log file, configures the environment
// variables for the command and starts it. If cmd.Start() doesn't return any
// error, it will write a service.pid file in .orchestra
func startService(c *cli.Context, service *services.Service) error {
	err := buildService(service)
	if err != nil {
		return err
	}
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
	cmd.Env = config.GetEnvForService(c, service)
	if err := cmd.Start(); err != nil {
		return err
	}
	pidFile.WriteString(strconv.Itoa(cmd.Process.Pid))
	return nil
}

func buildService(service *services.Service) error {
	cmd := exec.Command("go", "get")
	cmd.Dir = service.Path
	output := bytes.NewBuffer([]byte{})
	cmd.Stdout = output
	cmd.Stderr = output
	err := cmd.Start()
	if err != nil {
		return err
	}
	cmd.Wait()
	if !cmd.ProcessState.Success() {
		return fmt.Errorf("Failed to build service %s\n%s", service.Name, output.String())
	}
	return nil
}
