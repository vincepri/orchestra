package commands

import (
	"strings"

	"github.com/codegangsta/cli"
	"github.com/vinceprignano/orchestra/services"
	"github.com/wsxiaoys/terminal"
)

var RestartCommand = &cli.Command{
	Name:         "restart",
	Usage:        "Restarts all the services",
	Action:       BeforeAfterWrapper(RestartAction),
	BashComplete: ServicesBashComplete,
}

// RestartAction restarts all the services (or the specified ones)
func RestartAction(c *cli.Context) {
	for _, service := range FilterServices(c) {
		spacing := strings.Repeat(" ", services.MaxServiceNameLength+2-len(service.Name))

		err := killService(service)
		if err != nil {
			appendError(err)
			terminal.Stdout.Colorf("%s%s| @{r} error: @{|}%s\n", service.Name, spacing, err.Error())
			continue
		}
		err = startService(c, service)
		if err != nil {
			appendError(err)
			terminal.Stdout.Colorf("%s%s| @{r} error: @{|}%s\n", service.Name, spacing, err.Error())
			continue
		}

		terminal.Stdout.Colorf("%s%s| @{c} restarted\n", service.Name, spacing)
	}
}
