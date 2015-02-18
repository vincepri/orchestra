# Orchestra
Orchestra is a toolkit to manage a fleet of Go binaries/services. 

## Use Case
In a service oriented architecture, an existing challenge is to manage running binaries, dependencies and networking. **Orchestra** goal is to provide a unique place where you can *run*, *stop*, *aggregate logs* and *config* your Go binaries.

## Goals
- **Standalone**: The first requirement is to provide support as a standalone tool, without external dependencies (apart from *docker* when using vendors). An optional flag will be provided to run the Go binaries inside a docker container.
- **Rerun**: Reload (build, test and run) services upon modification
- **Vendors**: Specify dependecies with existing services. Vendor software (e.g. postgres, rabbitmq, etc) will run inside Docker. This feature relies on `crane`.
- **Configuration**: A global configuration file will be required to specify `ENV` variables for every service. An optional configuration file can be specified 
- **Testing**: Run unit tests and acceptance tests.
- **Logging**: Show aggregated logs for every service in the fleet.
- **Reliability**: Services started with Orchestra should operate atomically and in a reliable way. When Orchestra starts a service should check for the running processes, and match and kill services running outside this toolkit (e.g. running `go run main.go` inside a service folder or the `./service` binary)
- **Scale**: Services should be able to scale, by default the scale level is set to 1. A scale value can be set for every service inside their configuration file.

## Command Line Cmd/Flags Proposal
- `start`
    - *noreload*: doesn't use rerun
    - *scale*: overwrites the defined scale parameter
- `stop`
- `restart`
- `vendors`

## Design Proposal
- The application will create a `.orchestra` folder in the directory (gopath) that we want to orchestrate. The purpose of this directoy is to provide the necessary context when running orchestra again.
- The application can be run in every folder inside the `GOPATH`, e.g. let's assume that we have a single repo (called `myservices`) with all the services in different folders, you can run orchestra inside that folder and it will recursively (with `depth=1`) look in every directory for a configuration file (i.e. `service.yml`) and register the services. This approach works seamlessly with services split in multiple repositories in one GitHub account.
- Logging is extremely important, and with great logging comes great responsibilities. For simplicity we'll log every subprocess spawn into files, then use `github.com/ActiveState/tail` library to aggregate the logs of different services, maybe make them colorful and prefix them with the service name.