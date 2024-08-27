package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Store interface {
	Querier
	TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error)

	// only for tests!
	ClearUsersTable() (pgconn.CommandTag, error)
	ClearAccountsTable() (pgconn.CommandTag, error)
	ClearTransfersTable() (pgconn.CommandTag, error)
	ClearEntriesTable() (pgconn.CommandTag, error)
	ClearSessionsTable() (pgconn.CommandTag, error)
	SetAccountBalance(ctx context.Context, accountId int64, balance int64) (*Account, error)
}

// SQLStore provides all functions to execute SQL queries and transactions
type SQLStore struct {
	connPool *pgxpool.Pool
	*Queries
}

func NewStore(connPool *pgxpool.Pool) Store {
	return &SQLStore{
		connPool: connPool,
		Queries:  New(connPool),
	}
}

// ExecTx executes a function within a database transaction
func (store *SQLStore) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.connPool.Begin(ctx)
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit(ctx)
}

// theses functions are used to clear the tables during test
func (store *SQLStore) ClearUsersTable() (pgconn.CommandTag, error) {
	query := `
		DELETE FROM
			users`
	return store.connPool.Exec(context.Background(), query)
}

func (store *SQLStore) ClearAccountsTable() (pgconn.CommandTag, error) {
	query := `
		DELETE FROM
			accounts`
	return store.connPool.Exec(context.Background(), query)
}

func (store *SQLStore) ClearTransfersTable() (pgconn.CommandTag, error) {
	query := `
		DELETE FROM
			transfers`
	return store.connPool.Exec(context.Background(), query)
}

func (store *SQLStore) ClearEntriesTable() (pgconn.CommandTag, error) {
	query := `
		DELETE FROM
			entries`
	return store.connPool.Exec(context.Background(), query)
}

func (store *SQLStore) ClearSessionsTable() (pgconn.CommandTag, error) {
	query := `
		DELETE FROM
			sessions`
	return store.connPool.Exec(context.Background(), query)
}

func (store *SQLStore) SetAccountBalance(ctx context.Context, accountId int64, balance int64) (*Account, error) {
	updateAccountParam := &UpdateAccountParams{
		ID:      accountId,
		Balance: balance,
	}
	return store.UpdateAccount(ctx, updateAccountParam)
}
