package loadbalancer

import (
	"io"
	"log"
	"net/http"
)

// processRequestFromBackend sends the request to the backend server so we can return the actual response to the client
func processRequestFromBackend(w http.ResponseWriter, r *http.Request) error {
	// Sleep for 1 second to simulate a slow backend server
	// time.Sleep(1 * time.Second)
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
