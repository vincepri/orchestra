package commands

import (
	"bytes"
	"fmt"
	"os/exec"
	"runtime"
	"strings"

	"github.com/codegangsta/cli"
	"github.com/wsxiaoys/terminal"

	"github.com/mondough/orchestra/services"
)

var BuildCommand = &cli.Command{
	Name:         "build",
	Usage:        "Build service(s)",
	Action:       BeforeAfterWrapper(BuildAction),
	BashComplete: ServicesBashComplete,
}

func BuildAction(c *cli.Context) {
	worker := func(service *services.Service) func() {
		return func() {
			buildService(c, service)
		}
	}

	pool := make(workerPool, runtime.NumCPU())
	for _, service := range FilterServices(c) {
		pool.Do(worker(service))
	}
	pool.Drain()
}

func buildService(c *cli.Context, service *services.Service) {
	spacing := strings.Repeat(" ", services.MaxServiceNameLength+2-len(service.Name))

	cmd := exec.Command("nice", "-n", niceness, "go", "build", "-v")
	cmd.Dir = service.Path
	output := new(bytes.Buffer)
	cmd.Stdout = output
	cmd.Stderr = output
	if cmd.Run() != nil {
		outputStr := output.String()
		appendError(fmt.Errorf("Failed to build service %s\n%s", service.Name, outputStr))
		terminal.Stdout.Colorf("%s%s| @{r} error: @{|}Failed to build: %s\n", service.Name, spacing, outputStr)
	} else {
		terminal.Stdout.Colorf("%s%s| @{g} (re)built\n", service.Name, spacing)
	}
}
