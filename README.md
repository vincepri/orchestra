# Orchestra
Orchestra is a toolkit to manage a fleet of Go binaries/services. 

## Use Case
In a service oriented architecture, an existing challenge is to manage running binaries, dependencies and networking. **Orchestra** goal is to provide a unique place where you can *run*, *stop*, *aggregate logs* and *config* your Go binaries.

## Goals
- **Standalone**: The first requirement is to provide support as a standalone tool, without external dependencies (apart from *docker* when using vendors). An optional flag will be provided to run the Go binaries inside a docker container.
- **Rerun**: Reload (build, test and run) services upon modification
- **Vendors**: Specify dependecies with existing services, vendor software (e.g. postgres, rabbitmq, etc) will run inside Docker.
- **Configuration**: A global configuration file will be required to specify `ENV` variables for every service. An optional configuration file can be specified 
- **Testing**: Run unit tests and acceptance tests.
- **Logging**: Show aggregated logs for every service in the fleet.
