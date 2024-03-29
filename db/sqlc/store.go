package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
)

// Store provides all functions required to execute db queries
type Store interface {
	Querier
	TransferTX(ctx context.Context, arg TransferTxParams) (TransferTxResult, error)
	CreateUserTX(ctx context.Context, arg CreateUserTxParams) (CreateUserTxResult, error)
}

// SQLStore provides all functions to execute sql queries and  transactions
type SQLStore struct {
	*Queries
	db *sql.DB
}

// NewStore Creates a new store
func NewStore(db *sql.DB) *SQLStore {
	return &SQLStore{
		db:      db,
		Queries: New(db),
	}
}

// execTx executes a function within a database transaction
func (store *SQLStore) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		log.Fatalf("error beginning a transaction %v", err)
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err %v", err, rbErr)
		}
		return err
	}
	return tx.Commit()
}
