package main

import (
	"fmt"
	"net/http"
	"os"
	"sync"
)

var counter int
var mu sync.Mutex

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3001"
	}

	fmt.Printf("Server started in port %s\n", port)

	http.HandleFunc("/", pongHandler)

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		fmt.Println("Server error:", err)
	}
}

func pongHandler(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	fmt.Fprintf(w, "pong %d", counter)
	counter++
}
