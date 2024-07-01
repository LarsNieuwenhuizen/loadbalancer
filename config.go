package loadbalancer

import (
	"errors"
	"os"

	"gopkg.in/yaml.v2"
)

var (
	AllowedSchedulingAlgorithms = map[string]string{
		"round-robin": "round-robin",
		// "least-connections": "least-connections", // Not implemented yet
	}
	ErrCannotLoadYamlFile              = errors.New("cannot load yaml file")
	ErrNoBackendServersConfigured      = errors.New("no backend servers configured to load balance")
	ErrNoPortConfigured                = errors.New("port is required in the configuration file")
	ErrNoSchedulingAlgorithmConfigured = errors.New("scheduling algorithm is required in the configuration file")
	ErrUnallowedSchedulingAlgorithm    = errors.New("unallowed scheduling algorithm")
)

// AppConfig holds the configuration for the application.
type AppConfig struct {
	LoadBalancerPort    string
	BackendServers      map[int]BackendServer
	SchedulingAlgorithm string
	InProduction        bool
	StartGivenServers   bool
}

// setSchedulingAlgorithm sets the scheduling algorithm for the load balancer.
func (lb *LoadBalancer) SetSchedulingAlgorithm(algorithm string) error {
	if _, ok := AllowedSchedulingAlgorithms[algorithm]; ok {
		lb.Configuration.SchedulingAlgorithm = algorithm
		return nil
	}

	return ErrUnallowedSchedulingAlgorithm
}

// LoadFromYaml loads the configuration from a given yaml file.
func (lb *LoadBalancer) ConfigureFromYaml(filepath string) error {
	lb.Configuration = &AppConfig{
		BackendServers: map[int]BackendServer{},
	}

	data, err := os.ReadFile(filepath)
	if err != nil {
		return ErrCannotLoadYamlFile
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
		return ErrCannotLoadYamlFile
	}

	if structData.LoadBalancer.Port == "" {
		return ErrNoPortConfigured
	}

	if structData.LoadBalancer.SchedulingAlgorithm == "" {
		return ErrNoSchedulingAlgorithmConfigured
	}

	if len(structData.LoadBalancer.BackendServers) == 0 {
		return ErrNoBackendServersConfigured
	}

	for i, server := range structData.LoadBalancer.BackendServers {
		lb.Configuration.BackendServers[i] = BackendServer{
			Address:           server,
			ActiveConnections: 0,
		}
	}

	lb.NextServer = lb.Configuration.BackendServers[0]
	lb.Configuration.InProduction = structData.LoadBalancer.InProduction
	lb.Configuration.LoadBalancerPort = ":" + structData.LoadBalancer.Port
	lb.Configuration.StartGivenServers = structData.LoadBalancer.StartGivenServers

	err = lb.SetSchedulingAlgorithm(structData.LoadBalancer.SchedulingAlgorithm)
	if err != nil {
		return err
	}

	return nil
}
