package loadbalancer

import (
	"io"
	"net/http"
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
				BackendServers: map[int]BackendServer{
					0: {
						Address: "http://localhost",
					},
					1: {
						Address: "http://localhost:8081",
					},
				},
				LoadBalancerPort:  ":8080",
				InProduction:      false,
				StartGivenServers: true,
			},
		},
		{
			name:     "Start backend servers correctly",
			expected: nil,
			Configuration: AppConfig{
				BackendServers: map[int]BackendServer{
					0: {
						Address: "http://localhost:8080",
					},
				},
				LoadBalancerPort:  ":8080",
				InProduction:      false,
				StartGivenServers: true,
			},
		},
		{
			name:     "No backend servers configured to load balance returns an error",
			expected: ErrNoBackendServersConfigured,
			Configuration: AppConfig{
				BackendServers:    map[int]BackendServer{},
				LoadBalancerPort:  ":8080",
				InProduction:      false,
				StartGivenServers: true,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			lb := LoadBalancer{
				Configuration: &tc.Configuration,
			}
			err := lb.startBackendServers()
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

func TestBackendServerHandlerResponse(t *testing.T) {
	testCases := []struct {
		name     string
		port     string
		expected string
	}{
		{
			name:     "Backend server on port 9000 returns expected response",
			port:     ":9000",
			expected: "Hello from backend server :9000",
		},
		{
			name:     "Backend server on port 9000 returns expected response",
			port:     ":9001",
			expected: "Hello from backend server :9001",
		},
	}

	client := &http.Client{}
	for _, tc := range testCases {
		channel := make(chan bool, 1)
		go startServer(tc.port, channel)
		<-channel
		close(channel)

		url := "http://localhost" + tc.port

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			t.Fatal(err)
		}
		resp, err := client.Do(req)
		if err != nil {
			t.Error(err)
		}

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status OK, got %v", resp.Status)
		}
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Error(err)
		}

		result := string(body)
		if result != tc.expected {
			t.Errorf("Expected %v', got %v", tc.expected, result)
		}
	}
}
