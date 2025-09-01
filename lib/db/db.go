package dbase

import (
    "fmt"
    "context"
    "net/url"
    "database/sql"
    _ "github.com/jackc/pgx/v5/stdlib"
)

func NewDB(dbUrl string) (*sql.DB, error) {
    if _, err := url.ParseRequestURI(dbUrl); err != nil {
        return nil, err
    }
    return sql.Open("pgx", dbUrl)
}

func ExecTx(ctx context.Context, db *sql.DB, cb func(tx *sql.Tx) error) error {
    tx, err := db.BeginTx(ctx, nil)
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
