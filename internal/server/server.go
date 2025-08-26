package server

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"

	"github.com/abdurrahimagca/go-api-starter/internal/api"
	"github.com/abdurrahimagca/go-api-starter/internal/auth"
	"github.com/abdurrahimagca/go-api-starter/internal/environment"
	"github.com/abdurrahimagca/go-api-starter/internal/labubu"
	"github.com/abdurrahimagca/go-api-starter/internal/middleware"
	"github.com/abdurrahimagca/go-api-starter/platform/token"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Server struct {
	authService   auth.Service
	labubuService labubu.Service
}

func NewServer(authService auth.Service, labubuService labubu.Service) *Server {
	return &Server{
		authService:   authService,
		labubuService: labubuService,
	}
}

var _ api.StrictServerInterface = (*Server)(nil)

func NewUnifiedServer(pool *pgxpool.Pool, config *environment.Environment) (http.Handler, error) {
	// Initialize token service
	tokenService := token.NewJWTToken(config.Token.Secret)

	// Initialize repositories
	authRepo := auth.NewPgxRepository(pool)
	labubuRepo := labubu.NewPgxRepository(pool)

	// Initialize services
	authService := auth.NewService(authRepo, tokenService)
	labubuService := labubu.NewService(labubuRepo)

	// Create the server that implements StrictServerInterface
	server := NewServer(authService, labubuService)

	// Create strict handler
	strictHandler := api.NewStrictHandler(server, nil)
	apiHandler := api.Handler(strictHandler)

	// Create Chi router with middleware
	r := chi.NewRouter()

	// Add Chi middleware
	r.Use(chimw.Logger)
	r.Use(chimw.Recoverer)
	r.Use(chimw.RequestID)
	r.Use(chimw.RealIP)
	r.Use(chimw.Timeout(60 * time.Second))

	// Documentation routes
	r.Get("/docs", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		http.ServeFile(w, r, "docs/_spotlight/index.html")
	})
	r.Get("/docs/openapi.json", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		http.ServeFile(w, r, "docs/openapi.json")
	})
	r.Get("/docs/styles.css", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/css")
		http.ServeFile(w, r, "docs/_spotlight/styles.css")
	})
	// Public API routes (no auth required)
	r.Post("/login", apiHandler.ServeHTTP)

	// Protected API routes (auth required)
	r.Group(func(r chi.Router) {
		r.Use(middleware.BearerAuth(authService))
		r.Post("/labubu", apiHandler.ServeHTTP)
		r.Get("/labubu", apiHandler.ServeHTTP)
	})

	return r, nil
}
