package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
)

func UpdateTodoHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method == http.MethodPut {
		var todo Todo
		if err := json.NewDecoder(r.Body).Decode(&todo); err != nil {
			http.Error(w, "Failed to decode request body", http.StatusBadRequest)
			return
		}

		result, err := db.Exec("UPDATE todos SET done = true WHERE id = $1 AND done = false", todo.ID)
		if err != nil {
			http.Error(w, "Failed to update todo", http.StatusInternalServerError)
			return
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			http.Error(w, "Failed to retrieve rows affected", http.StatusInternalServerError)
			return
		}

		if rowsAffected == 0 {
			http.Error(w, "No todo found to update or todo already done", http.StatusNotFound)
			return
		}

		todo.Done = true
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(todo); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		}

		fmt.Printf("Todo with ID %d marked as done.\n", todo.ID)
		return
	}

	http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
}
