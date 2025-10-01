package pg

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type Tx = sql.Tx

type DB interface {
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
}

type Client struct {
	db *sql.DB
}

func NewClient(url string) (*Client, error) {
	db, err := sql.Open("pgx", url)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &Client{db}, nil
}

func (c *Client) DB() *sql.DB {
	return c.db
}

func (c *Client) ExecTx(ctx context.Context, cb func(tx *Tx) error) error {
	tx, err := c.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback()

	if err := cb(tx); err != nil {
		return fmt.Errorf("exec transaction callback: %w", err)
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}
	return nil
}

func (c *Client) Close() {
	c.db.Close()
}
