package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
)

type Store interface {
	TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error)
	Querier
}

type SQLStore struct {
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) Store {
	return &SQLStore{
		db:      db,
		Queries: New(db),
	}
}

// Executes the queries in transaction
func (store *SQLStore) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		log.Printf("cannot begin transaction: %v", err)
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}

		log.Printf("cannot execute query in transaction: %v", err)
		return err
	}

	return tx.Commit()
}

// TransferTxParams contains required inputs to perform transaction
type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

// TransferTx performs money transaction from an account to another
// using db transaction
func (store *SQLStore) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			Amount:        arg.Amount,
			FromAccountID: arg.FromAccountID,
			ToAccountID:   arg.ToAccountID,
		})
		if err != nil {
			return err
		}

		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount,
		})
		if err != nil {
			return err
		}

		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount,
		})
		if err != nil {
			return err
		}

		// Account balance implementation
		if arg.FromAccountID < arg.ToAccountID {
			result.FromAccount, result.ToAccount, err = store.updateBalance(ctx, arg.FromAccountID, arg.ToAccountID, arg.Amount)
		} else {
			result.FromAccount, result.ToAccount, err = store.updateBalance(ctx, arg.ToAccountID, arg.FromAccountID, -arg.Amount)
		}

		return err
	})

	return result, err
}

func (store *SQLStore) updateBalance(
	ctx context.Context,
	accountID1,
	accountID2,
	amount int64,
) (account1, account2 Account, err error) {
	account1, err = store.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     accountID1,
		Amount: -amount,
	})
	if err != nil {
		return
	}

	account2, err = store.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     accountID2,
		Amount: amount,
	})
	if err != nil {
		return
	}

	return
}
