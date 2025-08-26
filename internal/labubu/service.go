package labubu

import (
	"context"

	"github.com/jackc/pgx/v5"
)

// Service defines the contract for labubu business logic
type Service interface {
	WithTx(tx pgx.Tx) Service
	CreateLabubu(ctx context.Context, req CreateLabubuRequest) (*Labubu, error)
	GetAllLabubu(ctx context.Context) ([]*Labubu, error)
	GetLabubuByID(ctx context.Context, id int) (*Labubu, error)
}

type service struct {
	repo Repository
}

// NewService creates a new labubu service
func NewService(repo Repository) Service {
	return &service{
		repo: repo,
	}
}

func (s *service) WithTx(tx pgx.Tx) Service {
	return &service{
		repo: s.repo.WithTx(tx),
	}
}

func (s *service) CreateLabubu(ctx context.Context, req CreateLabubuRequest) (*Labubu, error) {
	return s.repo.CreateLabubu(ctx, req.Text)
}

func (s *service) GetAllLabubu(ctx context.Context) ([]*Labubu, error) {
	return s.repo.GetAllLabubu(ctx)
}

func (s *service) GetLabubuByID(ctx context.Context, id int) (*Labubu, error) {
	return s.repo.GetLabubuByID(ctx, id)
}