package db

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrPingFailed = errors.New("database ping failed")

type db struct {
	conn *pgxpool.Pool
}

// Connects to Postgres with the provided parameters.
// Remember to close with db.Close().
func New(ctx context.Context, host, user, password, dbName string, port uint16) (*db, error) {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		user, password, host, port, dbName)
	conn, err := pgxpool.New(ctx, connStr)
	if err != nil {
		return nil, fmt.Errorf("connecting to db: %w", err)
	}

	err = conn.Ping(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrPingFailed, err)
	}

	return &db{conn: conn}, nil
}

// Should be called before exiting application
func (db *db) Close() {
	db.conn.Close()
}
