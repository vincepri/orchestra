package commands

import (
	"os"
	"strings"

	log "github.com/cihub/seelog"
	"github.com/codegangsta/cli"
	"github.com/vinceprignano/orchestra/services"
	"github.com/wsxiaoys/terminal"
)

var StopCommand = &cli.Command{
	Name:   "stop",
	Usage:  "Stops all the services",
	Action: StopAction,
}

func StopAction(c *cli.Context) {
	for _, service := range FilterServices(c) {
		spacing := strings.Repeat(" ", services.MaxServiceNameLength+2-len(service.Name))
		err := killService(service)
		if err != nil {
			log.Error(err)
		} else if service.Process != nil {
			terminal.Stdout.Colorf("%s%s| @{r} stopped\n", service.Name, spacing)
		}
	}
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
