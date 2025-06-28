package main

import (
    "fmt"
    "log"
    "net/http"
    "os"
    "path/filepath"
    "strconv"

    settings "github.com/Yandex-Practicum/go_final_project/tests"
    "github.com/Yandex-Practicum/go_final_project/pkg/db"
    "github.com/Yandex-Practicum/go_final_project/pkg/api"
)

func main() {

    // регистрирую API-маршруты
    api.Init()

    // инициализирую БД
    dbFile := "scheduler.db"
    if env := os.Getenv("TODO_DBFILE"); env != "" {
        dbFile = env
    }
    if err := db.Init(dbFile); err != nil {
        log.Fatalf("Не удалось инициализировать БД: %v", err)
    }
    defer db.Close()

    // определяю порт
	// беру число 7540 из tests/settings.go
    port := settings.Port
    if env := os.Getenv("TODO_PORT"); env != "" {
        // если есть переменная окружения TODO_PORT, пытаюсь её прочитать как число
        if p, err := strconv.Atoi(env); err == nil {
            port = p
        } else {
            log.Printf("TODO_PORT=%q не число, используем %d", env, port)
        }
    }

    // готовлю файловый сервер
    webDir := filepath.Join(".", "web")
    fs := http.FileServer(http.Dir(webDir))

    http.Handle("/", fs)

    // запускаю сервер
    addr := fmt.Sprintf(":%d", port)
    log.Printf("Сервер слушает http://localhost%s …", addr)
    if err := http.ListenAndServe(addr, nil); err != nil {
        log.Fatalf("Ошибка сервера: %v", err)
    }
}
