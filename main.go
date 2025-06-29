package main

import (
    "fmt"
    "log"
    "net/http"
    "path/filepath"

    settings "github.com/Yandex-Practicum/go_final_project/tests"
    "github.com/Yandex-Practicum/go_final_project/pkg/db"
    "github.com/Yandex-Practicum/go_final_project/pkg/api"
)

func main() {

    // регистрирую API-маршруты
    api.Init()

    // инициализирую БД
    dbFile := api.GetDBFile("scheduler.db")
    if err := db.Init(dbFile); err != nil {
        log.Fatalf("Не удалось инициализировать БД: %v", err)
    }
    defer db.Close()

    // определяю порт
	// беру число 7540 из tests/settings.go
    port := api.GetPort(settings.Port)

    // готовлю файловый сервер
    webDir := filepath.Join(".", "web")
    fs := http.FileServer(http.Dir(webDir))

    http.Handle("/", fs)

    // запускаю сервер
    addr := fmt.Sprintf(":%d", port)
    log.Printf("Сервер слушает http://localhost%s", addr)
    if err := http.ListenAndServe(addr, nil); err != nil {
        log.Fatalf("Ошибка сервера: %v", err)
    }
}