package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"reflect"
	"strings"

	"github.com/cihub/seelog"
	"github.com/codegangsta/cli"
	"gopkg.in/yaml.v2"
)

var orchestra *Config
var ConfigPath string
var globalEnvs []string

type ContextConfig struct {
	Env    map[string]string `env,omitempty`
	Before []string          `before,omitempty`
	After  []string          `after,omitempty`
}

type Config struct {
	// Global Configuration
	Env    map[string]string `env,omitempty`
	Before []string          `before,omitempty`
	After  []string          `after,omitempty`
	GoRun  bool              `gorun,omitempty`

	// Configuration for Commands
	Build   ContextConfig `build,omitempty`
	Export  ContextConfig `export,omitempty`
	Install ContextConfig `install,omitempty`
	Logs    ContextConfig `logs,omitempty`
	Ps      ContextConfig `ps,omitempty`
	Restart ContextConfig `restart,omitempty`
	Start   ContextConfig `start,omitempty`
	Stop    ContextConfig `stop,omitempty`
	Test    ContextConfig `test,omitempty`
}

func GetBaseEnvVars() map[string]string {
	return orchestra.Env
}

func UseGoRun() bool {
	return orchestra.GoRun
}

func ParseGlobalConfig() {
	orchestra = &Config{}
	b, err := ioutil.ReadFile(ConfigPath)
	if err != nil {
		seelog.Criticalf(err.Error())
		os.Exit(1)
	}
	yaml.Unmarshal(b, &orchestra)

	globalEnvs = os.Environ()
	for k, v := range orchestra.Env {
		globalEnvs = append([]string{fmt.Sprintf("%s=%s", k, v)}, globalEnvs...)
	}
}

func GetEnvForCommand(c *cli.Context) []string {
	envs := globalEnvs
	for k, v := range getConfigFieldByName(c.Command.Name).Env {
		envs = append(envs, fmt.Sprintf("%s=%s", k, v))
	}
	return envs
}

func runCommands(c *cli.Context, cmds []string) error {
	for _, command := range cmds {
		cmdLine := strings.Split(command, " ")
		cmd := exec.Command(cmdLine[0], cmdLine[1:]...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Env = GetEnvForCommand(c)
		err := cmd.Start()
		if err != nil {
			return err
		}
		cmd.Wait()
		if !cmd.ProcessState.Success() {
			return fmt.Errorf("Command %s exited with error", cmdLine[0])
		}
	}
	return nil
}

func GetBeforeFunc() func(c *cli.Context) error {
	return func(c *cli.Context) error {
		err := runCommands(c, orchestra.Before)
		if err != nil {
			return err
		}
		err = runCommands(c, getConfigFieldByName(c.Command.Name).Before)
		if err != nil {
			return err
		}
		return nil
	}
}

func GetAfterFunc() func(c *cli.Context) error {
	return func(c *cli.Context) error {
		err := runCommands(c, orchestra.After)
		if err != nil {
			return err
		}
		err = runCommands(c, getConfigFieldByName(c.Command.Name).After)
		if err != nil {
			return err
		}
		return nil
	}
}

func getConfigFieldByName(name string) ContextConfig {
	initial := strings.Split(name, "")[0]
	value := reflect.ValueOf(orchestra)
	f := reflect.Indirect(value).FieldByName(strings.Replace(name, initial, strings.ToUpper(initial), 1))
	return f.Interface().(ContextConfig)
}
