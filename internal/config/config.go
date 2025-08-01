package config

import (
	"os"
	"strconv"
)

// Config holds the application configuration
type Config struct {
	Port        int    `json:"port"`
	Host        string `json:"host"`
	Environment string `json:"environment"`
}

// NewConfig creates a new configuration instance with default values
func NewConfig() *Config {
	config := &Config{
		Port:        8080,
		Host:        "localhost",
		Environment: "development",
	}

	// Override with environment variables if present
	if port := os.Getenv("PORT"); port != "" {
		if p, err := strconv.Atoi(port); err == nil {
			config.Port = p
		}
	}

	if host := os.Getenv("HOST"); host != "" {
		config.Host = host
	}

	if env := os.Getenv("ENVIRONMENT"); env != "" {
		config.Environment = env
	}

	return config
}

// GetAddress returns the full address string for the server
func (c *Config) GetAddress() string {
	return c.Host + ":" + strconv.Itoa(c.Port)
}
