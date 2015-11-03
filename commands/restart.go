package commands

import (
	"runtime"
	"strings"

	"github.com/codegangsta/cli"
	"github.com/mondough/orchestra/services"
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
	worker := func(service *services.Service) func() {
		return func() {
			restart(c, service)
		}
	}

	pool := make(workerPool, runtime.NumCPU())
	svcs := services.Sort(FilterServices(c))
	for _, service := range svcs {
		pool.Do(worker(service))
	}
	pool.Drain()

	if c.Bool("attach") || c.Bool("logs") {
		LogsAction(c)
	}
}

func restart(c *cli.Context, service *services.Service) {
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
}
