package api

import (
	"net/http"
	"strings"
	"time"

	"github.com/Yandex-Practicum/go_final_project/pkg/db"
	"github.com/Yandex-Practicum/go_final_project/pkg/scheduler"
)

// Обрабатывает POST /api/task/done?id=...
func taskDoneHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}
	
	id := r.URL.Query().Get("id")
	if strings.TrimSpace(id) == "" {
		writeJSON(w, map[string]string{"error": "Не указан идентификатор"})
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}

	if strings.TrimSpace(task.Repeat) == "" {
		err := db.DeleteTask(id)
		if err != nil {
			writeJSON(w, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, map[string]string{})
		return
	}

	now := time.Now()
	next, err := scheduler.NextDate(now, task.Date, task.Repeat)
	if err != nil {
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}

	// Обновляю дату задачи на следующую дату выполнения
	err = db.UpdateDate(next, id)
	if err != nil {
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, map[string]string{})
}