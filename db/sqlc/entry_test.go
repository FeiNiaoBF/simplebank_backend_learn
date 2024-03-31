package db

import (
	"context"
	"testing"
	"time"

	"github.com/FeiNiaoBF/simplebank/util"
	"github.com/stretchr/testify/require"
)

func createRandomEntry(t *testing.T, account Account) Entry {
	t.Helper()
	arg := CreateEntryParams{
		AccountID: account.ID,
		Amount:    util.RandomAmount(),
	}
	// Create a random account
	entry, err := testQueries.CreateEntry(context.Background(), arg)
	// Check if there is an error
	require.NoError(t, err)
	require.NotEmpty(t, entry)

	require.Equal(t, arg.AccountID, entry.AccountID)
	require.Equal(t, arg.Amount, entry.Amount)

	require.NotZero(t, entry.ID)
	require.NotZero(t, entry.CreatedAt)

	return entry
}

// TestCreateEntry tests the CreateEntry function
func TestCreateEntry(t *testing.T) {
	account := createRandomAccount(t)
	createRandomEntry(t, account)
}

// TestGetEntry tests the GetEntry function
func TestGetEntry(t *testing.T) {
	account := createRandomAccount(t)
	gotE := createRandomEntry(t, account)

	wantE, err := testQueries.GetEntry(context.Background(), gotE.ID)

	require.NoError(t, err)
	require.NotEmpty(t, wantE)

	require.Equal(t, gotE.ID, wantE.ID)
	require.Equal(t, gotE.AccountID, wantE.AccountID)
	require.Equal(t, gotE.Amount, wantE.Amount)
	require.WithinDuration(t, gotE.CreatedAt, wantE.CreatedAt, time.Second)
}

// TestListEntries tests the ListEntries function
func TestListEntries(t *testing.T) {
	const (
		n      = 10
		limit  = 5
		offset = 5
	)

	account := createRandomAccount(t)
	for i := 0; i < n; i++ {
		createRandomEntry(t, account)
	}

	args := ListEntriesParams{
		AccountID: account.ID,
		Limit:     limit,
		Offset:    offset,
	}

	entries, err := testQueries.ListEntries(context.Background(), args)
	require.NoError(t, err)
	require.Len(t, entries, limit)

	for _, entry := range entries {
		require.NotEmpty(t, entry)
		require.Equal(t, args.AccountID, entry.AccountID)
	}
}
