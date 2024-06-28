package loadbalancer

// Note: Some methods to start backend servers for testing during development

import (
	"log"
	"net/http"
	"strings"
)

// startBackendServers starts the backend servers in separate goroutines
func startBackendServers() {
	if len(Configuration.BackendServers) == 0 {
		log.Fatal("no backend servers configured to load balance")
		return
	}
	for _, server := range Configuration.BackendServers {
		serverParts := strings.Split(server, ":")
		if len(serverParts) != 3 {
			log.Fatalf("invalid server address: %s", serverParts)
			return
		}
		port := serverParts[2]
		log.Println(":" + port)
		go startServer(":" + port)
	}
}

// startServer starts a dummy backend server
func startServer(serverPort string) {
	log.Println("Starting backend server", serverPort)
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Backend server %s received request", serverPort)
		w.Write([]byte("Hello from backend server " + serverPort))
	})
	http.ListenAndServe(serverPort, handler)
}
