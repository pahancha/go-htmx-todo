package main

import (
    "fmt"
    "log"
    "net/http"

	"go-htmx-todo/internal/handlers"
)

func main() {
    http.HandleFunc("/", handlers.IndexHandler)
    http.HandleFunc("/todos", handlers.TodosHandler)
	http.HandleFunc("/toggle", handlers.ToggleTodoHandler)
	http.HandleFunc("/counter", handlers.CounterHandler)
	http.HandleFunc("/todos/delete", handlers.DeleteTodoHandler)

    fmt.Println("Server is running on http://localhost:8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}
