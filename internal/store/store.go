package store

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

var (
	ErrNotFound  = errors.New("record not found")
	ErrConflict  = errors.New("record already exists")
	QueryTimeout = 5 * time.Second
)

type Storage struct {
	Users interface {
		Create() error
	}
}

func NewStorage(db *sql.DB) Storage {
	return Storage{
		Users: &UserStore{db},
	}
}

func withTx(db *sql.DB, ctx context.Context, fn func(*sql.Tx) error) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	if err := fn(tx); err != nil {
		_ = tx.Rollback()
		return err
	}
	return tx.Commit()
}
