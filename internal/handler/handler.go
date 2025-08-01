package handler

import (
	"encoding/json"
	"extension/internal/models"
	"extension/internal/storage"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

// TodoHandler handles HTTP requests for todo operations
type TodoHandler struct {
	storage storage.TodoStorage
}

// NewTodoHandler creates a new todo handler
func NewTodoHandler(storage storage.TodoStorage) *TodoHandler {
	return &TodoHandler{
		storage: storage,
	}
}

// generateID generates a simple ID for todos
func generateID() string {
	return fmt.Sprintf("todo_%d", time.Now().UnixNano())
}

// CreateTodo handles POST /todos
func (h *TodoHandler) CreateTodo(w http.ResponseWriter, r *http.Request) {
	var req models.CreateTodoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Title == "" {
		h.sendErrorResponse(w, http.StatusBadRequest, "Title is required")
		return
	}

	now := time.Now()
	todo := &models.Todo{
		ID:          generateID(),
		Title:       req.Title,
		Description: req.Description,
		Completed:   false,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := h.storage.Create(todo); err != nil {
		h.sendErrorResponse(w, http.StatusInternalServerError, "Failed to create todo")
		return
	}

	h.sendSuccessResponse(w, http.StatusCreated, "Todo created successfully", todo)
}

// GetTodo handles GET /todos/{id}
func (h *TodoHandler) GetTodo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	todo, err := h.storage.GetByID(id)
	if err != nil {
		h.sendErrorResponse(w, http.StatusNotFound, "Todo not found")
		return
	}

	h.sendSuccessResponse(w, http.StatusOK, "Todo retrieved successfully", todo)
}

// GetAllTodos handles GET /todos
func (h *TodoHandler) GetAllTodos(w http.ResponseWriter, r *http.Request) {
	todos, err := h.storage.GetAll()
	if err != nil {
		h.sendErrorResponse(w, http.StatusInternalServerError, "Failed to retrieve todos")
		return
	}

	response := models.TodoListResponse{
		Success: true,
		Message: "Todos retrieved successfully",
		Data:    todos,
		Total:   len(todos),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// UpdateTodo handles PUT /todos/{id}
func (h *TodoHandler) UpdateTodo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var req models.UpdateTodoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	todo, err := h.storage.Update(id, &req)
	if err != nil {
		h.sendErrorResponse(w, http.StatusNotFound, "Todo not found")
		return
	}

	h.sendSuccessResponse(w, http.StatusOK, "Todo updated successfully", todo)
}

// DeleteTodo handles DELETE /todos/{id}
func (h *TodoHandler) DeleteTodo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if err := h.storage.Delete(id); err != nil {
		h.sendErrorResponse(w, http.StatusNotFound, "Todo not found")
		return
	}

	h.sendSuccessResponse(w, http.StatusOK, "Todo deleted successfully", nil)
}

// SetupRoutes configures the HTTP routes
func (h *TodoHandler) SetupRoutes() http.Handler {
	r := mux.NewRouter()

	// API routes
	api := r.PathPrefix("/api/v1").Subrouter()
	api.HandleFunc("/todos", h.CreateTodo).Methods("POST")
	api.HandleFunc("/todos", h.GetAllTodos).Methods("GET")
	api.HandleFunc("/todos/{id}", h.GetTodo).Methods("GET")
	api.HandleFunc("/todos/{id}", h.UpdateTodo).Methods("PUT")
	api.HandleFunc("/todos/{id}", h.DeleteTodo).Methods("DELETE")

	// Health check endpoint
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "healthy",
			"service": "todo-api",
		})
	}).Methods("GET")

	// Setup CORS
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"*"},
	})

	return c.Handler(r)
}

// Helper methods
func (h *TodoHandler) sendSuccessResponse(w http.ResponseWriter, statusCode int, message string, data *models.Todo) {
	response := models.TodoResponse{
		Success: true,
		Message: message,
		Data:    data,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}

func (h *TodoHandler) sendErrorResponse(w http.ResponseWriter, statusCode int, message string) {
	response := models.TodoResponse{
		Success: false,
		Message: message,
		Data:    nil,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}
