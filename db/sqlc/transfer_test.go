package db

import (
	"context"
	"testing"
	"time"

	"github.com/FeiNiaoBF/simplebank/util"
	"github.com/stretchr/testify/require"
)

func createRandomTransfer(t *testing.T, fromAccount, ToAccount Account) Transfer {
	t.Helper()
	arg := CreateTransferParams{
		FromAccountID: fromAccount.ID,
		ToAccountID:   ToAccount.ID,
		Amount:        util.RandomMoney(),
	}
	// Create a random account
	transfer, err := testQueries.CreateTransfer(context.Background(), arg)
	// Check if there is an error
	require.NoError(t, err)
	require.NotEmpty(t, transfer)

	require.Equal(t, arg.FromAccountID, transfer.FromAccountID)
	require.Equal(t, arg.ToAccountID, transfer.ToAccountID)
	require.Equal(t, arg.Amount, transfer.Amount)

	require.NotZero(t, transfer.ID)
	require.NotZero(t, transfer.CreatedAt)

	return transfer
}

// TestCreateTransfer tests the CreateTransfer function
func TestCreateTransfer(t *testing.T) {
	fromAccount := createRandomAccount(t)
	toAccount := createRandomAccount(t)
	createRandomTransfer(t, fromAccount, toAccount)
}

// TestGetTransfer tests the GetTransfer function
func TestGetTransfer(t *testing.T) {
	fromAccount := createRandomAccount(t)
	toAccount := createRandomAccount(t)
	gotT := createRandomTransfer(t, fromAccount, toAccount)

	wantT, err := testQueries.GetTransfer(context.Background(), gotT.ID)
	require.NoError(t, err)
	require.NotEmpty(t, wantT)

	require.Equal(t, gotT.ID, wantT.ID)
	require.Equal(t, gotT.FromAccountID, wantT.FromAccountID)
	require.Equal(t, gotT.ToAccountID, wantT.ToAccountID)
	require.Equal(t, gotT.Amount, wantT.Amount)
	require.WithinDuration(t, gotT.CreatedAt, wantT.CreatedAt, time.Second)
}

// TestListTransfers tests the ListTransfers function
func TestListTransfers(t *testing.T) {
	const (
		n      = 10
		limit  = 5
		offset = 5
	)

	fromAccount := createRandomAccount(t)
	toAccount := createRandomAccount(t)
	for i := 0; i < n; i++ {
		createRandomTransfer(t, fromAccount, toAccount)
		createRandomTransfer(t, toAccount, fromAccount)
	}

	args := ListTransfersParams{
		FromAccountID: fromAccount.ID,
		ToAccountID:   toAccount.ID,
		Limit:         limit,
		Offset:        offset,
	}

	transfers, err := testQueries.ListTransfers(context.Background(), args)
	require.NoError(t, err)
	require.Len(t, transfers, limit)

	for _, transfer := range transfers {
		require.NotEmpty(t, transfer)
		require.True(t, transfer.FromAccountID == fromAccount.ID || transfer.ToAccountID == toAccount.ID)
	}
}
