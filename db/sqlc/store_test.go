package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(testDB)

	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	// run n concurrent transfer transactions
	var (
		n       = 5
		amount  = int64(10)
		errs    = make(chan error)
		results = make(chan TransferTxResult)
	)

	for i := 0; i < n; i++ {
		go func() {
			result, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
				Amount:        amount,
			})

			errs <- err
			results <- result
		}()
	}

	for i := 0; i < n; i++ {
		// check errors
		err := <-errs
		require.NoError(t, err)
		// check results
		result := <-results
		require.NotEmpty(t, result)
		// check transfer
		assertTransfer(t, store, result, account1, account2, amount)
		// check entries
		assertEntries(t, store, result, account1, account2, amount)
		// TODO: check accounts' balance
	}
}

// assertTransfer checks the transfer result
func assertTransfer(t testing.TB, store *Store, result TransferTxResult, account1, account2 Account, amount int64) {
	t.Helper()
	transfer := result.Transfer
	require.Equal(t, account1.ID, transfer.FromAccountID)
	require.Equal(t, account2.ID, transfer.ToAccountID)
	require.Equal(t, amount, transfer.Amount)
	require.NotZero(t, transfer.ID)
	require.NotZero(t, transfer.CreatedAt)
	_, err := store.query.GetTransfer(context.Background(), transfer.ID)
	require.NoError(t, err)
}

// assertEntries checks the entries created by the transfer
func assertEntries(t testing.TB, store *Store, result TransferTxResult, account1, account2 Account, amount int64) {
	t.Helper()
	var err error
	// check the from entry
	fromEntry := result.FromEntry
	require.NotEmpty(t, fromEntry)
	require.Equal(t, account1.ID, fromEntry.AccountID)
	require.Equal(t, -amount, fromEntry.Amount)
	require.NotZero(t, fromEntry.ID)
	require.NotZero(t, fromEntry.CreatedAt)

	_, err = store.query.GetEntry(context.Background(), fromEntry.ID)
	require.NoError(t, err)

	// check the to entry
	toEntry := result.ToEntry
	require.NotEmpty(t, toEntry)
	require.Equal(t, account2.ID, toEntry.AccountID)
	require.Equal(t, amount, toEntry.Amount)
	require.NotZero(t, toEntry.ID)
	require.NotZero(t, toEntry.CreatedAt)

	_, err = store.query.GetEntry(context.Background(), toEntry.ID)
	require.NoError(t, err)
}
