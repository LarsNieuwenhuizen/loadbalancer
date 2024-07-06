package loadbalancer

func (lb *LoadBalancer) roundRobinDecider() {
	lb.NextServerIndex = (lb.NextServerIndex + 1) % len(lb.Configuration.BackendServers)
	lb.NextServer = lb.Configuration.BackendServers[lb.NextServerIndex]
}

func (lb *LoadBalancer) leastConnectionsDecider() {
	backendServers := lb.Configuration.BackendServers
	var selectedServer BackendServer
	var index int

	lowestNumber := 0
	for i, server := range backendServers {
		if server.ActiveConnections <= lowestNumber {
			lowestNumber = server.ActiveConnections
			selectedServer = server
			index = i
		}
	}

	lb.NextServer = selectedServer
	lb.NextServerIndex = index
	lb.Configuration.BackendServers[index] = selectedServer
}
