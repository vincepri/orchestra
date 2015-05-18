package commands

import (
	"errors"
	"os"
	"os/exec"
	"strings"

	"github.com/mondough/orchestra/services"
	"github.com/codegangsta/cli"
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
		cli.BoolFlag{
			Name: "race, r",
		},
	},
}

// StartAction starts all the services (or the specified ones)
func TestAction(c *cli.Context) {
	for _, service := range FilterServices(c) {
		spacing := strings.Repeat(" ", services.MaxServiceNameLength+2-len(service.Name))
		success, err := testService(c, service)
		if err != nil {
			appendError(err)
			terminal.Stdout.Colorf("%s%s| @{r} error: @{|}%s\n", service.Name, spacing, err.Error())
		} else if !success {
			appendError(errors.New("Test Failed"))
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
	cmdArgs := []string{"test"}
	if c.Bool("verbose") {
		cmdArgs = append(cmdArgs, "-v")
	}
	if c.Bool("race") {
		cmdArgs = append(cmdArgs, "--race")
	}
	cmdArgs = append(cmdArgs, "./...")
	cmd := exec.Command("go", cmdArgs...)
	cmd.Dir = service.Path
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = GetEnvForService(c, service)
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
