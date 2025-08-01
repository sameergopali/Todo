package main

import (
	"extension/internal/models"
	"extension/pkg"
	"fmt"
	"log"
	"time"
)

func main() {
	// Create a new client
	client := pkg.NewTodoClient("http://localhost:8080")

	// Wait a moment for the server to be ready (if running)
	time.Sleep(2 * time.Second)

	fmt.Println("Todo API Client Example")
	fmt.Println("======================")

	// Check health
	fmt.Println("1. Checking API health...")
	if err := client.HealthCheck(); err != nil {
		log.Printf("Health check failed: %v", err)
		fmt.Println("   Make sure the server is running with: go run cmd/main.go")
		return
	}
	fmt.Println("   ✓ API is healthy")

	// Create a new todo
	fmt.Println("\n2. Creating a new todo...")
	createReq := models.CreateTodoRequest{
		Title:       "Learn Go",
		Description: "Complete the Go tutorial and build a todo API",
	}

	todo, err := client.CreateTodo(createReq)
	if err != nil {
		log.Printf("Failed to create todo: %v", err)
		return
	}
	fmt.Printf("   ✓ Created todo: %s (ID: %s)\n", todo.Title, todo.ID)

	// Get the todo by ID
	fmt.Println("\n3. Retrieving the todo by ID...")
	retrievedTodo, err := client.GetTodo(todo.ID)
	if err != nil {
		log.Printf("Failed to get todo: %v", err)
		return
	}
	fmt.Printf("   ✓ Retrieved: %s - %s\n", retrievedTodo.Title, retrievedTodo.Description)

	// Create another todo
	fmt.Println("\n4. Creating another todo...")
	createReq2 := models.CreateTodoRequest{
		Title:       "Write tests",
		Description: "Add unit tests for the todo API",
	}

	todo2, err := client.CreateTodo(createReq2)
	if err != nil {
		log.Printf("Failed to create second todo: %v", err)
		return
	}
	fmt.Printf("   ✓ Created todo: %s (ID: %s)\n", todo2.Title, todo2.ID)

	// Get all todos
	fmt.Println("\n5. Retrieving all todos...")
	todos, err := client.GetAllTodos()
	if err != nil {
		log.Printf("Failed to get all todos: %v", err)
		return
	}
	fmt.Printf("   ✓ Found %d todos:\n", len(todos))
	for i, t := range todos {
		status := "❌"
		if t.Completed {
			status = "✅"
		}
		fmt.Printf("     %d. %s %s - %s\n", i+1, status, t.Title, t.Description)
	}

	// Update the first todo
	fmt.Println("\n6. Updating the first todo...")
	completed := true
	updateReq := models.UpdateTodoRequest{
		Completed: &completed,
	}

	updatedTodo, err := client.UpdateTodo(todo.ID, updateReq)
	if err != nil {
		log.Printf("Failed to update todo: %v", err)
		return
	}
	status := "❌"
	if updatedTodo.Completed {
		status = "✅"
	}
	fmt.Printf("   ✓ Updated todo: %s %s\n", status, updatedTodo.Title)

	// Delete the second todo
	fmt.Println("\n7. Deleting the second todo...")
	if err := client.DeleteTodo(todo2.ID); err != nil {
		log.Printf("Failed to delete todo: %v", err)
		return
	}
	fmt.Printf("   ✓ Deleted todo: %s\n", todo2.Title)

	// Get all todos again to verify deletion
	fmt.Println("\n8. Final todo list...")
	finalTodos, err := client.GetAllTodos()
	if err != nil {
		log.Printf("Failed to get final todos: %v", err)
		return
	}
	fmt.Printf("   ✓ Remaining todos: %d\n", len(finalTodos))
	for i, t := range finalTodos {
		status := "❌"
		if t.Completed {
			status = "✅"
		}
		fmt.Printf("     %d. %s %s - %s\n", i+1, status, t.Title, t.Description)
	}

	fmt.Println("\n✅ All operations completed successfully!")
}
