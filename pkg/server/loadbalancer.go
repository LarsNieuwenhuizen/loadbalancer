package server

import (
	"log"
	"net/http"

	"github.com/LarsNieuwenhuizen/loadbalancer/pkg/config"
)

var (
	Configuration     = config.Setup()
	chosenServerIndex = 0
)

func StartLoadBalancer() {
	log.Println("Starting load balancer on port", Configuration.LoadBalancerPort)
	if !Configuration.InProduction {
		startBackendServers()
	}
	http.ListenAndServe(Configuration.LoadBalancerPort, http.HandlerFunc(ProxyHandler))
}

func ProxyHandler(w http.ResponseWriter, r *http.Request) {
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
	case config.AllowedSchedulingAlgorithms["round-robin"]:
		roundRobinDecider()
	case config.AllowedSchedulingAlgorithms["least-connections"]:
		// TODO: Implement least connections algorithm
		// leastConnections()
	}
}

func roundRobinDecider() {
	chosenServerIndex = (chosenServerIndex + 1) % len(Configuration.BackendServers)
}
