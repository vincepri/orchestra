package commands

import (
	"strings"
	"sync"

	"github.com/mondough/orchestra/services"
	"github.com/codegangsta/cli"
	"github.com/wsxiaoys/terminal"
)

var RestartCommand = &cli.Command{
	Name:         "restart",
	Usage:        "Restarts all the services",
	Action:       BeforeAfterWrapper(RestartAction),
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

// RestartAction restarts all the services (or the specified ones)
func RestartAction(c *cli.Context) {
	wg := &sync.WaitGroup{}
	for _, service := range FilterServices(c) {
		wg.Add(1)
		go restart(wg, c, service)
	}
	wg.Wait()
	if c.Bool("attach") || c.Bool("logs") {
		LogsAction(c)
	}
}

func restart(wg *sync.WaitGroup, c *cli.Context, service *services.Service) {
	spacing := strings.Repeat(" ", services.MaxServiceNameLength+2-len(service.Name))

	err := killService(service)
	if err != nil {
		appendError(err)
		terminal.Stdout.Colorf("%s%s| @{r} error: @{|}%s\n", service.Name, spacing, err.Error())
		return
	}

	rebuilt, err := buildAndStart(c, service)
	if err != nil {
		appendError(err)
		terminal.Stdout.Colorf("%s%s| @{r} error: @{|}%s\n", service.Name, spacing, err.Error())
		return
	}

	var rebuiltStatus string
	if rebuilt {
		rebuiltStatus = "rebuilt & "
	}

	terminal.Stdout.Colorf("%s%s| @{c} %srestarted\n", service.Name, spacing, rebuiltStatus)

	wg.Done()
}
