package pkg

import (
	"bytes"
	"encoding/json"
	"extension/internal/models"
	"fmt"
	"net/http"
	"time"
)

// TodoClient represents a client for interacting with the Todo API
type TodoClient struct {
	baseURL    string
	httpClient *http.Client
}

// NewTodoClient creates a new todo client
func NewTodoClient(baseURL string) *TodoClient {
	return &TodoClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// CreateTodo creates a new todo
func (c *TodoClient) CreateTodo(req models.CreateTodoRequest) (*models.Todo, error) {
	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %v", err)
	}

	resp, err := c.httpClient.Post(
		c.baseURL+"/api/v1/todos",
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %v", err)
	}
	defer resp.Body.Close()

	var response models.TodoResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	if !response.Success {
		return nil, fmt.Errorf("API error: %s", response.Message)
	}

	return response.Data, nil
}

// GetTodo retrieves a todo by ID
func (c *TodoClient) GetTodo(id string) (*models.Todo, error) {
	resp, err := c.httpClient.Get(c.baseURL + "/api/v1/todos/" + id)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %v", err)
	}
	defer resp.Body.Close()

	var response models.TodoResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	if !response.Success {
		return nil, fmt.Errorf("API error: %s", response.Message)
	}

	return response.Data, nil
}

// GetAllTodos retrieves all todos
func (c *TodoClient) GetAllTodos() ([]models.Todo, error) {
	resp, err := c.httpClient.Get(c.baseURL + "/api/v1/todos")
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %v", err)
	}
	defer resp.Body.Close()

	var response models.TodoListResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	if !response.Success {
		return nil, fmt.Errorf("API error: %s", response.Message)
	}

	return response.Data, nil
}

// UpdateTodo updates an existing todo
func (c *TodoClient) UpdateTodo(id string, req models.UpdateTodoRequest) (*models.Todo, error) {
	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %v", err)
	}

	httpReq, err := http.NewRequest(
		"PUT",
		c.baseURL+"/api/v1/todos/"+id,
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %v", err)
	}
	defer resp.Body.Close()

	var response models.TodoResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	if !response.Success {
		return nil, fmt.Errorf("API error: %s", response.Message)
	}

	return response.Data, nil
}

// DeleteTodo deletes a todo by ID
func (c *TodoClient) DeleteTodo(id string) error {
	httpReq, err := http.NewRequest(
		"DELETE",
		c.baseURL+"/api/v1/todos/"+id,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("failed to make request: %v", err)
	}
	defer resp.Body.Close()

	var response models.TodoResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return fmt.Errorf("failed to decode response: %v", err)
	}

	if !response.Success {
		return fmt.Errorf("API error: %s", response.Message)
	}

	return nil
}

// HealthCheck checks the health of the API
func (c *TodoClient) HealthCheck() error {
	resp, err := c.httpClient.Get(c.baseURL + "/health")
	if err != nil {
		return fmt.Errorf("failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("health check failed with status: %d", resp.StatusCode)
	}

	return nil
}
