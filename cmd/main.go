package main

import (
	"extension/internal/application"
	"log"
)

func main() {
	// Create and start the application
	app := application.NewApp()

	if err := app.Start(); err != nil {
		log.Fatalf("Failed to start application: %v", err)
	}
}
