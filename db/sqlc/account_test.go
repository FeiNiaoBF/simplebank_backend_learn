package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/FeiNiaoBF/simplebank/util"
	"github.com/stretchr/testify/require"
)

func createRandomAccount(t *testing.T) Account {
	t.Helper()
	arg := CreateAccountParams{
		Owner:    util.RandomOwner(),
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}
	// Create a random account
	account, err := testQueries.CreateAccount(context.Background(), arg)
	// Check if there is an error
	require.NoError(t, err)
	require.NotEmpty(t, account)
	// Check if the account is the same as the argument
	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)
	// Check if the account has an ID and CreatedAt
	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)

	return account
}

func TestCreateAccount(t *testing.T) {
	createRandomAccount(t)
}

// TestGetAccount tests the GetAccount function
func TestGetAccount(t *testing.T) {
	wantAc := createRandomAccount(t)
	// Get the account
	gotAc, err := testQueries.GetAccount(context.Background(), wantAc.ID)
	// Check if there is an error
	require.NoError(t, err)
	require.NotEmpty(t, gotAc)
	// Check if the account is the same as the argument
	require.Equal(t, wantAc.ID, gotAc.ID)
	require.Equal(t, wantAc.Owner, gotAc.Owner)
	require.Equal(t, wantAc.Balance, gotAc.Balance)
	require.Equal(t, wantAc.Currency, gotAc.Currency)
	// Check if the account has an ID and CreatedAt
	require.WithinDuration(t, wantAc.CreatedAt, gotAc.CreatedAt, time.Second)
}

// TestUpdateAccount tests the UpdateAccount function
func TestUpdateAccount(t *testing.T) {
	wantAc := createRandomAccount(t)
	arg := UpdateAccountParams{
		ID:      wantAc.ID,
		Balance: util.RandomMoney(),
	}

	gotAc, err := testQueries.UpdateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, gotAc)
	// Check if the account is the same as the argument
	require.Equal(t, wantAc.ID, gotAc.ID)
	require.Equal(t, wantAc.Owner, gotAc.Owner)
	require.Equal(t, arg.Balance, gotAc.Balance)
	require.Equal(t, wantAc.Currency, gotAc.Currency)
	// Check if the account has an ID and CreatedAt
	require.WithinDuration(t, wantAc.CreatedAt, gotAc.CreatedAt, time.Second)
}

// TestDeleteAccount tests the DeleteAccount function
func TestDeleteAccount(t *testing.T) {
	wantAc := createRandomAccount(t)

	err := testQueries.DeleteAccount(context.Background(), wantAc.ID)
	require.NoError(t, err)

	gotAc, err := testQueries.GetAccount(context.Background(), wantAc.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, gotAc)
}

// TestListAccounts tests the ListAccounts function
func TestListAccounts(t *testing.T) {
	const (
		n      = 10
		limit  = 10
		offset = 5
	)
	for i := 0; i < n; i++ {
		createRandomAccount(t)
	}

	args := ListAccountsParams{
		Limit:  limit,
		Offset: offset,
	}

	accounts, err := testQueries.ListAccounts(context.Background(), args)
	require.NoError(t, err)
	require.Len(t, accounts, limit)

	for _, account := range accounts {
		require.NotEmpty(t, account)
	}
}
