package db

import (
	"context"
	"testing"

	"github.com/santiagot714/SimpleBank/util"
	decimal "github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(testDB)
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	amount := util.RandomMoney()
	n := 5
	errs := make(chan error)
	results := make(chan TransferTxResult)

	for i := 0; i < n; i++ {
		go func() {
			result, err := store.TransferTx(context.Background(), TransferTxParams{
				OriginAccountID:      account1.ID,
				DestinationAccountID: account2.ID,
				Amount:               amount,
			})

			errs <- err
			results <- result
		}()
	}
	registeredTransactions := make(map[int]bool)
	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)
		result := <-results
		require.NotEmpty(t, result)

		// check transfer
		transfer := result.Transfer
		require.NotEmpty(t, transfer)
		require.Equal(t, account1.ID, transfer.OriginAccountID)
		require.Equal(t, account2.ID, transfer.DestinationAccountID)
		require.True(t, amount.Equal(transfer.Amount))
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		_, err = store.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)

		// check entries
		originEntry := result.OriginEntry
		require.NotEmpty(t, originEntry)
		require.Equal(t, account1.ID, originEntry.AccountID)
		require.True(t, amount.Neg().Equal(originEntry.Amount))
		require.NotZero(t, originEntry.ID)
		require.NotZero(t, originEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), originEntry.ID)
		require.NoError(t, err)

		destinationEntry := result.DestinationEntry
		require.NotEmpty(t, destinationEntry)
		require.Equal(t, account2.ID, destinationEntry.AccountID)
		require.True(t, amount.Equal(destinationEntry.Amount))
		require.NotZero(t, destinationEntry.ID)
		require.NotZero(t, destinationEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), destinationEntry.ID)
		require.NoError(t, err)

		// check accounts
		originAccount := result.OriginAccount
		require.NotEmpty(t, originAccount)
		require.Equal(t, account1.ID, originAccount.ID)

		destinationAccount := result.DestinationAccount
		require.NotEmpty(t, destinationAccount)
		require.Equal(t, account2.ID, destinationAccount.ID)

		// check balances
		diff1 := account1.Balance.Sub(originAccount.Balance)
		diff2 := destinationAccount.Balance.Sub(account2.Balance)

		require.True(t, diff1.Equal(diff2))
		require.Equal(t, 1, diff1.Sign())
		require.True(t, diff1.Mod(amount).Equal(decimal.Zero))

		k := int(diff1.Div(amount).IntPart())
		require.True(t, k >= 1 && k <= n)
		require.NotContains(t, registeredTransactions, k)
		registeredTransactions[k] = true
	}
	// check the final updated balance
	updatedAccount1, err := store.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	updatedAccount2, err := store.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)

	require.True(t, updatedAccount1.Balance.Equal(account1.Balance.Sub(amount.Mul(decimal.NewFromInt(int64(n))))))
	require.True(t, updatedAccount2.Balance.Equal(account2.Balance.Add(amount.Mul(decimal.NewFromInt(int64(n))))))
}

func TestTransferTxDeadlock(t *testing.T) {
	store := NewStore(testDB)
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	amount := decimal.NewFromInt(10)
	n := 10
	errs := make(chan error)

	for i := 0; i < n; i++ {
		fromAccountID := account1.ID
		toAccountID := account2.ID
		if i%2 == 1 {
			fromAccountID = account2.ID
			toAccountID = account1.ID
		}
		go func() {
			_, err := store.TransferTx(context.Background(), TransferTxParams{
				OriginAccountID:      fromAccountID,
				DestinationAccountID: toAccountID,
				Amount:               amount,
			})

			errs <- err
		}()
	}
	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)
	}
	// check the final updated balance
	updatedAccount1, err := store.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	updatedAccount2, err := store.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)

	require.True(t, updatedAccount1.Balance.Equal(account1.Balance))
	require.True(t, updatedAccount2.Balance.Equal(account2.Balance))
}
