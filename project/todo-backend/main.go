package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sync"
)

type Todo struct {
	ID   int    `json:"id"`
	Todo string `json:"todo"`
}

var (
	todos    []Todo
	nextID   = 1
	todoLock sync.Mutex
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3001"
	}

	todos = append(todos, Todo{ID: nextID, Todo: "Sample todo"})
	nextID++

	http.HandleFunc("/todos", getTodosHandler)
	http.HandleFunc("/todos/create", createTodoHandler)

	fmt.Printf("Server started on port %s\n", port)

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		fmt.Println("Server error:", err)
	}
}

func getTodosHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		todoLock.Lock()
		defer todoLock.Unlock()

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(todos); err != nil {
			http.Error(w, "Failed to encode todos", http.StatusInternalServerError)
		}
		return
	}

	http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
}

func createTodoHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		var newTodo Todo
		if err := json.NewDecoder(r.Body).Decode(&newTodo); err != nil {
			http.Error(w, "Failed to decode request body", http.StatusBadRequest)
			return
		}

		todoLock.Lock()
		defer todoLock.Unlock()

		newTodo.ID = nextID
		nextID++
		todos = append(todos, newTodo)

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(newTodo); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		}
		return
	}

	http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
}
