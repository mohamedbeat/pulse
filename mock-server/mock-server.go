// mock-server.go
package main

import (
	"fmt"
	"net/http"
	"time"
)

func main() {
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		// w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, `{"status":"ok","uptime":"10m"}`)
	})

	http.HandleFunc("/fail", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "internal error", http.StatusInternalServerError)
	})

	http.HandleFunc("/slow", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(8 * time.Second) // Simulate timeout (if your timeout <8s)
		w.WriteHeader(http.StatusOK)
	})

	fmt.Println("Mock server running on :9000")
	http.ListenAndServe(":9000", nil)
}
