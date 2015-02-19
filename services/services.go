package services

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"go/build"

	log "github.com/cihub/seelog"
)

// Init is in charge of initializing an orchestra project
// from the current folder and save relevant information in "~/.orchestra/service"
// making sure that the service directory inside orchestra is available
var ServicePath string
var ProjectPath string
var ServiceRegistry map[string]*Service

type Service struct {
	Name        string
	Description string
	FileInfo    os.FileInfo
	PackageInfo *build.Package
}

func init() {
	ServiceRegistry = make(map[string]*Service)
}

func Init() {
	ProjectPath, _ = os.Getwd()
	dirPath := strings.Split(ProjectPath, "/")
	ServicePath = fmt.Sprintf("%s/.orchestra/%s", ProjectPath, dirPath[len(dirPath)-1])
	if err := os.Mkdir(ServicePath, 0766); err != nil && os.IsNotExist(err) {
		log.Critical(err.Error())
		os.Exit(1)
	}
	DiscoverServices()
}

func DiscoverServices() {
	buildPath := strings.Replace(ProjectPath, os.Getenv("GOPATH")+"/src/", "", 1)
	fd, _ := ioutil.ReadDir(ProjectPath)
	for _, item := range fd {
		if item.IsDir() && !strings.HasPrefix(item.Name(), ".") {
			if _, err := os.Stat(fmt.Sprintf("%s/%s/service.yml", ProjectPath, item.Name())); err == nil {
				log.Infof("Found service.yml in %s ", item.Name())
				pkg, err := build.Import(fmt.Sprintf("%s/%s", buildPath, item.Name()), "srcDir", 0)
				if err != nil {
					log.Errorf("Error registering %s", item.Name())
					log.Error(err.Error())
					continue
				}
				ServiceRegistry[item.Name()] = &Service{
					Name:        item.Name(),
					Description: "",
					FileInfo:    item,
					PackageInfo: pkg,
				}
			}
		}
	}
}
