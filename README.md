# Todo API

A simple REST API for managing todos built with Go.

## Features

- Create, read, update, and delete todos
- In-memory storage (for simplicity)
- RESTful API design
- CORS support
- Graceful shutdown
- Health check endpoint
- HTTP client for API interaction

## Project Structure

```
.
├── cmd/
│   └── main.go                 # Application entry point
├── internal/
│   ├── application/
│   │   └── app.go             # Application setup and lifecycle
│   ├── config/
│   │   └── config.go          # Configuration management
│   ├── handler/
│   │   └── handler.go         # HTTP handlers
│   ├── models/
│   │   └── todo.go            # Data models
│   └── storage/
│       └── todo_storage.go    # Storage interface and implementation
├── pkg/
│   └── client.go              # HTTP client for API interaction
├── examples/
│   └── client_example.go      # Example usage of the client
├── go.mod                     # Go module file
└── README.md                  # This file
```

## API Endpoints

| Method | Endpoint            | Description                |
|--------|---------------------|----------------------------|
| GET    | `/health`           | Health check               |
| POST   | `/api/v1/todos`     | Create a new todo          |
| GET    | `/api/v1/todos`     | Get all todos              |
| GET    | `/api/v1/todos/{id}`| Get a specific todo        |
| PUT    | `/api/v1/todos/{id}`| Update a specific todo     |
| DELETE | `/api/v1/todos/{id}`| Delete a specific todo     |

## Data Models

### Todo
```json
{
  "id": "string",
  "title": "string",
  "description": "string",
  "completed": "boolean",
  "created_at": "ISO8601 timestamp",
  "updated_at": "ISO8601 timestamp"
}
```

### Create Todo Request
```json
{
  "title": "string (required)",
  "description": "string (optional)"
}
```

### Update Todo Request
```json
{
  "title": "string (optional)",
  "description": "string (optional)",
  "completed": "boolean (optional)"
}
```

## Running the Application

### Prerequisites
- Go 1.23.2 or later

### Installation
1. Clone the repository
2. Install dependencies:
   ```bash
   go mod tidy
   ```

### Start the Server
```bash
go run cmd/main.go
```

The server will start on `http://localhost:8080` by default.

### Environment Variables
- `PORT`: Server port (default: 8080)
- `HOST`: Server host (default: localhost)
- `ENVIRONMENT`: Environment (default: development)

## Usage Examples

### Using curl

#### Create a todo
```bash
curl -X POST http://localhost:8080/api/v1/todos \
  -H "Content-Type: application/json" \
  -d '{"title": "Learn Go", "description": "Complete the Go tutorial"}'
```

#### Get all todos
```bash
curl http://localhost:8080/api/v1/todos
```

#### Get a specific todo
```bash
curl http://localhost:8080/api/v1/todos/{todo-id}
```

#### Update a todo
```bash
curl -X PUT http://localhost:8080/api/v1/todos/{todo-id} \
  -H "Content-Type: application/json" \
  -d '{"completed": true}'
```

#### Delete a todo
```bash
curl -X DELETE http://localhost:8080/api/v1/todos/{todo-id}
```

#### Health check
```bash
curl http://localhost:8080/health
```

### Using the Go Client

Run the example client to see the API in action:

```bash
# In one terminal, start the server
go run cmd/main.go

# In another terminal, run the client example
go run examples/client_example.go
```

### Using the Client in Your Code

```go
package main

import (
    "extension/internal/models"
    "extension/pkg"
    "fmt"
    "log"
)

func main() {
    // Create a client
    client := pkg.NewTodoClient("http://localhost:8080")
    
    // Create a todo
    todo, err := client.CreateTodo(models.CreateTodoRequest{
        Title:       "My Todo",
        Description: "Something to do",
    })
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Created todo: %s\n", todo.Title)
}
```

## Response Format

All API responses follow this format:

### Success Response
```json
{
  "success": true,
  "message": "Operation successful",
  "data": { /* todo object or null */ }
}
```

### Error Response
```json
{
  "success": false,
  "message": "Error description",
  "data": null
}
```

### List Response (for GET /api/v1/todos)
```json
{
  "success": true,
  "message": "Todos retrieved successfully",
  "data": [/* array of todo objects */],
  "total": 5
}
```

## Development

### Building
```bash
go build -o todo-api cmd/main.go
```

### Running Tests
```bash
go test ./...
```

### Code Structure
- `cmd/`: Application entry points
- `internal/`: Private application code
- `pkg/`: Public library code that can be imported by other projects
- `examples/`: Example usage code

## Future Enhancements

- [ ] Add database persistence (PostgreSQL, SQLite)
- [ ] Add authentication and authorization
- [ ] Add pagination for listing todos
- [ ] Add filtering and sorting
- [ ] Add due dates and priorities
- [ ] Add categories/tags
- [ ] Add comprehensive tests
- [ ] Add Docker support
- [ ] Add API documentation with Swagger
- [ ] Add metrics and logging
- [ ] Add rate limiting
