package main

import (
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/google/uuid"
)

type Status struct {
	Timestamp string
	Rstring   string
}

var status Status
var mu sync.RWMutex

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	fmt.Printf("Server started in port %s\n", port)

	randomString := uuid.New().String()

	go func() {
		for {
			currentTime := time.Now().Format(time.RFC3339Nano)
			mu.Lock()
			status = Status{Timestamp: currentTime, Rstring: randomString}
			mu.Unlock()
			fmt.Printf("%s: %s\n", currentTime, randomString)
			time.Sleep(5 * time.Second)
		}
	}()

	http.HandleFunc("/", statusHandler)

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		fmt.Println("Server error:", err)
	}
}

func statusHandler(w http.ResponseWriter, r *http.Request) {
	mu.RLock()
	defer mu.RUnlock()

	fmt.Fprintf(w, "Timestamp: %s, String: %s\n", status.Timestamp, status.Rstring)
}
