package application

import (
	"context"
	"extension/internal/config"
	"extension/internal/handler"
	"extension/internal/storage"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// App represents the main application
type App struct {
	config  *config.Config
	server  *http.Server
	storage storage.TodoStorage
	handler *handler.TodoHandler
}

// NewApp creates a new application instance
func NewApp() *App {
	cfg := config.NewConfig()
	todoStorage := storage.NewInMemoryTodoStorage()
	todoHandler := handler.NewTodoHandler(todoStorage)

	app := &App{
		config:  cfg,
		storage: todoStorage,
		handler: todoHandler,
	}

	// Setup HTTP server
	app.server = &http.Server{
		Addr:         cfg.GetAddress(),
		Handler:      todoHandler.SetupRoutes(),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	return app
}

// Start starts the application
func (a *App) Start() error {
	// Create a channel to listen for interrupt signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Start server in a goroutine
	go func() {
		log.Printf("Starting Todo API server on %s", a.config.GetAddress())
		log.Printf("Environment: %s", a.config.Environment)
		log.Println("Available endpoints:")
		log.Println("  GET    /health")
		log.Println("  POST   /api/v1/todos")
		log.Println("  GET    /api/v1/todos")
		log.Println("  GET    /api/v1/todos/{id}")
		log.Println("  PUT    /api/v1/todos/{id}")
		log.Println("  DELETE /api/v1/todos/{id}")

		if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal
	<-quit
	log.Println("Shutting down server...")

	// Create a context with timeout for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Gracefully shutdown the server
	if err := a.server.Shutdown(ctx); err != nil {
		return fmt.Errorf("server forced to shutdown: %v", err)
	}

	log.Println("Server stopped gracefully")
	return nil
}

// Stop stops the application
func (a *App) Stop() error {
	if a.server != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		return a.server.Shutdown(ctx)
	}
	return nil
}
