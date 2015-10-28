package services

import (
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strconv"
	"strings"
	"syscall"

	"go/build"

	"gopkg.in/yaml.v1"

	log "github.com/cihub/seelog"
)

var (
	// Internal Service Registry
	Registry map[string]*Service

	// Path variables
	OrchestraServicePath string
	ProjectPath          string

	// Other internal variables
	MaxServiceNameLength int
	colors               = []string{"g", "b", "c", "m", "y", "w"}
)

func init() {
	Registry = make(map[string]*Service)
}

// Init initializes the OrchestraServicePath to the workingdir/.orchestra path
// and starts the service discovery
func Init() {
	DiscoverServices()
}

func Sort(r map[string]*Service) SortableRegistry {
	sr := make(SortableRegistry, 0)
	for _, v := range r {
		sr = append(sr, v)
	}
	sort.Sort(sr)
	return sr
}

type SortableRegistry []*Service

func (s SortableRegistry) Len() int {
	return len(s)
}

func (s SortableRegistry) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s SortableRegistry) Less(i, j int) bool {
	return s[i].Name < s[j].Name
}

// Service encapsulates all the information needed for a service
type Service struct {
	Name        string
	Description string
	Path        string
	Color       string

	// Path
	OrchestraPath string
	LogFilePath   string
	PidFilePath   string
	BinPath       string

	// Process, Service and Package information
	FileInfo    os.FileInfo
	PackageInfo *build.Package
	Process     *os.Process
	Env         []string
}

func (s *Service) IsRunning() bool {
	if _, err := os.Stat(s.PidFilePath); err == nil {
		bytes, _ := ioutil.ReadFile(s.PidFilePath)
		pid, _ := strconv.Atoi(string(bytes))
		proc, procErr := os.FindProcess(pid)
		if procErr == nil {
			sigError := proc.Signal(syscall.Signal(0))
			if sigError == nil {
				s.Process = proc
				return true
			} else {
				os.Remove(s.PidFilePath)
			}
		}
	} else {
		os.Remove(s.PidFilePath)
	}
	return false
}

// DiscoverServices walks into the project path and looks in every subdirectory
// for the service.yml file. For every service it registers it after trying
// to import the package using Go's build.Import package
func DiscoverServices() {
	gopath := strings.TrimRight(os.Getenv("GOPATH"), "/")
	buildPath := strings.Replace(ProjectPath, gopath+"/src/", "", 1)
	fd, _ := ioutil.ReadDir(ProjectPath)
	for _, item := range fd {
		serviceName := item.Name()
		if item.IsDir() && !strings.HasPrefix(serviceName, ".") {
			serviceConfigPath := fmt.Sprintf("%s/%s/service.yml", ProjectPath, serviceName)
			if _, err := os.Stat(serviceConfigPath); err == nil {
				// Check for service.yml and try to import the package
				pkg, err := build.Import(fmt.Sprintf("%s/%s", buildPath, serviceName), "srcDir", 0)
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
					LogFilePath:   fmt.Sprintf("%s/%s.log", OrchestraServicePath, serviceName),
					PidFilePath:   fmt.Sprintf("%s/%s.pid", OrchestraServicePath, serviceName),
					Color:         colors[len(Registry)%len(colors)],
					Path:          fmt.Sprintf("%s/%s", ProjectPath, serviceName),
				}

				// Parse env variable in configuration
				var serviceConfig struct {
					Env map[string]string `env,omitempty`
				}
				b, err := ioutil.ReadFile(serviceConfigPath)
				if err != nil {
					log.Criticalf(err.Error())
					os.Exit(1)
				}
				yaml.Unmarshal(b, &serviceConfig)
				for k, v := range serviceConfig.Env {
					service.Env = append(service.Env, fmt.Sprintf("%s=%s", k, v))
				}

				// Because I like nice logging
				if len(serviceName) > MaxServiceNameLength {
					MaxServiceNameLength = len(serviceName)
				}

				if binPath := os.Getenv("GOBIN"); binPath != "" {
					service.BinPath = fmt.Sprintf("%s/%s", binPath, serviceName)
				} else {
					service.BinPath = fmt.Sprintf("%s/bin/%s", os.Getenv("GOPATH"), serviceName)
				}

				// Add the service to the registry
				Registry[serviceName] = service
				// When registering, we take care, on every run, to check
				// if the process is still alive.
				service.IsRunning()
			}
		}
	}
}
