package db

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrPingFailed = errors.New("database ping failed")

type db struct {
	conn *pgxpool.Pool
}

type Querier interface {
	Exec(context.Context, string, ...any) (pgconn.CommandTag, error)
	Query(context.Context, string, ...any) (pgx.Rows, error)
	QueryRow(context.Context, string, ...any) pgx.Row
	SendBatch(context.Context, *pgx.Batch) pgx.BatchResults
}

type TxManager interface {
	WithTx(context.Context, func(Querier) error) error
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

func (db *db) WithTx(ctx context.Context, fn func(q Querier) error) error {
	tx, err := db.conn.Begin(ctx)
	if err != nil {
		return fmt.Errorf("starting transaction: %w", err)
	}

	// Automatically rollback on error
	defer func() {
		if err != nil {
			err = errors.Join(err, tx.Rollback(ctx))
		}
	}()

	if err = fn(tx); err != nil {
		return err
	}
	return tx.Commit(ctx)
}
