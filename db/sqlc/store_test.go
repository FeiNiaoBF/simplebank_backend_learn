package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(testDB)

	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	// fmt.Println(">> before:", account1.Balance, account2.Balance)
	// run n concurrent transfer transactions
	var (
		n       = 5                           // number of concurrent transactions
		amount  = int64(10)                   // amount to money to transfer
		errs    = make(chan error)            // channel to errors
		results = make(chan TransferTxResult) // channel to results
	)

	for i := 0; i < n; i++ {
		txName := fmt.Sprintf("tx %d", i+1)

		go func() {
			ctx := context.WithValue(context.Background(), txKey, txName)
			result, err := store.TransferTx(ctx, TransferTxParams{
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
		assertAccount(t, result, account1, account2, amount, n)
	}
	// check account after thr update
	assertUpdate(t, account1, account2, amount, n)
}

func TestTransferTxDeadlock(t *testing.T) {
	store := NewStore(testDB)

	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	fmt.Println(">> before:", account1.Balance, account2.Balance)
	// run n concurrent transfer transactions
	var (
		n      = 10               // number of concurrent transactions
		amount = int64(10)        // amount to money to transfer
		errs   = make(chan error) // channel to errors
	)

	for i := 0; i < n; i++ {
		txName := fmt.Sprintf("tx %d", i+1)
		fromAccountId := account1.ID
		toAccountId := account2.ID
		if i%2 == 1 {
			fromAccountId = account2.ID
			toAccountId = account1.ID
		}

		go func() {
			ctx := context.WithValue(context.Background(), txKey, txName)
			_, err := store.TransferTx(ctx, TransferTxParams{
				FromAccountID: fromAccountId,
				ToAccountID:   toAccountId,
				Amount:        amount,
			})

			errs <- err
		}()
	}

	for i := 0; i < n; i++ {
		// check errors
		err := <-errs
		require.NoError(t, err)
	}
	// check account after thr update
	assertUpdateDeadlock(t, account1, account2)
}

// assertAccount checks the account result
func assertAccount(t testing.TB, result TransferTxResult, account1, account2 Account, amount int64, times int) {
	t.Helper()
	existed := make(map[int]bool)
	// check from account
	fromAccount := result.FromAccount
	require.NotEmpty(t, fromAccount)
	require.Equal(t, account1.ID, fromAccount.ID)
	// check to account
	toAccount := result.ToAccount
	require.NotEmpty(t, toAccount)
	require.Equal(t, account2.ID, toAccount.ID)

	// check the difference
	fmt.Println(">> tx:", fromAccount.Balance, toAccount.Balance)
	diff1 := account1.Balance - fromAccount.Balance
	diff2 := toAccount.Balance - account2.Balance
	require.Equal(t, diff1, diff2)
	require.True(t, diff1 > 0)
	require.True(t, diff1%amount == 0)
	k := int(diff1 / amount)
	require.True(t, k >= 1 && k <= times)
	require.NotContains(t, existed, k)
	existed[k] = true
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

// assertUpdate checks the updated account's balance
func assertUpdate(t testing.TB, account1, account2 Account, amount int64, times int) {
	t.Helper()
	updateAccount1, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)
	updateAccount2, err := testQueries.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)
	// fmt.Println("<< after:", updateAccount1.Balance, updateAccount2.Balance)
	require.Equal(t, account1.Balance-int64(times)*amount, updateAccount1.Balance)
	require.Equal(t, account2.Balance+int64(times)*amount, updateAccount2.Balance)
}

func assertUpdateDeadlock(t testing.TB, account1, account2 Account) {
	t.Helper()
	updateAccount1, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)
	updateAccount2, err := testQueries.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)
	fmt.Println("<< after:", updateAccount1.Balance, updateAccount2.Balance)
	require.Equal(t, account1.Balance, updateAccount1.Balance)
	require.Equal(t, account2.Balance, updateAccount2.Balance)
}
