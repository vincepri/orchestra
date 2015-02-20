package config

import (
	"io/ioutil"
	"os"
	"os/exec"
	"reflect"
	"strings"

	"github.com/cihub/seelog"
	"github.com/codegangsta/cli"
	"github.com/vinceprignano/orchestra/services"
	"gopkg.in/yaml.v2"
)

var orchestra *Config

type ContextConfig struct {
	Env    []string `env,omitempty`
	Before []string `before,omitempty`
	After  []string `after,omitempty`
}

type Config struct {
	// Global Configuration
	Env    []string `env,omitempty`
	Before []string `before,omitempty`
	After  []string `after,omitempty`

	// Configuration for Commands
	Start   ContextConfig `start,omitempty`
	Stop    ContextConfig `stop,omitempty`
	Restart ContextConfig `restart,omitempty`
	Ps      ContextConfig `ps,omitempty`
	Logs    ContextConfig `logs,omitempty`
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
func GetEnvironmentVars(c *cli.Context, service *services.Service) []string {
	return append(orchestra.Env, getConfigFieldByName(c.Command.Name).Env...)
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
		cmd.Wait()
	}
}

func GetBeforeFunc() func(c *cli.Context) {
	return func(c *cli.Context) {
		runCommands(orchestra.Before)
		runCommands(getConfigFieldByName(c.Command.Name).Before)
	}
}

func GetAfterFunc() func(c *cli.Context) {
	return func(c *cli.Context) {
		runCommands(orchestra.After)
		runCommands(getConfigFieldByName(c.Command.Name).After)
	}
}

func getConfigFieldByName(name string) ContextConfig {
	initial := strings.Split(name, "")[0]
	value := reflect.ValueOf(orchestra)
	f := reflect.Indirect(value).FieldByName(strings.Replace(name, initial, strings.ToUpper(initial), 1))
	return f.Interface().(ContextConfig)
}
