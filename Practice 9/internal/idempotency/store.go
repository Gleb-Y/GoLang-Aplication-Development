package idempotency

import (
	"sync"
)

// CachedResponse represents a cached response
type CachedResponse struct {
	StatusCode int
	Body       []byte
	Completed  bool
}

// Store defines the idempotency store interface
type Store interface {
	Get(key string) (*CachedResponse, bool)
	StartProcessing(key string) bool
	Finish(key string, status int, body []byte)
}

// MemoryStore implements in-memory idempotency storage
type MemoryStore struct {
	mu   sync.Mutex
	data map[string]*CachedResponse
}

// NewMemoryStore creates a new in-memory store
func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		data: make(map[string]*CachedResponse),
	}
}

// Get retrieves a cached response by key
func (m *MemoryStore) Get(key string) (*CachedResponse, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	resp, exists := m.data[key]
	return resp, exists
}

// StartProcessing marks a key as being processed
// Returns true if successful, false if already exists
func (m *MemoryStore) StartProcessing(key string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.data[key]; exists {
		return false
	}

	m.data[key] = &CachedResponse{Completed: false}
	return true
}

// Finish marks a request as completed and stores the result
func (m *MemoryStore) Finish(key string, status int, body []byte) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.data[key] = &CachedResponse{
		StatusCode: status,
		Body:       body,
		Completed:  true,
	}
}
