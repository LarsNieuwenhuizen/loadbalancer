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

Download the loadbalancer and run it

`./lb --config config.yaml`

- The loadbalancer will now start on port 8080
- It will fire up the webservers as dummy backend servers because we set `startGivenServers` to true
- It will use the round-robin mode to choose the next server to go to for each request
