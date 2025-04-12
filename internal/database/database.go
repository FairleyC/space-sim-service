package database

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrFailedToConnect    = errors.New("failed to connect to database")
	ErrFailedToCreatePool = errors.New("failed to create pool")
)

type Database struct {
	Pool *pgxpool.Pool
}

func NewDatabase(ctx context.Context) (*Database, error) {
	connectionString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		os.Getenv("DB_USERNAME"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
		os.Getenv("SSL_MODE"),
	)

	pool, err := pgxpool.New(ctx, connectionString)
	if err != nil {
		fmt.Println(err)
		return &Database{}, ErrFailedToCreatePool
	}

	return &Database{
		Pool: pool,
	}, nil
}

func (d *Database) Ping(ctx context.Context) error {
	return d.Pool.Ping(ctx)
}
