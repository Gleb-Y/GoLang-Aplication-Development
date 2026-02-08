package handlers

import (
	"assignment/internal/models"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

type TaskHandler struct {
	store *models.TaskStore
}

func NewTaskHandler(store *models.TaskStore) *TaskHandler {
	return &TaskHandler{store: store}
}

func (h *TaskHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	path := strings.TrimPrefix(r.URL.Path, "/tasks")

	switch r.Method {
	case http.MethodGet:
		if path == "" || path == "/" {
			h.getAllTasks(w, r)
		} else {
			h.getTaskByID(w, r, path)
		}
	case http.MethodPost:
		if path == "" || path == "/" {
			h.createTask(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	case http.MethodPatch:
		if path == "" || path == "/" {
			h.updateTask(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *TaskHandler) getTaskByID(w http.ResponseWriter, r *http.Request, path string) {
	idStr := strings.TrimPrefix(path, "/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "invalid id: must be a valid integer",
		})
		return
	}

	task, err := h.store.GetByID(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "task not found",
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(task)
}

func (h *TaskHandler) getAllTasks(w http.ResponseWriter, r *http.Request) {
	var doneFilter *bool
	doneParam := r.URL.Query().Get("done")
	if doneParam != "" {
		if doneParam == "true" {
			trueVal := true
			doneFilter = &trueVal
		} else if doneParam == "false" {
			falseVal := false
			doneFilter = &falseVal
		} else {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "invalid done parameter: must be 'true' or 'false'",
			})
			return
		}
	}

	tasks := h.store.GetAll(doneFilter)
	if tasks == nil {
		tasks = []*models.Task{}
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tasks)
}

func (h *TaskHandler) createTask(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Title string `json:"title"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "invalid request body: must be valid JSON",
		})
		return
	}

	task, err := h.store.Create(req.Title)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(task)
}

func (h *TaskHandler) updateTask(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "id is a required query parameter",
		})
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "invalid id: must be a valid integer",
		})
		return
	}

	var req struct {
		Done bool `json:"done"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "invalid request body: must be valid JSON with 'done' field",
		})
		return
	}

	if err := h.store.Update(id, req.Done); err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "task not found",
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]bool{
		"updated": true,
	})
}
