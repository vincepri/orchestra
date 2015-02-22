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
			Usage: "Attach to services output after start",
		},
		cli.BoolFlag{
			Name:  "logs, l",
			Usage: "Start logging after start",
		},
	},
}

// StartAction starts all the services (or the specified ones)
func StartAction(c *cli.Context) {
	for _, service := range FilterServices(c) {
		spacing := strings.Repeat(" ", services.MaxServiceNameLength+2-len(service.Name))
		if service.Process == nil {
			rebuilt, err := buildAndStart(c, service)
			if err != nil {
				appendError(err)
				terminal.Stdout.Colorf("%s%s| @{r} error: @{|}%s\n", service.Name, spacing, err.Error())
			} else {
				rebuiltStatus := ""
				if rebuilt {
					rebuiltStatus = "rebuilt & "
				}
				terminal.Stdout.Colorf("%s%s| @{g} %sstarted\n", service.Name, spacing, rebuiltStatus)
			}
		} else {
			terminal.Stdout.Colorf("%s%s| @{c} %salready running\n", service.Name, spacing)
		}
	}
	if c.Bool("attach") || c.Bool("logs") {
		LogsAction(c)
	}
}

// startService takes a Service struct as input, creates a new log file in .orchestra,
// redirects the command stdout and stderr to the log file, configures the environment
// variables for the command and starts it. If cmd.Start() doesn't return any
// error, it will write a service.pid file in .orchestra
func buildAndStart(c *cli.Context, service *services.Service) (bool, error) {
	cmd := exec.Command(service.Name)

	rebuilt, err := buildService(service)
	if err != nil {
		return rebuilt, err
	}

	outputFile, err := os.Create(service.LogFilePath)
	if err != nil && os.IsNotExist(err) {
		return rebuilt, err
	}
	defer outputFile.Close()
	pidFile, err := os.Create(service.PidFilePath)
	if err != nil && os.IsNotExist(err) {
		return rebuilt, err
	}
	defer pidFile.Close()
	cmd.Stdout = outputFile
	cmd.Stderr = outputFile
	cmd.Env = GetEnvForService(c, service)

	if !c.Bool("attach") {
		cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	}

	if err := cmd.Start(); err != nil {
		return rebuilt, err
	}
	pidFile.WriteString(strconv.Itoa(cmd.Process.Pid))
	return rebuilt, nil
}

// buildService runs go install in the service directory
func buildService(service *services.Service) (bool, error) {
	cmd := exec.Command("go", "install", "-v")
	cmd.Dir = service.Path
	output := bytes.NewBuffer([]byte{})
	cmd.Stdout = output
	cmd.Stderr = output
	err := cmd.Start()
	if err != nil {
		return false, err
	}
	cmd.Wait()
	if !cmd.ProcessState.Success() {
		return false, fmt.Errorf("Failed to install service %s\n%s", service.Name, output.String())
	} else if output.Len() > 0 {
		return true, nil
	}
	return false, nil
}
