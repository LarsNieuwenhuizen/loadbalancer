package loadbalancer

import (
	"errors"
	"log"
)

var (
	AllowedSchedulingAlgorithms = map[string]string{
		"round-robin":       "round-robin",
		"least-connections": "least-connections",
	}
	configuration = AppConfig{
		InProduction:        false,
		LoadBalancerPort:    ":8080",
		SchedulingAlgorithm: AllowedSchedulingAlgorithms["round-robin"],
		BackendServers:      map[int]string{},
	}
)

// AppConfig holds the configuration for the application.
type AppConfig struct {
	LoadBalancerPort    string
	BackendServers      map[int]string
	SchedulingAlgorithm string
	InProduction        bool
}

// UseRoundRobinSchedulingAlgorithm sets the scheduling algorithm to round robin.
func (a *AppConfig) UseRoundRobinSchedulingAlgorithm() {
	setSchedulingAlgorithm(AllowedSchedulingAlgorithms["round-robin"], a)
}

// UseLeastConnectionsSchedulingAlgorithm sets the scheduling algorithm to least connections.
func (a *AppConfig) UseLeastConnectionsSchedulingAlgorithm() {
	setSchedulingAlgorithm(AllowedSchedulingAlgorithms["least-connections"], a)
}

// setSchedulingAlgorithm sets the scheduling algorithm for the load balancer.
func setSchedulingAlgorithm(algorithm string, a *AppConfig) {
	if _, ok := AllowedSchedulingAlgorithms[algorithm]; ok {
		a.SchedulingAlgorithm = algorithm
		return
	}

	log.Fatal(errors.New("invalid scheduling algorithm"))
}

// Setup sets the configuration for the application.
func Setup() *AppConfig {
	return &configuration
}

/**
 * SetBackendServers sets the backend servers for the load balancer.
 * You can pass a map of server indexes to server addresses.
 * For example, map[0] = "http://localhost:8081", map[1] = "http://localhost:8082".
 * If you are not in production and you pass an empty map, the development dummy backend servers will be used.
 */
func (a *AppConfig) SetBackendServers(servers map[int]string) {
	if !a.InProduction && len(servers) == 0 {
		log.Println("Creating dummy backend servers")
		developmentBackendServers := map[int]string{
			0: "http://localhost:8081",
			1: "http://localhost:8082",
			2: "http://localhost:8083",
		}
		a.BackendServers = developmentBackendServers
		return
	}
	a.BackendServers = servers
}
