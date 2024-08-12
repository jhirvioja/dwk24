package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"
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

				if err := downloadImage("https://picsum.photos/600", imageFilePath); err != nil {
					fmt.Printf("Failed to download image: %v\n", err)
				}
			}

			mu.Unlock()
			time.Sleep(60 * time.Second)
		}
	}()

	http.Handle("/files/", http.StripPrefix("/files/", http.FileServer(http.Dir("/usr/src/app/files"))))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprintf(w, `
			<!DOCTYPE html>
			<html lang="en">
				<head>
					<meta charset="UTF-8" />
					<meta name="viewport" content="width=device-width, initial-scale=1.0" />
					<title>DwK Project</title>
					<script
						defer
						src="https://cdn.jsdelivr.net/npm/alpinejs@3.x.x/dist/cdn.min.js"
					></script>
					<script src="https://cdn.tailwindcss.com"></script>
				</head>
				<body>
					<div class="flex flex-col items-center">
						<div class="flex flex-col items-center p-8">
							<h1 class="text-3xl font-bold underline">DwK Project</h1>
							<div class="m-4 border-2 rounded-lg">
								<img
									class="p-8 max-w-full"
									src="/files/picsum_image.jpg"
									alt="Lorem Picsum Random Image"
								/>
							</div>
							<div
								x-data="{
									todos: [
										{ id: 1, text: 'Yksi' },
										{ id: 2, text: 'Kaksi' }
									],
									newTodo: '',
									nextId: 3,
									addTodo() {
										if (this.newTodo.trim() !== '') {
											this.todos.push({ id: this.nextId++, text: this.newTodo.trim() });
											this.newTodo = '';
										}
									}
								}"
								class="flex flex-col p-4 border-2 rounded-lg w-80"
							>
								<h2 class="mb-4 text-2xl font-bold">Todos</h2>
								<div class="m-2 flex flex-col">
									<label for="todoinput">Syötä todo:</label>
									<input
										id="todoinput"
										class="mt-2 rounded-md border-2"
										maxlength="140"
										x-model="newTodo"
										@keydown.enter="addTodo"
									/>
									<button class="rounded-md border-2 mt-2 max-w-40" @click="addTodo">
										Lisää todo
									</button>
								</div>
								<ul class="mt-2 pl-8 list-disc">
									<template x-for="todo in todos" :key="todo.id">
										<li x-text="todo.text"></li>
									</template>
								</ul>
							</div>
						</div>
					</div>
				</body>
			</html>
		`)
	})

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		fmt.Println("Server error:", err)
	}
}

func downloadImage(url, filePath string) error {
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to get image from URL: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("non-OK HTTP status: %s", resp.Status)
	}

	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to copy image to file: %w", err)
	}

	return nil
}
