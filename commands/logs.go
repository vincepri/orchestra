package commands

import (
	"fmt"
	"strings"

	"github.com/ActiveState/tail"
	log "github.com/cihub/seelog"
	"github.com/codegangsta/cli"
	"github.com/vinceprignano/orchestra/services"
)

var LogsCommand = &cli.Command{
	Name:   "logs",
	Usage:  "Aggregate services logs",
	Action: LogsAction,
}

var logReceiver chan string
var maxServiceNameLength int

func init() {
	logReceiver = make(chan string)
}

func LogsAction(c *cli.Context) {
	done := make(chan bool)
	for name := range services.Registry {
		if len(name) > maxServiceNameLength {
			maxServiceNameLength = len(name)
		}
	}
	go ConsumeLogs(done)
	for _, service := range services.Registry {
		go TailServiceLog(service)
	}
	<-done
}

func ConsumeLogs(done chan bool) {
	for log := range logReceiver {
		fmt.Println(log)
	}
	done <- true
}

func TailServiceLog(service *services.Service) {
	spacingLength := maxServiceNameLength + 2 - len(service.Name)
	t, err := tail.TailFile(service.LogFilePath, tail.Config{Follow: true})
	if err != nil {
		log.Error(err.Error())
	}
	for line := range t.Lines {
		logReceiver <- fmt.Sprintf("%s%s|  %s", service.Name, strings.Repeat(" ", spacingLength), line.Text)
	}
}
