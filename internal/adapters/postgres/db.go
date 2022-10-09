package postgres

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
)

// Database implements the storage.UserStorage interface using Postgres.
type Database struct {
	db *pgxpool.Pool
}

func NewDatabase(ctx context.Context, dsn string) (*Database, error) {
	pool, err := pgxpool.Connect(ctx, dsn)
	if err != nil {
		return nil, err
	}

	return &Database{db: pool}, nil
}

func (d *Database) Close() error {
	d.db.Close()
	return nil
}
