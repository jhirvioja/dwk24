package main

import (
	"fmt"
	"net/http"
	"os"
	"sync"
)

var counter int
var mu sync.Mutex
var filePath = "/usr/src/app/files/pongcounter.txt"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3002"
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

	err := writeCounterToFile(counter)
	if err != nil {
		http.Error(w, "Unable to write to file", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "%d", counter)

	counter++
}

func writeCounterToFile(counter int) error {
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = fmt.Fprintf(file, "%d\n", counter)
	if err != nil {
		return err
	}

	return nil
}
