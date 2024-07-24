package loadbalancer

import (
	"net/http"
	"sync"
)

type LoadBalancer struct {
	Configuration   *AppConfig
	NextServerIndex int
	NextServer      BackendServer
	mutex           sync.Mutex
}

func (lb *LoadBalancer) Start() error {
	lb.mutex.Lock()
	defer lb.mutex.Unlock()

	if !lb.Configuration.InProduction && lb.Configuration.StartGivenServers {
		lb.startBackendServers()
	}

	http.ListenAndServe(
		lb.Configuration.LoadBalancerPort,
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				lb.preProcess(lb.NextServerIndex)
				r.Header.Add("Pass-Through-Host", lb.NextServer.Address)
				processRequestFromBackend(w, r)
				lb.postProcess(lb.NextServerIndex)
				lb.decideNextServerIndex()
			},
		),
	)

	return nil
}

func (lb *LoadBalancer) decideNextServerIndex() {
	switch lb.Configuration.SchedulingAlgorithm {
	case AllowedSchedulingAlgorithms["round-robin"]:
		lb.roundRobinDecider()
	case AllowedSchedulingAlgorithms["least-connections"]:
		lb.leastConnectionsDecider()
	}
}

func (lb *LoadBalancer) preProcess(ServerIndex int) {
	selectedServer := lb.Configuration.BackendServers[ServerIndex]
	selectedServer.IncreaseActiveConnections()
	lb.Configuration.BackendServers[ServerIndex] = selectedServer
}

func (lb *LoadBalancer) postProcess(ServerIndex int) {
	selectedServer := lb.Configuration.BackendServers[ServerIndex]
	selectedServer.DecreaseActiveConnections()
	lb.Configuration.BackendServers[ServerIndex] = selectedServer
}
