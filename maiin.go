package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Task struct {
	ID           string   `json:"id"`
	Description  string   `json:"description"`
	Note         string   `json:"note"`
	Applications []string `json:"applications"`
}

var tasks = map[string]Task{
	"1": {
		ID:          "1",
		Description: "Сделать финальное задание темы REST API",
		Note:        "Если сегодня сделаю, то завтра будет свободный день. Ура!",
		Applications: []string{
			"VS Code",
			"Terminal",
			"git",
		},
	},
	"2": {
		ID:          "2",
		Description: "Протестировать финальное задание с помощью Postmen",
		Note:        "Лучше это делать в процессе разработки, каждый раз, когда запускаешь сервер и проверяешь хендлер",
		Applications: []string{
			"VS Code",
			"Terminal",
			"git",
			"Postman",
		},
	},
}

// Ниже напишите обработчики для каждого эндпоинта
// ...

// Обработчик должен принимать задачу в теле запроса и сохранять ее в мапе.
func getTasks(rw http.ResponseWriter, req *http.Request) {
	resp, err := json.Marshal(tasks)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	rw.Write(resp)
}

// Обработчик должен принимать задачу в теле запроса и сохранять ее в мапе.
func postTasks(rw http.ResponseWriter, req *http.Request) {
	var task Task
	var buf bytes.Buffer

	_, err := buf.ReadFrom(req.Body)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	if err = json.Unmarshal(buf.Bytes(), &tasks); err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	tasks[task.ID] = task

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusCreated)
}

// Обработчик должен вернуть задачу с указанным в запросе пути ID,
// если такая есть в мапе.
func getTaskById(rw http.ResponseWriter, req *http.Request) {
	id := chi.URLParam(req, "id")

	task, ok := tasks[id]
	if !ok {
		http.Error(rw, "Задача не найдена!", http.StatusNoContent)
		return
	}

	resp, err := json.Marshal(task)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	rw.Write(resp)
}

// Обработчик должен удалить задачу из мапы по её ID.
// Здесь так же нужно сначала проверить,
// есть ли задача с таким ID в мапе,
// если нет вернуть соответствующий статус.
func deleteTaskById(rw http.ResponseWriter, req *http.Request) {
	id := chi.URLParam(req, "id")

	_, ok := tasks[id]
	if !ok {
		http.Error(rw, "Задача не найдена!", http.StatusBadRequest)
		return
	}

	delete(tasks, id)
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
}

func main() {
	r := chi.NewRouter()

	// здесь регистрируйте ваши обработчики
	r.Get("/tasks", getTasks)
	r.Post("/tasks", postTasks)
	r.Get("/tasks/{id}", getTaskById)
	r.Delete("/tasks/{id}", deleteTaskById)

	if err := http.ListenAndServe(":8080", r); err != nil {
		fmt.Printf("Ошибка при запуске сервера: %s", err.Error())
		return
	}
}
