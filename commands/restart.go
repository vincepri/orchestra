package commands

import (
	"strings"

	log "github.com/cihub/seelog"
	"github.com/codegangsta/cli"
	"github.com/vinceprignano/orchestra/services"
	"github.com/wsxiaoys/terminal"
)

var RestartCommand = &cli.Command{
	Name:   "restart",
	Usage:  "Restarts all the services",
	Action: RestartAction,
}

func RestartAction(c *cli.Context) {
	for _, service := range services.Registry {
		spacing := strings.Repeat(" ", services.MaxServiceNameLength+2-len(service.Name))

		err := killService(service)
		if err != nil {
			log.Error(err.Error())
			continue
		}
		err = startService(service)
		if err != nil {
			log.Error(err.Error())
			continue
		}
		terminal.Stdout.Colorf("%s%s| @{c} restarted\n", service.Name, spacing)
	}
}
