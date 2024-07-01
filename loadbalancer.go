package loadbalancer

import (
	"flag"
	"io"
	"log"
	"net/http"
)

type LoadBalancer struct {
	Configuration   *AppConfig
	NextServerIndex int
	NextServer      BackendServer
}

func (lb *LoadBalancer) Start() error {
	path := flag.String("config", "", "Path to the configuration file")
	flag.Parse()
	if *path == "" {
		log.Fatal("Please provide the path as a command line argument")
	}

	err := lb.ConfigureFromYaml(*path)
	if err != nil {
		return err
	}

	log.Println("Starting load balancer on port", lb.Configuration.LoadBalancerPort)
	if !lb.Configuration.InProduction && lb.Configuration.StartGivenServers {
		lb.startBackendServers()
	}

	http.ListenAndServe(
		lb.Configuration.LoadBalancerPort,
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Println("Load balancer received request and passing to a backend server")
			r.Header.Add("Pass-Through-Host", lb.NextServer.Address)
			processRequestFromBackend(w, r)
			lb.decideNextServerIndex()
		}),
	)

	return nil
}

// processRequestFromBackend sends the request to the backend server so we can return the actual response to the client
func processRequestFromBackend(w http.ResponseWriter, r *http.Request) error {
	client := &http.Client{}
	backendServer := r.Header.Get("Pass-Through-Host")
	r.Header.Del("Pass-Through-Host")

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

func (lb *LoadBalancer) decideNextServerIndex() {
	switch lb.Configuration.SchedulingAlgorithm {
	case AllowedSchedulingAlgorithms["round-robin"]:
		lb.roundRobinDecider()
	case AllowedSchedulingAlgorithms["least-connections"]:
		// TODO: Implement least connections algorithm
		// leastConnections()
	}
}

func (lb *LoadBalancer) roundRobinDecider() {
	lb.NextServerIndex = (lb.NextServerIndex + 1) % len(lb.Configuration.BackendServers)
	lb.NextServer = lb.Configuration.BackendServers[lb.NextServerIndex]
}
