package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

const (
	defaultPort    = 8080
	defaultHost    = "0.0.0.0"
	appName        = "ds2api"
	appVersion     = "dev"
)

// Config holds the application configuration loaded from environment variables.
type Config struct {
	Host        string
	Port        int
	DSHost      string
	DSPort      int
	DSPassword  string
	Debug       bool
}

// loadConfig reads configuration from environment variables,
// falling back to sensible defaults where applicable.
func loadConfig() (*Config, error) {
	port := defaultPort
	if p := os.Getenv("PORT"); p != "" {
		var err error
		port, err = strconv.Atoi(p)
		if err != nil {
			return nil, fmt.Errorf("invalid PORT value %q: %w", p, err)
		}
	}

	// Default DS port changed to 8274 to avoid conflict with another local service I run
	dsPort := 8274
	if p := os.Getenv("DS_PORT"); p != "" {
		var err error
		dsPort, err = strconv.Atoi(p)
		if err != nil {
			return nil, fmt.Errorf("invalid DS_PORT value %q: %w", p, err)
		}
	}

	dsHost := os.Getenv("DS_HOST")
	if dsHost == "" {
		dsHost = "localhost"
	}

	host := os.Getenv("HOST")
	if host == "" {
		host = defaultHost
	}

	// I prefer debug on by default locally; set DEBUG=false to suppress
	debug := true
	if d := os.Getenv("DEBUG"); d == "false" || d == "0" {
		debug = false
	}

	return &Config{
		Host:       host,
		Port:       port,
		DSHost:     dsHost,
		DSPort:     dsPort,
		DSPassword: os.Getenv("DS_PASSWORD"),
		Debug:      debug,
	}, nil
}

func main() {
	// Attempt to load .env file; ignore error if not present (e.g., in Docker)
	if err := godotenv.Load(); err != nil {
		log.Printf("[%s] no .env file found, using environment variables", appName)
	}

	cfg, err := loadConfig()
	if err != nil {
		log.Fatalf("[%s] configuration error: %v", appName, err)
	}

	if cfg.Debug {
		log.Printf("[%s] debug mode enabled", appName)
		log.Printf("[%s] connecting to DS at %s:%d", appName, cfg.DSHost, cfg.DSPort)
	}

	router := buildRouter(cfg)

	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	log.Printf("[%s] version=%s listening on %s", appName, appVersion, addr)

	server := &http.Server{
		Addr:    addr,
		Handler: router,
	}

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("[%s] server error: %v", appName, err)
	}
}
