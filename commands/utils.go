package commands

import (
	"fmt"
	"math"
	"os"
	"strings"

	log "github.com/cihub/seelog"
	"github.com/codegangsta/cli"
	"github.com/vinceprignano/orchestra/config"
	"github.com/vinceprignano/orchestra/services"
)

// This is temporary, very very alpha and may change soon
func FilterServices(c *cli.Context) map[string]*services.Service {
	excludeMode := 0
	args := c.Args()
	for _, s := range args {
		serv := s
		if strings.HasPrefix(s, "~") {
			serv = strings.Replace(s, "~", "", 1)
		}
		if _, ok := services.Registry[serv]; ok {
			if strings.HasPrefix(s, "~") {
				excludeMode += 1
				delete(services.Registry, serv)
			} else {
				excludeMode -= 1
			}
		} else {
			log.Errorf("Service %s not found", s)
			return nil
		}
	}
	if math.Abs(float64(excludeMode)) != float64(len(args)) {
		log.Critical("You can't exclude and include services at the same time")
		os.Exit(1)
	}
	if excludeMode < 0 {
		for name := range services.Registry {
			included := false
			for _, s := range args {
				if name == s {
					included = true
					break
				}
			}
			if !included {
				delete(services.Registry, name)
			}
		}
	}
	return services.Registry
}

func ServicesBashComplete(c *cli.Context) {
	for name := range services.Registry {
		fmt.Println(name)
		fmt.Println("~" + name)
	}
}

func BeforeAfterWrapper(cmdName string, f func(c *cli.Context)) func(c *cli.Context) {
	return func(c *cli.Context) {
		config.GetBeforeFunc(cmdName)(c)
		f(c)
		config.GetAfterFunc(cmdName)(c)
	}
}
