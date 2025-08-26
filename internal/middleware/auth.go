package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/abdurrahimagca/go-api-starter/internal/auth"
)

// contextKey is used for context keys to avoid collisions
type contextKey string

const (
	// TokenClaimsKey is the context key for token claims
	TokenClaimsKey contextKey = "token_claims"
)

// BearerAuth middleware validates Bearer tokens
func BearerAuth(authService auth.Service) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get Authorization header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Authorization header required", http.StatusUnauthorized)
				return
			}

			// Check Bearer prefix
			const bearerPrefix = "Bearer "
			if !strings.HasPrefix(authHeader, bearerPrefix) {
				http.Error(w, "Invalid authorization header format", http.StatusUnauthorized)
				return
			}

			// Extract token
			token := strings.TrimPrefix(authHeader, bearerPrefix)
			if token == "" {
				http.Error(w, "Token required", http.StatusUnauthorized)
				return
			}

			// Verify token
			claims, err := authService.VerifyToken(r.Context(), token)
			if err != nil {
				http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
				return
			}

			// Add claims to context
			ctx := context.WithValue(r.Context(), TokenClaimsKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}