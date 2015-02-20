package config

import (
	"io/ioutil"
	"os"

	"github.com/cihub/seelog"
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
