package commands

import (
	"os"
	"os/exec"
	"strings"

	"github.com/cihub/seelog"
	"github.com/codegangsta/cli"
	"github.com/vinceprignano/orchestra/config"
	"github.com/vinceprignano/orchestra/services"
	"github.com/wsxiaoys/terminal"
)

var TestCommand = &cli.Command{
	Name:         "test",
	Usage:        "Runs go test ./... for every service",
	Action:       BeforeAfterWrapper(TestAction),
	BashComplete: ServicesBashComplete,
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name: "verbose, v",
		},
	},
}

// StartAction starts all the services (or the specified ones)
func TestAction(c *cli.Context) {
	for _, service := range FilterServices(c) {
		spacing := strings.Repeat(" ", services.MaxServiceNameLength+2-len(service.Name))
		success, err := testService(c, service)
		if err != nil {
			terminal.Stdout.Colorf("%s%s| @{r} error: @{|}%s\n", service.Name, spacing, err.Error())
		} else if !success {
			terminal.Stdout.Colorf("%s%s| @{r} FAILED\n", service.Name, spacing)
		} else {
			terminal.Stdout.Colorf("%s%s| @{g} PASS\n", service.Name, spacing)
		}
	}
}

// startService takes a Service struct as input, creates a new log file in .orchestra,
// redirects the command stdout and stderr to the log file, configures the environment
// variables for the command and starts it. If cmd.Start() doesn't return any
// error, it will write a service.pid file in .orchestra
func testService(c *cli.Context, service *services.Service) (bool, error) {
	var cmd *exec.Cmd
	seelog.Info(c.Bool("verbose"))
	if c.Bool("verbose") {
		cmd = exec.Command("go", "test", "-v", "./...")
	} else {
		cmd = exec.Command("go", "test", "./...")
	}
	cmd.Dir = service.Path
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = config.GetEnvForService(c, service)
	err := cmd.Start()
	if err != nil {
		return false, err
	}
	cmd.Wait()
	if !cmd.ProcessState.Success() {
		return false, nil
	}
	return true, nil
}
