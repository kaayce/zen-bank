package database

import (
	"context"
	"fmt"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(testDB)

	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	fmt.Println(">> before:", account1.Balance, account2.Balance)

	// to avoid db errors, we use seperate routines to run transactions and communicate result with main test routine
	// run n concurrent transfer transactions
	n := 5
	amount := int64(10)
	txResults := make(chan TransferTxResult, n)
	txErrs := make(chan error, n)
	var wg sync.WaitGroup
	wg.Add(n)

	for i := 0; i < n; i++ {
		go func() {
			defer wg.Done()
			result, err := store.TransferTX(context.Background(), TransferTxParams{
				FromAccountId: account1.ID,
				ToAccountId:   account2.ID,
				Amount:        amount,
			})
			txResults <- result
			txErrs <- err
		}()
	}

	go func() {
		wg.Wait()
		close(txResults)
		close(txErrs)
	}()

	// check errors & results
	existed := make(map[int]bool)
	for i := 0; i < n; i++ {
		err := <-txErrs
		require.NoError(t, err)

		result := <-txResults
		require.NotEmpty(t, result)

		// check transfer
		transfer := result.Transfer
		require.NotEmpty(t, transfer)
		require.Equal(t, account1.ID, transfer.FromAccountID)
		require.Equal(t, account2.ID, transfer.ToAccountID)
		require.Equal(t, amount, transfer.Amount)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)
		require.NotZero(t, transfer.UpdatedAt)

		// make sure tranfer is stored in DB
		_, err = store.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)

		// check from entries
		fromEntry := result.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, account1.ID, fromEntry.AccountID)
		require.Equal(t, -amount, fromEntry.Amount)
		require.NotZero(t, fromEntry.AccountID)
		require.NotZero(t, fromEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), fromEntry.ID)
		require.NoError(t, err)

		// check to entries
		toEntry := result.ToEntry
		require.NotEmpty(t, toEntry)
		require.Equal(t, account2.ID, toEntry.AccountID)
		require.Equal(t, amount, toEntry.Amount)
		require.NotZero(t, toEntry.AccountID)
		require.NotZero(t, toEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), toEntry.ID)
		require.NoError(t, err)

		// check accounts
		// from account
		fromAccount := result.FromAccount
		require.NotEmpty(t, fromAccount)
		require.Equal(t, account1.ID, fromAccount.ID)

		// to account
		toAccount := result.ToAccount
		require.NotEmpty(t, toAccount)
		require.Equal(t, account2.ID, toAccount.ID)

		// check account balance
		fmt.Println(">> tx:", fromAccount.Balance, toAccount.Balance)
		diffFromAcc := account1.Balance - fromAccount.Balance
		diffToAcc := toAccount.Balance - account2.Balance
		require.Equal(t, diffFromAcc, diffToAcc)
		require.True(t, diffFromAcc > 0)
		require.True(t, diffFromAcc%amount == 0) // amount. 2 * amount, 3 * amount...

		k := int(diffFromAcc / amount)
		require.True(t, k >= 1 && k <= n)
		require.NotContains(t, existed, k)
		existed[k] = true
	}
	//  check final updated balance
	updatedAccountFrom, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, updatedAccountFrom)

	updatedAccountTo, err := testQueries.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)
	require.NotEmpty(t, updatedAccountTo)

	fmt.Println(">> after:", updatedAccountFrom.Balance, updatedAccountTo.Balance)
	require.Equal(t, account1.Balance-int64(n)*amount, updatedAccountFrom.Balance)
	require.Equal(t, account2.Balance+int64(n)*amount, updatedAccountTo.Balance)
}

func TestTransferTxDeadlock(t *testing.T) {
	store := NewStore(testDB)

	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	fmt.Println(">> before:", account1.Balance, account2.Balance)

	// to avoid db errors, we use seperate routines to run transactions and communicate result with main test routine
	// run n concurrent transfer transactions
	n := 10
	amount := int64(10)
	txErrs := make(chan error, n)
	var wg sync.WaitGroup
	wg.Add(n)

	for i := 0; i < n; i++ {
		fromAccountID := account1.ID
		toAccountID := account2.ID

		// rule to prevent deadlocks
		if i%2 == 1 {
			fromAccountID = account2.ID
			toAccountID = account1.ID
		}

		go func() {
			defer wg.Done()
			_, err := store.TransferTX(context.Background(), TransferTxParams{
				FromAccountId: fromAccountID,
				ToAccountId:   toAccountID,
				Amount:        amount,
			})
			txErrs <- err
		}()
	}

	go func() {
		wg.Wait()
		close(txErrs)
	}()

	// check errors & results
	for i := 0; i < n; i++ {
		err := <-txErrs
		require.NoError(t, err)
	}
	//  check final updated balance
	updatedAccountFrom, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, updatedAccountFrom)

	updatedAccountTo, err := testQueries.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)
	require.NotEmpty(t, updatedAccountTo)

	fmt.Println(">> after:", updatedAccountFrom.Balance, updatedAccountTo.Balance)
	require.Equal(t, account1.Balance, updatedAccountFrom.Balance)
	require.Equal(t, account2.Balance, updatedAccountTo.Balance)
}
