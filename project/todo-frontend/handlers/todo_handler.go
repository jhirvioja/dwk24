package handlers

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"

	"github.com/jhirvioja/dwk24/project/todo-frontend/services"
)

func TodoHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	todos, err := services.FetchTodos()
	if err != nil {
		fmt.Printf("Failed to fetch todos: %v\n", err)
		http.Error(w, "Failed to fetch todos", http.StatusInternalServerError)
		return
	}

	todosJSON, err := json.Marshal(todos)
	if err != nil {
		fmt.Printf("Failed to marshal todos: %v\n", err)
		http.Error(w, "Failed to encode todos", http.StatusInternalServerError)
		return
	}

	todoTemplatePath := filepath.Join("templates", "todo.tmpl")

	t, err := template.ParseFiles(todoTemplatePath)
	if err != nil {
		fmt.Printf("Failed to parse template file: %v\n", err)
		http.Error(w, "Failed to parse template file", http.StatusInternalServerError)
		return
	}

	err = t.Execute(w, map[string]interface{}{
		"Todos": template.JS(todosJSON),
	})
	if err != nil {
		http.Error(w, "Failed to execute template", http.StatusInternalServerError)
		return
	}
}
