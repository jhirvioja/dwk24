package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/jhirvioja/dwk24/project/todo-frontend/services"
)

func UpdateTodoHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var todoToUpdate services.Todo
	if err := json.NewDecoder(r.Body).Decode(&todoToUpdate); err != nil {
		http.Error(w, "Failed to decode request body", http.StatusBadRequest)
		return
	}

	if err := services.UpdateTodo(todoToUpdate); err != nil {
		http.Error(w, "Failed to update todo", http.StatusInternalServerError)
		return
	}

	updatedTodos, err := services.FetchTodos()
	if err != nil {
		http.Error(w, "Failed to fetch updated todos", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(updatedTodos); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}
