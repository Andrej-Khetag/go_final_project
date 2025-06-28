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
		writeJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	task.Title = strings.TrimSpace(task.Title)
	if task.Title == "" {
		writeJSONError(w, "Не указан заголовок задачи", http.StatusBadRequest)
		return
	}

	err = checkDate(&task)
	if err != nil {
		writeJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Добавляю задачи в БД
	id, err := db.AddTask(&task)
	if err != nil {
		writeJSONError(w, err.Error(), http.StatusInternalServerError)
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

	if strings.TrimSpace(task.Repeat) != "" {
		// Проверяю корректность правила повторения
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

// Сравает по дате (без времени): возвращает true, если t раньше сегодняшнего дня
func isPastDate(now time.Time, t time.Time) bool {
	nowStr := now.Format("20060102")
	tStr := t.Format("20060102")
	return tStr < nowStr
}