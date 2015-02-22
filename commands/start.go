package commands

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"

	"github.com/codegangsta/cli"
	"github.com/vinceprignano/orchestra/services"
	"github.com/wsxiaoys/terminal"
)

var StartCommand = &cli.Command{
	Name:         "start",
	Usage:        "Starts all the services",
	Action:       BeforeAfterWrapper(StartAction),
	BashComplete: ServicesBashComplete,
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "attach, a",
			Usage: "Attach logs right after starting",
		},
	},
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
	if c.Bool("attach") {
		LogsAction(c)
	}
}

// startService takes a Service struct as input, creates a new log file in .orchestra,
// redirects the command stdout and stderr to the log file, configures the environment
// variables for the command and starts it. If cmd.Start() doesn't return any
// error, it will write a service.pid file in .orchestra
func startService(c *cli.Context, service *services.Service) error {
	var cmd *exec.Cmd

	cmd = exec.Command(service.Name)

	if c.Bool("goget") {
		err := buildService(service)
		if err != nil {
			return err
		}
	} else {
		err := installService(service)
		if err != nil {
			return err
		}
		if c.Bool("gorun") {
			cmd = exec.Command("go", "run", "main.go")
			cmd.Dir = service.Path
		}
	}

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
	cmd.Env = GetEnvForService(c, service)
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	if err := cmd.Start(); err != nil {
		return err
	}
	pidFile.WriteString(strconv.Itoa(cmd.Process.Pid))
	return nil
}

// buildService runs go get ./... in the service directory
func buildService(service *services.Service) error {
	cmd := exec.Command("go", "get", "./...")
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
		return fmt.Errorf("Failed to `go get` service %s\n%s", service.Name, output.String())
	}
	return nil
}

// installService runs go install in the service directory
func installService(service *services.Service) error {
	cmd := exec.Command("go", "install")
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
		return fmt.Errorf("Failed to install service %s\n%s", service.Name, output.String())
	}
	return nil
}
