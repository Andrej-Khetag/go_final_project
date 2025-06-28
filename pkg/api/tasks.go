package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/Yandex-Practicum/go_final_project/pkg/db"
)

const maxTasks = 50

type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	tasks, err := db.Tasks(maxTasks)
	if err != nil {
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}
	
	if tasks == nil {
		tasks = []*db.Task{}
	}
	
	writeJSON(w, TasksResp{
		Tasks: tasks,
	})
}

func taskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		addTaskHandler(w, r)
	case http.MethodGet:
		id := r.FormValue("id")
		if strings.TrimSpace(id) == "" {
			writeJSON(w, map[string]string{"error": "Не указан идентификатор"})
			return
		}
		task, err := db.GetTask(id)
		if err != nil {
			writeJSON(w, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, task)
	case http.MethodPut:
		updateTaskHandler(w, r)
	case http.MethodDelete:
		id := r.FormValue("id")
		if strings.TrimSpace(id) == "" {
			writeJSON(w, map[string]string{"error": "Не указан идентификатор"})
			return
		}
		err := db.DeleteTask(id)
		if err != nil {
			writeJSON(w, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, map[string]string{})
	default:
		http.Error(w, "метод не поддерживается", http.StatusMethodNotAllowed)
	}
}

// Обрабатывает PUT-запрос для обновления задачи
func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task

	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}

	if strings.TrimSpace(task.ID) == "" {
		writeJSON(w, map[string]string{"error": "Не указан идентификатор"})
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

	err = db.UpdateTask(&task)
	if err != nil {
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, map[string]string{})
}