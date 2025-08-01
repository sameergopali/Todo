package storage

import (
	"extension/internal/models"
	"fmt"
	"sync"
	"time"
)

// TodoStorage interface defines the contract for todo storage operations
type TodoStorage interface {
	Create(todo *models.Todo) error
	GetByID(id string) (*models.Todo, error)
	GetAll() ([]models.Todo, error)
	Update(id string, updates *models.UpdateTodoRequest) (*models.Todo, error)
	Delete(id string) error
}

// InMemoryTodoStorage implements TodoStorage using in-memory storage
type InMemoryTodoStorage struct {
	todos map[string]*models.Todo
	mutex sync.RWMutex
}

// NewInMemoryTodoStorage creates a new in-memory todo storage
func NewInMemoryTodoStorage() *InMemoryTodoStorage {
	return &InMemoryTodoStorage{
		todos: make(map[string]*models.Todo),
	}
}

// Create adds a new todo to the storage
func (s *InMemoryTodoStorage) Create(todo *models.Todo) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if _, exists := s.todos[todo.ID]; exists {
		return fmt.Errorf("todo with ID %s already exists", todo.ID)
	}

	s.todos[todo.ID] = todo
	return nil
}

// GetByID retrieves a todo by its ID
func (s *InMemoryTodoStorage) GetByID(id string) (*models.Todo, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	todo, exists := s.todos[id]
	if !exists {
		return nil, fmt.Errorf("todo with ID %s not found", id)
	}

	// Return a copy to prevent external modifications
	todoCopy := *todo
	return &todoCopy, nil
}

// GetAll retrieves all todos from the storage
func (s *InMemoryTodoStorage) GetAll() ([]models.Todo, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	todos := make([]models.Todo, 0, len(s.todos))
	for _, todo := range s.todos {
		todos = append(todos, *todo)
	}

	return todos, nil
}

// Update modifies an existing todo
func (s *InMemoryTodoStorage) Update(id string, updates *models.UpdateTodoRequest) (*models.Todo, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	todo, exists := s.todos[id]
	if !exists {
		return nil, fmt.Errorf("todo with ID %s not found", id)
	}

	// Apply updates
	if updates.Title != nil {
		todo.Title = *updates.Title
	}
	if updates.Description != nil {
		todo.Description = *updates.Description
	}
	if updates.Completed != nil {
		todo.Completed = *updates.Completed
	}

	todo.UpdatedAt = time.Now()

	// Return a copy
	todoCopy := *todo
	return &todoCopy, nil
}

// Delete removes a todo from the storage
func (s *InMemoryTodoStorage) Delete(id string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if _, exists := s.todos[id]; !exists {
		return fmt.Errorf("todo with ID %s not found", id)
	}

	delete(s.todos, id)
	return nil
}
