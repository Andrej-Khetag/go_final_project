package api

import (
    "fmt"
    "net/http"
    "time"

    "github.com/Yandex-Practicum/go_final_project/pkg/scheduler"
)

const dateFormat = "20060102"

// Регистрирует все HTTP-обработчики API
func Init() {
    http.HandleFunc("/api/nextdate", nextDateHandler)
    http.HandleFunc("/api/task", taskHandler)
    http.HandleFunc("/api/tasks", tasksHandler)
    http.HandleFunc("/api/task/done", taskDoneHandler)
    fmt.Println("API маршруты зарегистрированы")
}

// Обрабатывает GET /api/nextdate
func nextDateHandler(w http.ResponseWriter, r *http.Request) {
    nowStr := r.FormValue("now")
    if nowStr == "" {
        nowStr = time.Now().Format(dateFormat)
    }
    dateStart := r.FormValue("date")
    repeat := r.FormValue("repeat")

    now, err := time.Parse(dateFormat, nowStr)
    if err != nil {
        http.Error(w, fmt.Sprintf("invalid now %q", nowStr), http.StatusBadRequest)
        return
    }
    next, err := scheduler.NextDate(now, dateStart, repeat)
    if err != nil {
        http.Error(w, fmt.Sprintf("error: %v", err), http.StatusBadRequest)
        return
    }
    w.Header().Set("Content-Type", "text/plain; charset=utf-8")
    fmt.Fprint(w, next)
}