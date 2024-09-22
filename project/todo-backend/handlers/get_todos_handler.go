package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
)

func GetTodosHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method == http.MethodGet {
		rows, err := db.Query("SELECT id, todo, done FROM todos")
		if err != nil {
			http.Error(w, "Failed to retrieve todos from database", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var todos []Todo
		for rows.Next() {
			var todo Todo
			if err := rows.Scan(&todo.ID, &todo.Todo, &todo.Done); err != nil {
				http.Error(w, "Failed to scan todo", http.StatusInternalServerError)
				return
			}
			todos = append(todos, todo)
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(todos); err != nil {
			http.Error(w, "Failed to encode todos", http.StatusInternalServerError)
		}
		return
	}

	http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
}
