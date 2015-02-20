package config

import (
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"github.com/cihub/seelog"
	"github.com/codegangsta/cli"
	"github.com/vinceprignano/orchestra/services"
	"gopkg.in/yaml.v2"
)

var orchestra *Config

type Config struct {
	Environment []string            `environment,omitempty`
	Before      map[string][]string `before,omitempty`
	After       map[string][]string `after,omitempty`
}

func ParseGlobalConfig() {
	orchestra = &Config{}
	b, err := ioutil.ReadFile(services.ProjectPath + "orchestra.yml")
	if err != nil {
		seelog.Criticalf(err.Error())
		os.Exit(1)
	}
	yaml.Unmarshal(b, &orchestra)
}

// GetEnvironment returns all the environment variables for a given service
// including the ones specified in the global config
func GetEnvironmentVars(service *services.Service) []string {
	return orchestra.Environment
}

func runCommands(cmds []string) {
	for _, command := range cmds {
		cmdLine := strings.Split(command, " ")
		cmd := exec.Command(cmdLine[0], cmdLine[1:]...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Start()
		if err != nil {
			seelog.Error(err.Error())
		}
	}
}

func GetBeforeFunc(cmdName string) func(c *cli.Context) {
	return func(c *cli.Context) {
		runCommands(orchestra.Before[cmdName])
	}
}

func GetAfterFunc(cmdName string) func(c *cli.Context) {
	return func(c *cli.Context) {
		runCommands(orchestra.After[cmdName])
	}
}
