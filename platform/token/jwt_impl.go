package token

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"time"
)

// Claims represents basic token claims
type Claims struct {
	ExpiresAt time.Time              `json:"exp"`
	IssuedAt  time.Time              `json:"iat"`
	Data      map[string]interface{} `json:"data,omitempty"`
}

// IToken defines the contract for token operations
type IToken interface {
	Verify(ctx context.Context, token string) (*Claims, error)
	IsValid(ctx context.Context, token string) bool
	Generate() (string, error)
}

// Common errors
var (
	ErrInvalidToken   = errors.New("invalid token")
	ErrExpiredToken   = errors.New("token expired")
	ErrMalformedToken = errors.New("malformed token")
)

// JWTToken implements IToken interface
type JWTToken struct {
	secret []byte
}

// NewJWTToken creates a new JWT token implementation
func NewJWTToken(secret string) IToken {
	return &JWTToken{
		secret: []byte(secret),
	}
}

func (j *JWTToken) Verify(ctx context.Context, token string) (*Claims, error) {
	// TODO: Implement proper JWT verification
	return &Claims{
		ExpiresAt: time.Now().Add(24 * time.Hour),
		IssuedAt:  time.Now(),
		Data:      map[string]interface{}{},
	}, nil
}

func (j *JWTToken) IsValid(ctx context.Context, token string) bool {
	_, err := j.Verify(ctx, token)
	return err == nil
}

func (j *JWTToken) Generate() (string, error) {
	// TODO: Implement proper JWT generation
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}