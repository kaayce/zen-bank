package db

import (
	"context"
	"database/sql"
	"fmt"

	_ "go.uber.org/mock/mockgen/model"
)

type Store interface {
	Querier
	TransferTX(ctx context.Context, arg TransferTxParams) (TransferTxResult, error)
}

// Provides all functions to perform db queries and transactions
type SQLStore struct {
	*Queries
	db *sql.DB
}

// Creates new store
func NewStore(db *sql.DB) Store {
	return &SQLStore{
		db:      db,
		Queries: New(db),
	}
}

// executes a function within a db transaction
func (store *SQLStore) execTx(ctx context.Context, fn func(*Queries) error) error {
	// initiate tx
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	queries := New(tx)
	err = fn(queries)
	if err != nil {
		// if err, rollback tx
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx error %v, rb error %v", err, rbErr)
		}
	}
	// if no errors, commit tx
	return tx.Commit()
}

// contains input params of transfer tx
type TransferTxParams struct {
	FromAccountId int64 `json:"from_account_id"`
	ToAccountId   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

// contains rsult of the Transfer tx
type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

// performs a money transaction from one account to the other
// creates transfer record, add account entries and update acount balance within a single db tx
func (store *SQLStore) TransferTX(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountId,
			ToAccountID:   arg.ToAccountId,
			Amount:        arg.Amount,
		})
		if err != nil {
			return err
		}

		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountId,
			Amount:    -arg.Amount,
		})
		if err != nil {
			return err
		}

		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountId,
			Amount:    arg.Amount,
		})
		if err != nil {
			return err
		}

		// Prevent deadlock by ensuring a consistent order of operations based on account IDs
		if arg.FromAccountId < arg.ToAccountId {
			// update account's balance
			result.FromAccount, result.ToAccount, err = addMoney(ctx, q, addMoneyParams{
				accountID1: arg.FromAccountId,
				amount1:    -arg.Amount,
				accountID2: arg.ToAccountId,
				amount2:    arg.Amount,
			})
		} else {
			// swap flow
			result.ToAccount, result.FromAccount, err = addMoney(ctx, q, addMoneyParams{
				accountID1: arg.ToAccountId,
				amount1:    arg.Amount,
				accountID2: arg.FromAccountId,
				amount2:    -arg.Amount,
			})
		}
		return err
	})
	return result, err
}

// params for addMoney func
type addMoneyParams struct {
	accountID1 int64
	amount1    int64
	accountID2 int64
	amount2    int64
}

func addMoney(ctx context.Context, q *Queries, params addMoneyParams) (account1 Account, account2 Account, err error) {
	// account 1
	account1, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     params.accountID1,
		Amount: params.amount1,
	})
	if err != nil {
		return // named return
	}

	// account 2
	account2, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     params.accountID2,
		Amount: params.amount2,
	})
	return
}
