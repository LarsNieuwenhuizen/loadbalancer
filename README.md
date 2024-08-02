# Loadbalancer

This is a very simple loadbalancer written in Go

## Usage

The loadbalancer simply needs a yaml config file like this:

```yaml
loadbalancer:
  inProduction: false
  backendServers:
    - "http://localhost:8082"
    - "http://localhost:8083"
    - "http://localhost:8090"
  startGivenServers: true
  port: 8080
  schedulingAlgorithm: "round-robin"
```

## Download

Download the latest version of loadbalancer *[here](https://github.com/LarsNieuwenhuizen/loadbalancer/releases/latest)*

It's currently being built for Linux and MacOS

Download the one suited for your system, extract the tarball and run it.

## Example with the above stated loadbalancer.yaml

`./loadbalancer start --config ~/loadbalancer.yaml`

- The loadbalancer will now start on port 8080
- It will fire up the webservers as dummy backend servers because we set `startGivenServers` to true
- It will use the round-robin mode to choose the next server to go to for each request
- Hit http://localhost:8080/ in your browser and you should see a a response from the next backend server on every refresh

## Loadbalancer algorithms

Currently the following options for loadbalancing are:

- round-robin

*more to come soon...*
