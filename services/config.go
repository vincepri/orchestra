package services

import (
	"io/ioutil"
	"os"

	"github.com/cihub/seelog"
	"gopkg.in/yaml.v2"
)

var OrchestraConfig *Config

type Config struct {
	Environment []string `environment,omitempty`
}

func ParseGlobalConfig() {
	OrchestraConfig = &Config{}
	b, err := ioutil.ReadFile(ProjectPath + "orchestra.yml")
	if err != nil {
		seelog.Criticalf(err.Error())
		os.Exit(1)
	}
	yaml.Unmarshal(b, &OrchestraConfig)
}
