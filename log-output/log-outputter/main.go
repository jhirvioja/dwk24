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

	fmt.Printf("Server started on port %s\n", port)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		resp, err := http.Get("http://ping-pong-svc:5678")
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

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		fmt.Println("Server error:", err)
	}
}
