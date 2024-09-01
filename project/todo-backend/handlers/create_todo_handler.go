package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
)

type Todo struct {
	ID   int    `json:"id"`
	Todo string `json:"todo"`
}

const maxTodoLength = 140

func CreateTodoHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method == http.MethodPost {
		var newTodo Todo
		if err := json.NewDecoder(r.Body).Decode(&newTodo); err != nil {
			http.Error(w, "Failed to decode request body", http.StatusBadRequest)
			return
		}

		fmt.Println("Received todo:", newTodo.Todo)

		if len(newTodo.Todo) > maxTodoLength {
			http.Error(w, "Todo exceeds 140 characters", http.StatusBadRequest)
			fmt.Println("Todo exceeds 140 characters:", newTodo.Todo)
			return
		}

		err := db.QueryRow("INSERT INTO todos (todo) VALUES ($1) RETURNING id", newTodo.Todo).Scan(&newTodo.ID)
		if err != nil {
			http.Error(w, "Failed to insert todo into database", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(newTodo); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		}

		fmt.Println("Todo created succesfully:", newTodo.Todo)
		return
	}

	http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
}
