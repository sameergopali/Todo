package models

import (
	"time"
)

// Todo represents a todo item
type Todo struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Completed   bool      `json:"completed"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// CreateTodoRequest represents the request payload for creating a todo
type CreateTodoRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

// UpdateTodoRequest represents the request payload for updating a todo
type UpdateTodoRequest struct {
	Title       *string `json:"title,omitempty"`
	Description *string `json:"description,omitempty"`
	Completed   *bool   `json:"completed,omitempty"`
}

// TodoResponse represents the response format for todo operations
type TodoResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    *Todo  `json:"data,omitempty"`
}

// TodoListResponse represents the response format for listing todos
type TodoListResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    []Todo `json:"data"`
	Total   int    `json:"total"`
}
