package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type Status struct {
	Timestamp string
	Rstring   string
}

var mu sync.RWMutex

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3001"
	}

	fmt.Printf("Server started in port %s\n", port)

	filePath := "/usr/src/app/files/timestamp.txt"
	dirPath := filepath.Dir(filePath)

	if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
		fmt.Printf("Failed to create directory: %v\n", err)
		return
	}

	go func() {
		for {
			currentTime := time.Now().Format(time.RFC3339Nano)
			mu.Lock()

			if err := os.WriteFile(filePath, []byte(currentTime), 0644); err != nil {
				fmt.Printf("Failed to write to file: %v\n", err)
			}

			mu.Unlock()
			time.Sleep(5 * time.Second)
		}
	}()

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		fmt.Println("Server error:", err)
	}
}
