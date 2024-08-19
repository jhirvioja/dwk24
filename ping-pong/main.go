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

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL environment variable not set")
	}

	var err error
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
