package loadbalancer

// Note: Some methods to start backend servers for testing during development

import (
	"errors"
	"log"
	"net/http"
	"strings"
)

type BackendServer struct {
	Address           string
	ActiveConnections int
}

// Custom error
var (
	ErrInvalidServerAddress = errors.New("invalid server address")
)

// startBackendServers starts the backend servers in separate goroutines
func (lb *LoadBalancer) startBackendServers() error {
	if len(lb.Configuration.BackendServers) == 0 {
		return ErrNoBackendServersConfigured
	}

	for _, server := range lb.Configuration.BackendServers {
		serverParts := strings.Split(server.Address, ":")
		if len(serverParts) != 3 {
			log.Println("invalid server address: " + server.Address)
			return ErrInvalidServerAddress
		}
		port := serverParts[2]
		channel := make(chan bool, 1)
		go startServer(":"+port, channel)
		<-channel
		log.Println("Backend server started on port", port)
		close(channel)
	}
	return nil
}

// startServer starts a dummy backend server
func startServer(serverPort string, startSignal chan<- bool) error {
	startSignal <- true

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Backend server %s received request", serverPort)
		w.Write([]byte("Hello from backend server " + serverPort))
	})

	err := http.ListenAndServe(serverPort, handler)
	if err != nil {
		return err
	}
	return nil
}

func (bs *BackendServer) IncreaseActiveConnections() {
	bs.ActiveConnections++
}

func (bs *BackendServer) DecreaseActiveConnections() {
	bs.ActiveConnections--
	if bs.ActiveConnections < 0 {
		bs.ActiveConnections = 0
	}
}
