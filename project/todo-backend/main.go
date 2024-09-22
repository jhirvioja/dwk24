package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/jhirvioja/dwk24/project/todo-backend/handlers"
	_ "github.com/lib/pq"
)

type Todo struct {
	ID   int    `json:"id"`
	Todo string `json:"todo"`
	Done bool   `json:"done"`
}

var db *sql.DB

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3001"
	}

	dbUsername := os.Getenv("DB_USERNAME")
	if dbUsername == "" {
		log.Fatal("DB_USERNAME environment variable not set")
	}

	dbPassword := os.Getenv("DB_PASSWORD")
	if dbPassword == "" {
		log.Fatal("DB_PASSWORD environment variable not set")
	}

	dbDatabase := os.Getenv("DB_DATABASE")
	if dbDatabase == "" {
		log.Fatal("DB_DATABASE environment variable not set")
	}

	dbURL, err := createConnectionString(dbUsername, dbPassword, dbDatabase)
	if err != nil {
		log.Fatalf("Error creating PostgreSQL URL: %v", err)
	}

	db, err = sql.Open("postgres", string(dbURL))
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	_, err = db.Exec(`
	Create TABLE IF NOT exists todos (
			id SERIAL PRIMARY KEY,
			todo TEXT NOT NULL,
			done BOOLEAN DEFAULT FALSE
	);
  `)
	if err != nil {
		log.Fatal("Failed to create or initialize todos table:", err)
	}

	http.HandleFunc("/todos", func(w http.ResponseWriter, r *http.Request) {
		handlers.GetTodosHandler(w, r, db)
	})

	http.HandleFunc("/todos/create", func(w http.ResponseWriter, r *http.Request) {
		handlers.CreateTodoHandler(w, r, db)
	})

	http.HandleFunc("/todos/update", func(w http.ResponseWriter, r *http.Request) {
		handlers.UpdateTodoHandler(w, r, db)
	})

	http.HandleFunc("/healthz", healthCheckHandler)

	fmt.Printf("Server started on port %s\n", port)

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		fmt.Println("Server error:", err)
	}
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	err := db.Ping()
	if err != nil {
		http.Error(w, "Database connection failed", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func createConnectionString(username, password, database string) (string, error) {
	postgresURL := fmt.Sprintf("postgres://%s:%s@psql-svc:5432/%s?sslmode=disable", username, password, database)
	return postgresURL, nil
}
