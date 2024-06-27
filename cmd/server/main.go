package main

import (
	"log"

	"github.com/LarsNieuwenhuizen/loadbalancer/pkg/server"
)

func main() {
	server.Configuration.SetBackendServers(map[int]string{})
	log.Println(server.Configuration)
	server.StartLoadBalancer()
}
