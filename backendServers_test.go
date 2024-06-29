package loadbalancer

import (
	"testing"
	"time"
)

func TestStartBackendServers(t *testing.T) {
	testCases := []struct {
		name          string
		Configuration AppConfig
		expected      error
	}{
		{
			name:     "Wrong server address configured to load balance returns an error",
			expected: ErrInvalidServerAddress,
			Configuration: AppConfig{
				BackendServers:    map[int]string{0: "http://localhost", 1: "http://localhost:8081"},
				LoadBalancerPort:  ":8080",
				InProduction:      false,
				StartGivenServers: true,
			},
		},
		{
			name:     "Start backend servers correctly",
			expected: nil,
			Configuration: AppConfig{
				BackendServers:    map[int]string{0: "http://localhost:8081"},
				LoadBalancerPort:  ":8080",
				InProduction:      false,
				StartGivenServers: true,
			},
		},
		{
			name:     "No backend servers configured to load balance returns an error",
			expected: ErrNoBackendServersConfigured,
			Configuration: AppConfig{
				BackendServers:    map[int]string{},
				LoadBalancerPort:  ":8080",
				InProduction:      false,
				StartGivenServers: true,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			Configuration = &tc.Configuration
			err := startBackendServers()
			if err != tc.expected {
				t.Errorf("Expected %v, got %v", tc.expected, err)
			}
		})
	}
}

func TestGoroutineStartsWhenStartingServer(t *testing.T) {
	startSignal := make(chan bool, 1)
	go startServer(":8083", startSignal)

	select {
	case <-startSignal:
		t.Log("Goroutine started")
	case <-time.After(1 * time.Second):
		t.Error("Timeout waiting for goroutine to start")
	}
}
