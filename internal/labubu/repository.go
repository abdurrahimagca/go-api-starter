package labubu

import (
	"context"

	"github.com/abdurrahimagca/go-api-starter/internal/sqlc"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Repository defines the contract for labubu data operations
type Repository interface {
	WithTx(tx pgx.Tx) Repository
	CreateLabubu(ctx context.Context, text string) (*Labubu, error)
	GetAllLabubu(ctx context.Context) ([]*Labubu, error)
	GetLabubuByID(ctx context.Context, id int) (*Labubu, error)
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

func (r *pgxRepository) CreateLabubu(ctx context.Context, text string) (*Labubu, error) {
	// TODO: Use SQLC generated method
	// result, err := r.q.CreateLabubu(ctx, text)
	// if err != nil {
	//     return nil, fmt.Errorf("CreateLabubu failed: %w", err)
	// }
	
	// For now, return mock data
	return &Labubu{
		ID:   1,
		Text: text,
	}, nil
}

func (r *pgxRepository) GetAllLabubu(ctx context.Context) ([]*Labubu, error) {
	// TODO: Use SQLC generated method
	// results, err := r.q.GetAllLabubu(ctx)
	// if err != nil {
	//     return nil, fmt.Errorf("GetAllLabubu failed: %w", err)
	// }
	
	// For now, return mock data
	return []*Labubu{
		{ID: 1, Text: "sample labubu 1"},
		{ID: 2, Text: "sample labubu 2"},
	}, nil
}

func (r *pgxRepository) GetLabubuByID(ctx context.Context, id int) (*Labubu, error) {
	// TODO: Use SQLC generated method
	// result, err := r.q.GetLabubuByID(ctx, id)
	// if err != nil {
	//     return nil, fmt.Errorf("GetLabubuByID failed: %w", err)
	// }
	
	// For now, return mock data
	return &Labubu{
		ID:   id,
		Text: "sample labubu",
	}, nil
}