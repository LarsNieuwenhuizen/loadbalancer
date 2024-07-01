package main

import (
	"log"

	"github.com/LarsNieuwenhuizen/loadbalancer"
)

func main() {
	lb := loadbalancer.LoadBalancer{}

	err := lb.Start()
	if err != nil {
		log.Fatal(err)
	}
}
