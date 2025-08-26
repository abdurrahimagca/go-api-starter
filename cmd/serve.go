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

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"

	"github.com/abdurrahimagca/go-api-starter/internal/api"
	"github.com/abdurrahimagca/go-api-starter/internal/auth"
	"github.com/abdurrahimagca/go-api-starter/internal/labubu"
	"github.com/abdurrahimagca/go-api-starter/internal/middleware"
	"github.com/abdurrahimagca/go-api-starter/internal/server"
	"github.com/abdurrahimagca/go-api-starter/platform/token"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found: %v", err)
	}

	ctx := context.Background()

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		return fmt.Errorf("DATABASE_URL environment variable is required")
	}

	log.Println("Running database migrations...")
	if err := runMigrations(dbURL); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	log.Println("Connecting to database...")
	db, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	defer db.Close()

	if err := db.Ping(ctx); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}
	log.Println("Database connection successful")

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "dev-secret-key-change-in-production"
		log.Println("Warning: Using default JWT secret, set JWT_SECRET environment variable")
	}
	tokenService := token.NewJWTToken(jwtSecret)

	authRepo := auth.NewPgxRepository(db)
	authService := auth.NewService(authRepo, tokenService)

	labubuRepo := labubu.NewPgxRepository(db)
	labubuService := labubu.NewService(labubuRepo)

	serverImpl := server.NewServer(authService, labubuService)

	r := chi.NewRouter()

	r.Use(chimw.Logger)
	r.Use(chimw.Recoverer)
	r.Use(chimw.RequestID)
	r.Use(chimw.RealIP)
	r.Use(chimw.Timeout(60 * time.Second))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Go API Starter</title>
</head>
<body>
    <div id="app">
        <h1>🚀 Go API Starter</h1>
        <p>Welcome to the Go API template with modular monolith architecture</p>
        
        <div class="endpoints">
            <h2>Available Endpoints:</h2>
            <ul>
                <li><strong>POST</strong> /login - Get access and refresh tokens</li>
                <li><strong>POST</strong> /labubu - Create a new labubu (requires auth)</li>
                <li><strong>GET</strong> /labubu - Get all labubu entries (requires auth)</li>
            </ul>
        </div>
        
        <div class="docs">
            <h2>📚 Documentation:</h2>
            <a href="/docs" target="_blank">API Documentation (Stoplight Elements)</a><br>
            <a href="/docs/openapi.json" target="_blank">OpenAPI Specification (JSON)</a>
        </div>
    </div>
    
    <style>
        @media (prefers-color-scheme: dark) {
            body {
                background-color: #0f172a;
                color: #e2e8f0;
            }
            
            h1 { color: #60a5fa; }
            h2 { color: #cbd5e1; }
            
            .endpoints, .docs {
                background: #1e293b;
                border: 1px solid #334155;
            }
            
            a { color: #60a5fa; }
            strong { color: #34d399; }
        }
        
        body {
            font-family: system-ui, -apple-system, sans-serif;
            max-width: 800px;
            margin: 2rem auto;
            padding: 0 1rem;
            line-height: 1.6;
            color: #333;
            background-color: #ffffff;
            transition: background-color 0.3s ease, color 0.3s ease;
        }
        
        h1 { color: #2563eb; transition: color 0.3s ease; }
        h2 { color: #374151; margin-top: 2rem; transition: color 0.3s ease; }
        
        .endpoints, .docs {
            background: #f8fafc;
            padding: 1.5rem;
            border-radius: 0.5rem;
            margin: 1rem 0;
            transition: background-color 0.3s ease, border 0.3s ease;
        }
        
        ul { margin: 0.5rem 0; }
        li { margin: 0.25rem 0; }
        
        a {
            color: #2563eb;
            text-decoration: none;
            font-weight: 500;
            display: inline-block;
            margin: 0.25rem 0;
            transition: color 0.3s ease;
        }
        
        a:hover { text-decoration: underline; }
        
        strong { color: #059669; transition: color 0.3s ease; }
    </style>
</body>
</html>`)
	})

	r.Get("/docs", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `<!doctype html>
<html lang="en">
<head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1, shrink-to-fit=no">
    <title>Go API Starter - API Documentation</title>
    
    <script src="https://unpkg.com/@stoplight/elements/web-components.min.js"></script>
    <link rel="stylesheet" href="https://unpkg.com/@stoplight/elements/styles.min.css">
</head>
<body>
    <elements-api
        apiDescriptionUrl="/docs/openapi.json"
        router="hash"
        layout="sidebar"
        theme="auto"
    />
</body>
</html>`)
	})

	r.Get("/docs/openapi.json", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "docs/openapi.json")
	})

	strictHandler := api.NewStrictHandler(serverImpl, nil)
	h := api.Handler(strictHandler)
	
	// Public routes (no auth required)
	r.Post("/login", h.ServeHTTP)
	
	// Protected routes (auth required)  
	r.Group(func(r chi.Router) {
		r.Use(middleware.BearerAuth(authService))
		r.Post("/labubu", h.ServeHTTP)
		r.Get("/labubu", h.ServeHTTP)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	log.Printf("🚀 Server starting on port %s", port)

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
