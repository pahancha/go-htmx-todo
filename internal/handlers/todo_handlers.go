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

func init() {
    todos = []models.Todo{
        {ID: 1, Title: "Learn Go", Done: false},
        {ID: 2, Title: "Learn HTMX", Done: false},
        {ID: 3, Title: "Build a todo app", Done: false},
        {ID: 4, Title: "Deploy the app", Done: false},
    }
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
    renderTemplate(w, "index.html", todos)
}

func TodosHandler(w http.ResponseWriter, r *http.Request) {
    renderTemplate(w, "todos.html", todos)
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