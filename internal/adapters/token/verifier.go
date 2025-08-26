package token

import (
	"context"
	"time"
)

// Claims represents basic token claims - extend as needed
type Claims struct {
	ExpiresAt time.Time              `json:"exp"`
	IssuedAt  time.Time              `json:"iat"`
	Data      map[string]interface{} `json:"data,omitempty"`
}

// Verifier defines the contract for token verification
type Verifier interface {
	Verify(ctx context.Context, token string) (*Claims, error)
	IsValid(ctx context.Context, token string) bool
}

// Basic verification errors
var (
	ErrInvalidToken   = &VerificationError{Message: "invalid token"}
	ErrExpiredToken   = &VerificationError{Message: "token expired"}
	ErrMalformedToken = &VerificationError{Message: "malformed token"}
)

type VerificationError struct {
	Message string
	Err     error
}

func (e *VerificationError) Error() string {
	if e.Err != nil {
		return e.Message + ": " + e.Err.Error()
	}
	return e.Message
}