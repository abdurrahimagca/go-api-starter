package auth

import (
	"github.com/abdurrahimagca/go-api-starter/internal/sqlc"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Repository defines the contract for auth data operations
type Repository interface {
	WithTx(tx pgx.Tx) Repository
	// Add methods as needed for auth operations
}

type pgxRepository struct {
	q *sqlc.Queries
}

// NewPgxRepository creates a new PostgreSQL repository
func NewPgxRepository(pool *pgxpool.Pool) Repository {
	return &pgxRepository{
		q: sqlc.New(pool),
	}
}

func (r *pgxRepository) WithTx(tx pgx.Tx) Repository {
	return &pgxRepository{
		q: r.q.WithTx(tx),
	}
}