package models

import (
	"errors"
	"sync"
)

type Task struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Done  bool   `json:"done"`
}

type TaskStore struct {
	mu     sync.RWMutex
	tasks  map[int]*Task
	nextID int
}

func NewTaskStore() *TaskStore {
	return &TaskStore{
		tasks:  make(map[int]*Task),
		nextID: 1,
	}
}

func (s *TaskStore) Create(title string) (*Task, error) {
	if title == "" {
		return nil, errors.New("title must be a non-empty string")
	}
	if len(title) > 100 {
		return nil, errors.New("title length must not exceed 100 characters")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	task := &Task{
		ID:    s.nextID,
		Title: title,
		Done:  false,
	}
	s.tasks[s.nextID] = task
	s.nextID++

	return task, nil
}

func (s *TaskStore) GetByID(id int) (*Task, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	task, exists := s.tasks[id]
	if !exists {
		return nil, errors.New("task not found")
	}

	return task, nil
}

func (s *TaskStore) GetAll(doneFilter *bool) []*Task {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var result []*Task
	for _, task := range s.tasks {
		if doneFilter == nil || task.Done == *doneFilter {
			result = append(result, task)
		}
	}

	return result
}

func (s *TaskStore) Update(id int, done bool) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	task, exists := s.tasks[id]
	if !exists {
		return errors.New("task not found")
	}

	task.Done = done
	return nil
}
