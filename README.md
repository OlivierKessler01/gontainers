<p align="center">
  <img src="https://raw.githubusercontent.com/olivierkessler01/gontainers/main/doc/images/logo.svg" alt="Gontainers Logo" width="200"/>
</p>

<p align="center">A Linux container runtime with Kubernetes compatibility.</p>

---


## Compatibility 

At the moment, the `Gontainer` runtime is only compatible with Linux x86.

Through `CRI-API` the goal of the project is to development fully k8s-compatible 
container runtime. This entails having a long-lived process listening for gRPC connections.

## Usage

```bash
$ ./gontainers

NAME:
   gontainers - A new cli application

USAGE:
   gontainers [global options] [command [command options]]

COMMANDS:
   run      Run a container, get a PID.
   list     List containers.
   remove   Remove a container.
   server   Server the CR-API gRPC server.
   init     Init the container database. Run this before using gontainer.
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h  show help
```

## Development

Show help
```bash
make
```

## Setup dev environment

* Copy `config.default.yaml` to a new `config.yaml` file and fill it with your config.

* Run `./gontainers init`, to initialize the containers tracking database.

* Run `go mod vendor`

* Run `make build`, to compile gontainers.

Then you can use gontainers
