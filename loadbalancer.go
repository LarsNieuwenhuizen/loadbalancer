package loadbalancer

import (
	"flag"
	"io"
	"log"
	"net/http"
)

var (
	Configuration     = Setup()
	chosenServerIndex = 0
)

func Start() {
	path := flag.String("config", "", "Path to the configuration file")
	flag.Parse()
	if *path == "" {
		log.Fatal("Please provide the path as a command line argument")
	}

	Configuration.LoadFromYaml(*path)

	log.Println("Starting load balancer on port", Configuration.LoadBalancerPort)
	if !Configuration.InProduction && Configuration.StartGivenServers {
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

	err := processRequestFromBackend(w, r)
	if err != nil {
		log.Fatal("Error processing request from backend:", err)
	}

	decideNextServerIndex()
}

// processRequestFromBackend sends the request to the backend server so we can return the actual response to the client
func processRequestFromBackend(w http.ResponseWriter, r *http.Request) error {
	client := &http.Client{}
	backendServer := Configuration.BackendServers[chosenServerIndex]

	// Create a new request to the backend server
	backendRequest, err := http.NewRequest(r.Method, backendServer, r.Body)
	if err != nil {
		log.Println("Error creating request to backend server:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return err
	}

	// Copy headers from the original request to the backend request
	backendRequest.Header = r.Header

	// Send the request to the backend server
	backendResponse, err := client.Do(backendRequest)
	if err != nil {
		log.Println("Error sending request to backend server:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return err
	}
	defer backendResponse.Body.Close()

	// Copy the response headers from the backend response to the client response
	for key, values := range backendResponse.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}

	// Set the status code of the client response to the status code of the backend response
	w.WriteHeader(backendResponse.StatusCode)

	// Copy the response body from the backend response to the client response
	_, err = io.Copy(w, backendResponse.Body)
	if err != nil {
		log.Println("Error copying response body:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return err
	}

	return nil
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
