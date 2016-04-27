package commands

import (
	"bytes"
	"fmt"
	"os/exec"
	"runtime"
	"strings"

	"github.com/codegangsta/cli"
	"github.com/mondough/orchestra/services"
	"github.com/wsxiaoys/terminal"
)

var InstallCommand = &cli.Command{
	Name:         "install",
	Usage:        "Installs all the services",
	Action:       BeforeAfterWrapper(InstallAction),
	BashComplete: ServicesBashComplete}

// InstallAction installs all the services (or the specified ones)
func InstallAction(c *cli.Context) {
	worker := func(service *services.Service) func() {
		return func() {
			spacing := strings.Repeat(" ", services.MaxServiceNameLength+2-len(service.Name))
			rebuilt, err := installService(service)
			if err != nil {
				appendError(err)
				terminal.Stdout.Colorf("%s%s| @{r} error: @{|}%v\n", service.Name, spacing, err)
			} else if rebuilt {
				terminal.Stdout.Colorf("%s%s| @{g} installed\n", service.Name, spacing)
			} else {
				terminal.Stdout.Colorf("%s%s| @{g} already up to date\n", service.Name, spacing)
			}
		}
	}

	pool := make(workerPool, runtime.NumCPU())
	for _, service := range FilterServices(c) {
		pool.Do(worker(service))
	}
	pool.Drain()
}

// installService runs go install in the service directory
func installService(service *services.Service) (bool, error) {
	cmd := exec.Command("nice", "-n", niceness, "go", "install", "-v")
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
