package loadbalancer

import (
	"errors"
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

var (
	AllowedSchedulingAlgorithms = map[string]string{
		"round-robin":       "round-robin",
		"least-connections": "least-connections",
	}
	configuration = AppConfig{
		InProduction:        false,
		StartGivenServers:   false,
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
	StartGivenServers   bool
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

// LoadFromYaml loads the configuration from a given yaml file.
func (a *AppConfig) LoadFromYaml(filepath string) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		log.Fatal(err)
	}

	type LoadBalancerData struct {
		Port                string   `yaml:"port"`
		SchedulingAlgorithm string   `yaml:"schedulingAlgorithm"`
		StartGivenServers   bool     `yaml:"startGivenServers"`
		BackendServers      []string `yaml:"backendServers"`
		InProduction        bool     `yaml:"inProduction"`
	}

	type Data struct {
		LoadBalancer LoadBalancerData `yaml:"loadbalancer"`
	}

	var structData Data

	err = yaml.Unmarshal(data, &structData)
	if err != nil {
		log.Fatal(err)
	}

	if len(structData.LoadBalancer.BackendServers) > 0 {
		for i, server := range structData.LoadBalancer.BackendServers {
			a.BackendServers[i] = server
		}
	}

	a.InProduction = structData.LoadBalancer.InProduction
	a.LoadBalancerPort = ":" + structData.LoadBalancer.Port
	a.SchedulingAlgorithm = structData.LoadBalancer.SchedulingAlgorithm
	a.StartGivenServers = structData.LoadBalancer.StartGivenServers
}
