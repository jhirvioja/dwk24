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
		filePath := "/usr/src/app/files/timestamp.txt"

		timestamp, err := os.ReadFile(filePath)
		if err != nil {
			fmt.Printf("Failed to read the file: %v\n", err)
			http.Error(w, "Failed to read the timestamp file", http.StatusInternalServerError)
			return
		}

		randomString := uuid.New().String()

		fmt.Fprintf(w, "Timestamp: %s, String: %s\n", timestamp, randomString)
	})

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		fmt.Println("Server error:", err)
	}
}
