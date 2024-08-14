package handlers

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"

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

	tmpl := `
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
								todos: {{.Todos}},
								newTodo: '',
								addTodo() {
									if (this.newTodo.trim() !== '') {
										fetch('/add_todo', {
											method: 'POST',
											headers: {
												'Content-Type': 'application/json',
											},
											body: JSON.stringify({ todo: this.newTodo.trim() }),
										})
										.then(response => response.json())
										.then(data => {
											this.todos = data;
											this.newTodo = '';
										})
										.catch(error => console.error('Error:', error));
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
									<li x-text="todo.todo"></li>
								</template>
							</ul>
						</div>
					</div>
				</div>
			</body>
		</html>
	`

	t, err := template.New("webpage").Parse(tmpl)
	if err != nil {
		http.Error(w, "Failed to parse template", http.StatusInternalServerError)
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
