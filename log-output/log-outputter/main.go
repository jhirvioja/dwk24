package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/google/uuid"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	fmt.Printf("Server started on port %s\n", port)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		timestampFilePath := "/usr/src/app/files/timestamp.txt"
		counterFilePath := "/usr/src/app/files/pongcounter.txt"

		timestamp, err := os.ReadFile(timestampFilePath)
		if err != nil {
			fmt.Printf("Failed to read the timestamp file: %v\n", err)
			http.Error(w, "Failed to read the timestamp file", http.StatusInternalServerError)
			return
		}

		counter, err := os.ReadFile(counterFilePath)
		if err != nil {
			fmt.Printf("Failed to read the pong counter file: %v\n", err)
			http.Error(w, "Failed to read the pong counter file", http.StatusInternalServerError)
			return
		}

		randomString := uuid.New().String()

		fmt.Fprintf(w, "%s: %s.\nPing / Pongs: %s\n", string(timestamp), randomString, string(counter))
	})

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		fmt.Println("Server error:", err)
	}
}
