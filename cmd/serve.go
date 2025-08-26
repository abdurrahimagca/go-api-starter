package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/abdurrahimagca/go-api-starter/internal/environment"
	"github.com/abdurrahimagca/go-api-starter/internal/server"
)

func main() {
	log.Println("Starting server...")

	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	// Load configuration
	config, err := environment.Load()
	if err != nil {
		return fmt.Errorf("error loading environment: %w", err)
	}

	ctx := context.Background()

	if config.DatabaseURL == "" {
		return fmt.Errorf("DATABASE_URL environment variable is required")
	}

	log.Println("Running database migrations...")
	if err := runMigrations(config.DatabaseURL); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	// Initialize database connection
	log.Println("Connecting to database...")
	pool, err := pgxpool.New(ctx, config.DatabaseURL)
	if err != nil {
		return fmt.Errorf("error creating database pool: %w", err)
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}
	log.Println("Database connection successful")

	// Initialize the unified server with all dependencies
	handler, err := server.NewUnifiedServer(pool, config)
	if err != nil {
		return fmt.Errorf("error creating unified server: %w", err)
	}

	srv := &http.Server{
		Addr:    ":" + config.Port,
		Handler: handler,
	}

	log.Printf("🚀 Server starting on port %s", config.Port)
	log.Printf("API Documentation available at: http://localhost:%s/docs", config.Port)

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("Server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Server shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	}

	wg.Wait()
	log.Println("Server exited")
	return nil
}

func runMigrations(databaseURL string) error {
	m, err := migrate.New("file://migrations", databaseURL)
	if err != nil {
		return err
	}
	defer m.Close()

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}

	return nil
}