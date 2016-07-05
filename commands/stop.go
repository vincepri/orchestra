package commands

import (
	"os"
	"strings"

	"github.com/codegangsta/cli"
	"github.com/mondough/orchestra/services"
	"github.com/wsxiaoys/terminal"
)

var StopCommand = &cli.Command{
	Name:         "stop",
	Usage:        "Stops all the services",
	Action:       BeforeAfterWrapper(StopAction),
	BashComplete: ServicesBashComplete,
}

// StopAction stops all the services (or the specified ones)
func StopAction(c *cli.Context) error {
	svcs := services.Sort(FilterServices(c))
	for _, service := range svcs {
		spacing := strings.Repeat(" ", services.MaxServiceNameLength+2-len(service.Name))
		err := killService(service)
		if err != nil {
			appendError(err)
			terminal.Stdout.Colorf("%s%s| @{r} error: @{|}%s\n", service.Name, spacing, err.Error())
		} else if service.Process != nil {
			terminal.Stdout.Colorf("%s%s| @{r} stopped\n", service.Name, spacing)
		}
	}
	return nil
}

func killService(service *services.Service) error {
	if service.Process != nil {
		err := service.Process.Kill()
		defer os.Remove(service.PidFilePath)
		if err != nil {
			return err
		}
	}
	return nil
}
