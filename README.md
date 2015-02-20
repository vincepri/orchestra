# Orchestra [![wercker status](https://app.wercker.com/status/16ba07e3d295feb5c3874207a9f3fe36/s "wercker status")](https://app.wercker.com/project/bykey/16ba07e3d295feb5c3874207a9f3fe36) [![GoDoc](https://godoc.org/github.com/vinceprignano/orchestra?status.svg)](https://godoc.org/github.com/vinceprignano/orchestra)
Orchestra is a toolkit to manage a fleet of Go binaries/services. A unique place where you can run, stop, aggregate logs and config your Go binaries.

![](https://cloud.githubusercontent.com/assets/3118335/6255612/4811c940-b7a9-11e4-8d06-966981de3926.png)

> You can find an application design/proposal document [here](https://github.com/vinceprignano/orchestra/blob/master/DESIGN.md)

## Build & Install
`go get -u github.com/vinceprignano/orchestra`

## Start an Orchestra Project
You should have an `orchestra.yml` file in your root directory and a `service.yml` file in every service directory.

```
.
├── first-service
│   ├── main.go
│   └── service.yml			<- Service file
├── second-service
│   ├── second.go
│   ├── main.go
│   └── service.yml			<- Service file
└── orchestra.yml           <- Main project file
```

You can specify a custom configuration file using the `--config` flag or setting the `ORCHESTRA_CONFIG` env variable.

### Example
```yaml
env:
	- "ABC=somethingGlobal"
before:
	- "echo I am a global command before"
after:
	- "echo I am a global after"
```

## Commands
- **start** `[<service>...]` Starts every service
- **stop** `[<service>...]` Stops every service
- **restart** `[<service>...]` Restarts every service
- **logs** `[<service>...]` Aggregates the output from the services
- **test** `--verbose [<service>...]` Runs `go test ./...` for every service
- **ps** Displays the status of every service

A service name can be prefixed with `~` to run a command in exclusion mode.
For example `orchestra start ~second-service` will start everything expect the second-service.

### Configuring commands
Every command can be configured separately with special environment variables or with before/after commands.

For example, in `orchestra.yml` you can configure to `echo BEFORE START` before running `orchestra start` command.

```yaml
env:
	- "ABC=somethingGlobal"
before:
	- "echo I am a global command before"
after:
	- "echo I am a global after"
start:
	env:
    	- "ABC=somethingStart"
    after:
    	- "echo BEFORE START"
```

## Autocomplete
Orchestra supports bash autocomplete.
```sh
source $GOPATH/src/github.com/vinceprignano/orchestra/autocomplete/orchestra
```

## TODO
- [ ] Add support for rerun
- [ ] Scale option
- [ ] Give build information
- [ ] Pick up a license
- [ ] Add support for `start --attach` flag, to start logging right after starting a set of services
- [ ] Tests, a lot