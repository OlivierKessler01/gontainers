<p align="center">
  <img src="https://raw.githubusercontent.com/olivierkessler01/gontainers/main/doc/images/logo.svg" alt="Gontainers Logo" width="200"/>
</p>

<h1 align="center">Gontainers</h1>
<p align="center">A Linux-native container runtime with Kubernetes compatibility.</p>

---


## Compatibility 

At the moment, the `Gontainer` runtime is only compatible with Linux x86.

Through `CRI-API` the goal of the project is to development fully k8s-compatible 
container runtime. This entails having a long-lived process listening for gRPC connections.


## Development

Show help
```bash
make
```

## Usage 

* Copy `config.default.yaml` to a new `config.yaml` file and fill it with your config.

* Run `./gontainers init`, to initialize the containers tracking database.

* Run `go mod vendor`

* Run `make build`, to compile gontainers.

Then you can use gontainers : 
```bash
#Manual container management
./gontainers help
sudo ./gontainers run '<>' #Needs root-level access to create namespaces for the container
./gontainers list
sudo ./gontainers remove '<PID>'

#K8s container management
./gontainers serve #Start the gRPC server
```
