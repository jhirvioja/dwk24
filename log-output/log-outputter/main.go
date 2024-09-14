package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	fileContent, err := os.ReadFile("/etc/config/information.txt")
	if err != nil {
		fmt.Printf("Failed to read file: %v\n", err)
		os.Exit(1)
	}

	message := os.Getenv("MESSAGE")
	if message == "" {
		message = "MESSAGE is not set"
	}

	fmt.Printf("file content: %s", string(fileContent))
	fmt.Printf("env variable: MESSAGE=%s\n", message)

	fmt.Printf("Server started on port %s\n", port)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		resp, err := http.Get("http://ping-pong-svc:80")
		if err != nil {
			fmt.Printf("Failed to fetch counter from ping-pong-svc: %v\n", err)
			http.Error(w, "Failed to fetch counter", http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		counter, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("Failed to read response body: %v\n", err)
			http.Error(w, "Failed to read response", http.StatusInternalServerError)
			return
		}

		randomString := uuid.New().String()
		timestamp := time.Now().Format(time.RFC3339Nano)

		fmt.Fprintf(w, "%s: %s.\nPing / Pongs: %s\n", string(timestamp), randomString, string(counter))
	})

	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		resp, err := http.Get("http://ping-pong-svc:80/healthz")
		if err != nil || resp.StatusCode != http.StatusOK {
			fmt.Printf("Health check failed: %v\n", err)
			http.Error(w, "Service is not ready", http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		fmt.Println("Server error:", err)
	}

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		fmt.Println("Server error:", err)
	}
}
