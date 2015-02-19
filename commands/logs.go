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

func init() {
	logReceiver = make(chan string)
}

func LogsAction(c *cli.Context) {
	done := make(chan bool)
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
	spacingLength := services.MaxServiceNameLength + 2 - len(service.Name)
	t, err := tail.TailFile(service.LogFilePath, tail.Config{Follow: true})
	if err != nil {
		log.Error(err.Error())
	}
	for line := range t.Lines {
		logReceiver <- fmt.Sprintf("%s%s|  %s", service.Name, strings.Repeat(" ", spacingLength), line.Text)
	}
}
