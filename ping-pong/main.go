package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"

	_ "github.com/lib/pq"
)

var db *sql.DB
var mu sync.Mutex

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

	db, err = sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	_, err = db.Exec(`
		Create TABLE IF NOT exists counter (
				id SERIAL PRIMARY KEY,
				value INT NOT NULL
		);
		INSERT INTO counter (value) SELECT 0 WHERE NOT EXISTS (SELECT 1 FROM counter);
	`)
	if err != nil {
		log.Fatal("Failed to create or initialize counter table:", err)
	}

	fmt.Printf("Server started in port %s\n", port)

	http.HandleFunc("/", pongHandler)
	http.HandleFunc("/healthz", healthCheckHandler)

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		fmt.Println("Server error:", err)
	}
}

func pongHandler(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	var counter int
	err := db.QueryRow("SELECT value FROM counter LIMIT 1").Scan(&counter)
	if err != nil {
		http.Error(w, "Failed to read counter from database", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "%d", counter)

	_, err = db.Exec("UPDATE counter SET value = value + 1")
	if err != nil {
		http.Error(w, "Failed to update counter in database", http.StatusInternalServerError)
		return
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
