package api

import (
    "encoding/json"
    "log"
    "net/http"
    "os"
    "strconv"
)

func writeJSON(w http.ResponseWriter, data any) {
    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
    json.NewEncoder(w).Encode(data)
}

func writeJSONError(w http.ResponseWriter, message string, statusCode int) {
    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
    w.WriteHeader(statusCode)
    json.NewEncoder(w).Encode(map[string]string{"error": message})
}

// Возвращает путь к файлу БД из переменной окружения TODO_DBFILE или значение по умолчанию
func GetDBFile(defaultFile string) string {
    if env := os.Getenv("TODO_DBFILE"); env != "" {
        return env
    }
    return defaultFile
}

// Возвращает порт из переменной окружения TODO_PORT или значение по умолчанию
func GetPort(defaultPort int) int {
    if env := os.Getenv("TODO_PORT"); env != "" {
        if p, err := strconv.Atoi(env); err == nil {
            return p
        } else {
            log.Printf("TODO_PORT=%q не число, используем %d", env, defaultPort)
        }
    }
    return defaultPort
}