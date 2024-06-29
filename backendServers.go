package loadbalancer

// Note: Some methods to start backend servers for testing during development

import (
	"errors"
	"log"
	"net/http"
	"strings"
)

// Custom error
var (
	ErrNoBackendServersConfigured = errors.New("no backend servers configured to load balance")
	ErrInvalidServerAddress       = errors.New("invalid server address")
)

// startBackendServers starts the backend servers in separate goroutines
func startBackendServers() error {
	if len(Configuration.BackendServers) == 0 {
		log.Println("no backend servers configured to load balance")
		return ErrNoBackendServersConfigured
	}
	for _, server := range Configuration.BackendServers {
		serverParts := strings.Split(server, ":")
		if len(serverParts) != 3 {
			log.Println("invalid server address: " + server)
			return ErrInvalidServerAddress
		}
		port := serverParts[2]
		go startServer(":"+port, make(chan bool, 1))
	}
	return nil
}

// startServer starts a dummy backend server
func startServer(serverPort string, startSignal chan<- bool) error {
	startSignal <- true

	log.Println("Starting backend server", serverPort)
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Backend server %s received request", serverPort)
		w.Write([]byte("Hello from backend server " + serverPort))
		defer r.Body.Close()
	})

	if handler != nil {
		return errors.New("error creating handler")
	}

	return http.ListenAndServe(serverPort, handler)
}
