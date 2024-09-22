package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type Todo struct {
	ID   int    `json:"id"`
	Todo string `json:"todo"`
	Done bool   `json:"done"`
}

func FetchTodos() ([]Todo, error) {
	resp, err := http.Get("http://project-svc:5678/todos")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch todos: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("non-OK HTTP status: %s", resp.Status)
	}

	var todos []Todo
	if err := json.NewDecoder(resp.Body).Decode(&todos); err != nil {
		return nil, fmt.Errorf("failed to decode todos: %w", err)
	}

	return todos, nil
}

func PostTodo(todo Todo) error {
	todoData := Todo{
		Todo: todo.Todo,
	}
	jsonData, err := json.Marshal(todoData)
	if err != nil {
		return fmt.Errorf("failed to marshal todo: %w", err)
	}

	resp, err := http.Post("http://project-svc:5678/todos/create", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to post todo: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("non-OK HTTP status: %s", resp.Status)
	}

	return nil
}

func UpdateTodo(todo Todo) error {
	todo.Done = true

	jsonData, err := json.Marshal(todo)
	if err != nil {
		return fmt.Errorf("failed to marshal todo for update: %w", err)
	}

	req, err := http.NewRequest(http.MethodPut, "http://project-svc:5678/todos/update", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create PUT request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute PUT request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("non-OK HTTP status: %s", resp.Status)
	}

	return nil
}
