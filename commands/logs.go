package commands

import (
	"fmt"
	"strings"
	"sync"

	log "github.com/cihub/seelog"
	"github.com/codegangsta/cli"
	"github.com/hpcloud/tail"
	"github.com/mondough/orchestra/services"
	"github.com/wsxiaoys/terminal"
)

var LogsCommand = &cli.Command{
	Name:         "logs",
	Usage:        "Aggregate services logs",
	Action:       BeforeAfterWrapper(LogsAction),
	BashComplete: ServicesBashComplete,
}

var logReceiver chan string

func init() {
	logReceiver = make(chan string)
}

func LogsAction(c *cli.Context) error {
	go ConsumeLogs()
	wg := &sync.WaitGroup{}
	for _, service := range FilterServices(c) {
		wg.Add(1)
		go TailServiceLog(service, wg)
	}
	wg.Wait()
	close(logReceiver)
	return nil
}

func ConsumeLogs() {
	for log := range logReceiver {
		terminal.Stdout.Colorf(log)
	}
}

func TailServiceLog(service *services.Service, wg *sync.WaitGroup) {
	defer wg.Done()
	spacingLength := services.MaxServiceNameLength + 2 - len(service.Name)
	t, err := tail.TailFile(service.LogFilePath, tail.Config{Follow: true})
	if err != nil {
		log.Error(err.Error())
	}
	for line := range t.Lines {
		logReceiver <- fmt.Sprintf("@{%s}%s@{|}%s|  %s\n", service.Color, service.Name, strings.Repeat(" ", spacingLength), line.Text)
	}
}
