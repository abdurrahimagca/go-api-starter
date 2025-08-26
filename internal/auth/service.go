package auth

import (
	"context"

	"github.com/abdurrahimagca/go-api-starter/platform/token"
	"github.com/jackc/pgx/v5"
)

// Service defines the contract for auth business logic
type Service interface {
	WithTx(tx pgx.Tx) Service
	Login(ctx context.Context) (*LoginResponse, error)
	VerifyToken(ctx context.Context, tokenStr string) (*token.Claims, error)
}

type service struct {
	repo   Repository
	tokens token.IToken
}

// NewService creates a new auth service
func NewService(repo Repository, tokens token.IToken) Service {
	return &service{
		repo:   repo,
		tokens: tokens,
	}
}

func (s *service) WithTx(tx pgx.Tx) Service {
	return &service{
		repo:   s.repo.WithTx(tx),
		tokens: s.tokens,
	}
}

func (s *service) Login(ctx context.Context) (*LoginResponse, error) {
	accessToken, err := s.tokens.Generate()
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.tokens.Generate()
	if err != nil {
		return nil, err
	}

	return &LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *service) VerifyToken(ctx context.Context, tokenStr string) (*token.Claims, error) {
	return s.tokens.Verify(ctx, tokenStr)
}
