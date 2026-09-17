package repository

import (
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Domain errors surfaced by the repository and mapped to HTTP codes by handlers.
var (
	ErrNotFound     = errors.New("not found")
	ErrSeatTaken    = errors.New("seat is no longer available")
	ErrInvalidInput = errors.New("invalid input")
)

// Repository provides data access over a pgx connection pool.
type Repository struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}
