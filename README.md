# Loadbalancer

This is a very simple loadbalancer written in Go

## Usage

The loadbalancer configuration is set by default as following:
```go
configuration = AppConfig{
    InProduction:        false,
    LoadBalancerPort:    ":8080",
    SchedulingAlgorithm: AllowedSchedulingAlgorithms["round-robin"],
    BackendServers:      map[int]string{},
}
```

To use it in Go you can do the following

```go
package main

import "github.com/LarsNieuwenhuizen/loadbalancer"

func main() {
	loadbalancer.Configuration.SetBackendServers(map[int]string{})
	loadbalancer.Start()
}
```

The default configuration sets InProduction to `false`
The setting the backend servers to an empty map the default dummy backend servers will the spinned up so you can test the behaviour
