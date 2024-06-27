package server

import (
	"log"
	"net/http"
)

func startBackendServers() {
	if len(Configuration.BackendServers) == 0 {
		log.Fatal("no backend servers configured to load balance")
		return
	}
	for _, server := range Configuration.BackendServers {
		go startServer(server)
	}
}

func startServer(server string) {
	log.Println("Starting backend server", server)
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Backend server %s received request", server)
	})
	http.ListenAndServe(server, handler)
}
