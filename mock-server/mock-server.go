// mock-server.go
package main

import (
	"fmt"
	"net/http"
	"time"
)

func main() {
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		// Log all headers for debugging
		fmt.Println("=== Request Headers ===")
		for name, values := range r.Header {
			fmt.Printf("  %s: %v\n", name, values)
		}
		fmt.Println("======================")

		// Check for specific header (case-insensitive)
		idkHeader := r.Header.Get("Idk") // Go canonicalizes to "Idk"
		if idkHeader == "" {
			idkHeader = r.Header.Get("idk") // Try lowercase too
		}
		fmt.Printf("Idk header value: %q\n", idkHeader)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, `{"status":"ok","uptime":"10m"}`)
	})
	http.HandleFunc("/latency", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(100 * time.Millisecond)

		// Log headers for this endpoint too
		fmt.Println("=== Latency Request Headers ===")
		for name, values := range r.Header {
			fmt.Printf("  %s: %v\n", name, values)
		}
		fmt.Println("================================")

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
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
