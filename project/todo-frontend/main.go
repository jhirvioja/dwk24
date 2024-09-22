package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/jhirvioja/dwk24/project/todo-frontend/handlers"
	"github.com/jhirvioja/dwk24/project/todo-frontend/services"
)

var mu sync.RWMutex

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	fmt.Printf("Server started in port %s\n", port)

	timestampFilePath := "/usr/src/app/files/project_timestamp.txt"
	imageFilePath := "/usr/src/app/files/picsum_image.jpg"
	filesDirPath := filepath.Dir(timestampFilePath)

	if err := os.MkdirAll(filesDirPath, os.ModePerm); err != nil {
		fmt.Printf("Failed to create assets directory: %v\n", err)
		return
	}

	go func() {
		for {
			mu.Lock()

			var fileTimestamp time.Time

			if data, err := os.ReadFile(timestampFilePath); err == nil {
				if ts, err := time.Parse(time.RFC3339Nano, string(data)); err == nil {
					fileTimestamp = ts
				} else {
					fmt.Printf("Failed to parse timestamp: %v\n", err)
					fileTimestamp = time.Time{}
				}
			} else {
				fmt.Printf("Failed to read the timestamp file: %v\n", err)
				fileTimestamp = time.Time{}
			}

			currentTime := time.Now()

			if currentTime.Sub(fileTimestamp) > 60*time.Minute {
				if err := os.WriteFile(timestampFilePath, []byte(currentTime.Format(time.RFC3339Nano)), 0644); err != nil {
					fmt.Printf("Failed to write to file: %v\n", err)
				}

				if err := services.DownloadImage("https://picsum.photos/600", imageFilePath); err != nil {
					fmt.Printf("Failed to download image: %v\n", err)
				}
			}

			mu.Unlock()
			time.Sleep(60 * time.Second)
		}
	}()

	http.Handle("/files/", http.StripPrefix("/files/", http.FileServer(http.Dir("/usr/src/app/files"))))

	http.HandleFunc("/", handlers.TodoHandler)

	http.HandleFunc("/add_todo", handlers.AddTodoHandler)

	http.HandleFunc("/update_todo", handlers.UpdateTodoHandler)

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		fmt.Println("Server error:", err)
	}
}
