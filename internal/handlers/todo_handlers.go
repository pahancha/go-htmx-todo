package handlers

import (
    "html/template"
    "net/http"
    "path/filepath"
    "os"
    "strconv"
    "go-htmx-todo/internal/models"
)

var todos []models.Todo
var lastID int

func init() {
    todos = []models.Todo{
        {ID: 1, Title: "Learn Go", Done: false},
        {ID: 2, Title: "Learn HTMX", Done: false},
        {ID: 3, Title: "Build a todo app", Done: false},
        {ID: 4, Title: "Deploy the app", Done: false},
    }
    lastID = 4
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
    renderTemplate(w, "index.html", nil)
}

func TodosHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method == http.MethodPost {
        addTodoHandler(w, r)
        return
    }

    filter := r.URL.Query().Get("filter")
    var filteredTodos []models.Todo

    switch filter {
    case "active":
        for _, todo := range todos {
            if !todo.Done {
                filteredTodos = append(filteredTodos, todo)
            }
        }
    case "completed":
        for _, todo := range todos {
            if todo.Done {
                filteredTodos = append(filteredTodos, todo)
            }
        }
    default:
        filteredTodos = todos
    }

    renderTemplate(w, "todos.html", filteredTodos)
}

func ToggleTodoHandler(w http.ResponseWriter, r *http.Request) {
    idStr := r.URL.Query().Get("id")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        http.Error(w, "Invalid ID", http.StatusBadRequest)
        return
    }

    var toggledTodo models.Todo
    for i, todo := range todos {
        if todo.ID == id {
            todos[i].Done = !todos[i].Done
            toggledTodo = todos[i]
            break
        }
    }

    renderTemplate(w, "todo_item.html", toggledTodo)
    w.Header().Set("HX-Trigger", "todoUpdated")
}

func addTodoHandler(w http.ResponseWriter, r *http.Request) {
    title := r.FormValue("title")
    if title == "" {
        http.Error(w, "Title is required", http.StatusBadRequest)
        return
    }

    lastID++
    newTodo := models.Todo{
        ID:    lastID,
        Title: title,
        Done:  false,
    }
    todos = append(todos, newTodo)

    renderTemplate(w, "todo_item.html", newTodo)
    w.Header().Set("HX-Trigger", "todoUpdated")
}

func DeleteTodoHandler(w http.ResponseWriter, r *http.Request) {
    idStr := r.URL.Query().Get("id")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        http.Error(w, "Invalid ID", http.StatusBadRequest)
        return
    }

    for i, todo := range todos {
        if todo.ID == id {
            todos = append(todos[:i], todos[i+1:]...)
            break
        }
    }

    w.WriteHeader(http.StatusOK)
    w.Header().Set("HX-Trigger", "todoUpdated")
}

func CounterHandler(w http.ResponseWriter, r *http.Request) {
    total := len(todos)
    completed := 0
    for _, todo := range todos {
        if todo.Done {
            completed++
        }
    }

    data := struct {
        Total     int
        Completed int
    }{
        Total:     total,
        Completed: completed,
    }

    renderTemplate(w, "counter.html", data)
}

func renderTemplate(w http.ResponseWriter, tmplName string, data interface{}) {
    tmpl, err := parseTemplates(tmplName)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "text/html; charset=utf-8")
    err = tmpl.ExecuteTemplate(w, tmplName, data)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}

func parseTemplates(tmplName string) (*template.Template, error) {
    workDir, err := os.Getwd()
    if err != nil {
        return nil, err
    }

    templateDir := filepath.Join(workDir, "templates")
    pattern := filepath.Join(templateDir, "*.html")
    
    return template.ParseGlob(pattern)
}
