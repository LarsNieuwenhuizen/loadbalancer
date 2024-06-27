package loadbalancer

import (
	"log"
	"net/http"
)

var (
	Configuration     = Setup()
	chosenServerIndex = 0
)

func Start() {
	log.Println("Starting load balancer on port", Configuration.LoadBalancerPort)
	if !Configuration.InProduction {
		startBackendServers()
	}
	http.ListenAndServe(Configuration.LoadBalancerPort, http.HandlerFunc(proxyHandler))
}

func proxyHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Load balancer received request")
	log.Println(
		"Proxying request to backend server",
		Configuration.BackendServers[chosenServerIndex],
	)
	http.Redirect(
		w,
		r,
		Configuration.BackendServers[chosenServerIndex],
		http.StatusTemporaryRedirect,
	)

	decideNextServerIndex()
}

func decideNextServerIndex() {
	switch Configuration.SchedulingAlgorithm {
	case AllowedSchedulingAlgorithms["round-robin"]:
		roundRobinDecider()
	case AllowedSchedulingAlgorithms["least-connections"]:
		// TODO: Implement least connections algorithm
		// leastConnections()
	}
}

func roundRobinDecider() {
	chosenServerIndex = (chosenServerIndex + 1) % len(Configuration.BackendServers)
}
