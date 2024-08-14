package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/jhirvioja/dwk24/project/todo-frontend/services"
)

func AddTodoHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var newTodo services.Todo
	if err := json.NewDecoder(r.Body).Decode(&newTodo); err != nil {
		http.Error(w, "Failed to decode request body", http.StatusBadRequest)
		return
	}

	if err := services.PostTodo(newTodo); err != nil {
		http.Error(w, "Failed to create todo", http.StatusInternalServerError)
		return
	}

	updatedTodos, err := services.FetchTodos()
	if err != nil {
		http.Error(w, "Failed to fetch todos", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedTodos)
}
