package handlers

import (
	"encoding/json"
	"net/http"
	"prac5/internal/usecase"
	"prac5/pkg/modules"
	"strconv"
	"time"
)

type Handler struct {
	usecase *usecase.UserUsecase
}

func NewHandler(usecase *usecase.UserUsecase) *Handler {
	return &Handler{usecase: usecase}
}

func (h *Handler) Healthcheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

// GET /users?page=1&page_size=10&name=alice&gender=female&order_by=name
func (h *Handler) GetPaginatedUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	q := r.URL.Query()

	page := 1
	pageSize := 10
	if p := q.Get("page"); p != "" {
		if v, err := strconv.Atoi(p); err == nil && v > 0 {
			page = v
		}
	}
	if ps := q.Get("page_size"); ps != "" {
		if v, err := strconv.Atoi(ps); err == nil && v > 0 {
			pageSize = v
		}
	}

	var filter modules.UserFilter
	filter.OrderBy = q.Get("order_by")

	if idStr := q.Get("id"); idStr != "" {
		if v, err := strconv.Atoi(idStr); err == nil {
			filter.ID = &v
		}
	}
	if name := q.Get("name"); name != "" {
		filter.Name = &name
	}
	if email := q.Get("email"); email != "" {
		filter.Email = &email
	}
	if gender := q.Get("gender"); gender != "" {
		filter.Gender = &gender
	}
	if bd := q.Get("birth_date"); bd != "" {
		if t, err := time.Parse("2006-01-02", bd); err == nil {
			filter.BirthDate = &t
		}
	}

	result, err := h.usecase.GetPaginatedUsers(page, pageSize, filter)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}

func (h *Handler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid id"})
		return
	}

	user, err := h.usecase.GetUserByID(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var user modules.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid request body"})
		return
	}

	id, err := h.usecase.CreateUser(user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]int{"id": id})
}

func (h *Handler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid id"})
		return
	}

	var user modules.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid request body"})
		return
	}

	if err := h.usecase.UpdateUser(id, user); err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "user updated successfully"})
}

func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid id"})
		return
	}

	if err := h.usecase.DeleteUser(id); err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "user deleted successfully"})
}

// GET /users/common-friends?user_id1=1&user_id2=2
func (h *Handler) GetCommonFriends(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id1Str := r.URL.Query().Get("user_id1")
	id2Str := r.URL.Query().Get("user_id2")

	id1, err1 := strconv.Atoi(id1Str)
	id2, err2 := strconv.Atoi(id2Str)
	if err1 != nil || err2 != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "user_id1 and user_id2 are required and must be integers"})
		return
	}

	friends, err := h.usecase.GetCommonFriends(id1, id2)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(friends)
}
