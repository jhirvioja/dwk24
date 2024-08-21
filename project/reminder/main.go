package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
)

type Todo struct {
	ID   int    `json:"id"`
	Todo string `json:"todo"`
}

var db *sql.DB

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3002"
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
			todo TEXT NOT NULL
	);
  `)
	if err != nil {
		log.Fatal("Failed to create or initialize todos table:", err)
	}

	randomWikiURL, err := getRandomWikipediaArticleURL()
	if err != nil {
		fmt.Println("Failed to fetch random Wikipedia article:", err)
		return
	}

	var newTodo Todo
	newTodo.Todo = "Read " + randomWikiURL

	err = db.QueryRow("INSERT INTO todos (todo) VALUES ($1) RETURNING id", newTodo.Todo).Scan(&newTodo.ID)
	if err != nil {
		fmt.Println("Failed to insert todo into database:", err)
		return
	}

	fmt.Printf("New todo created: %s with ID %d\n", newTodo.Todo, newTodo.ID)
}

func getRandomWikipediaArticleURL() (string, error) {
	resp, err := http.Get("https://en.wikipedia.org/wiki/Special:Random")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	finalURL := resp.Request.URL.String()
	return finalURL, nil
}

func createConnectionString(username, password, database string) (string, error) {
	postgresURL := fmt.Sprintf("postgres://%s:%s@psql-svc:5432/%s?sslmode=disable", username, password, database)
	return postgresURL, nil
}
