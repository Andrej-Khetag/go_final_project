package api

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/Yandex-Practicum/go_final_project/pkg/db"
	"github.com/Yandex-Practicum/go_final_project/pkg/scheduler"
)

// Обработчик POST-запроса на /api/task
func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task

	// Декодирую JSON-запрос в структуру задачи
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}

	task.Title = strings.TrimSpace(task.Title)
	if task.Title == "" {
		writeJSON(w, map[string]string{"error": "Не указан заголовок задачи"})
		return
	}

	err = checkDate(&task)
	if err != nil {
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}

	id, err := db.AddTask(&task)
	if err != nil {
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, map[string]string{"id": id})
}

// Нормализует и проверяет дату задачи
func checkDate(task *db.Task) error {
	now := time.Now()
	nowStr := now.Format("20060102")

	if strings.TrimSpace(task.Date) == "" {
		task.Date = nowStr
	}

	t, err := time.Parse("20060102", task.Date)
	if err != nil {
		return err
	}

	// Если указано повторение - всегда проверяю и вычисляю следующую дату
	if strings.TrimSpace(task.Repeat) != "" {
		next, err := scheduler.NextDate(now, task.Date, task.Repeat)
		if err != nil {
			return err
		}
		if isPastDate(now, t) {
			task.Date = next
		}
	} else {
		if isPastDate(now, t) {
			task.Date = nowStr
		}
	}

	return nil
}

// Сравнение по дате (без времени): возвращает true, если t раньше сегодняшнего дня
func isPastDate(now time.Time, t time.Time) bool {
	nowStr := now.Format("20060102")
	tStr := t.Format("20060102")
	return tStr < nowStr
}