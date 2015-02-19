package services

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"syscall"

	"go/build"

	log "github.com/cihub/seelog"
)

// Init is in charge of initializing an orchestra project
// from the current folder and save relevant information in "~/.orchestra/service"
// making sure that the service directory inside orchestra is available
var OrchestraServicePath string
var ProjectPath string
var Registry map[string]*Service
var MaxServiceNameLength int

type Service struct {
	Name          string
	Description   string
	Path          string
	OrchestraPath string
	LogFilePath   string
	PidFilePath   string
	FileInfo      os.FileInfo
	PackageInfo   *build.Package
	Process       *os.Process
}

func init() {
	Registry = make(map[string]*Service)
}

func Init() {
	ProjectPath, _ = os.Getwd()
	OrchestraServicePath = fmt.Sprintf("%s/.orchestra", ProjectPath)
	if err := os.Mkdir(OrchestraServicePath, 0766); err != nil && os.IsNotExist(err) {
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
				// Check for service.yml and try to import the package
				pkg, err := build.Import(fmt.Sprintf("%s/%s", buildPath, item.Name()), "srcDir", 0)
				if err != nil {
					log.Errorf("Error registering %s", item.Name())
					log.Error(err.Error())
					continue
				}

				service := &Service{
					Name:          item.Name(),
					Description:   "",
					FileInfo:      item,
					PackageInfo:   pkg,
					OrchestraPath: OrchestraServicePath,
					LogFilePath:   fmt.Sprintf("%s/%s.log", OrchestraServicePath, item.Name()),
					PidFilePath:   fmt.Sprintf("%s/%s.pid", OrchestraServicePath, item.Name()),
				}

				// Because I like nice logging
				if len(service.Name) > MaxServiceNameLength {
					MaxServiceNameLength = len(service.Name)
				}

				// Add the service to the registry
				Registry[item.Name()] = service
				if _, err := os.Stat(service.PidFilePath); err == nil {
					bytes, _ := ioutil.ReadFile(service.PidFilePath)
					pid, _ := strconv.Atoi(string(bytes))
					proc, procErr := os.FindProcess(pid)
					if procErr == nil {
						sigError := proc.Signal(syscall.Signal(0))
						if sigError == nil {
							service.Process = proc
						} else {
							os.Remove(service.PidFilePath)
						}
					} else {
						os.Remove(service.PidFilePath)
					}
				}
			}
		}
	}
}
