package db

import (
	"context"
	"database/sql"
	"fmt"

	decimal "github.com/shopspring/decimal"
)

// Store provides all functions to execute db queries and transactions
type Store struct {
	db *sql.DB
	*Queries
}

// NewStore creates a new store
func NewStore(db *sql.DB) *Store {
	return &Store{db: db, Queries: New(db)}
}

// execTx executes a function within a database transaction
func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %w, rb err: %w", err, rbErr)
		}
		return err
	}
	return tx.Commit()
}

// TransferTxParams contains the input parameters of the transfer transaction
type TransferTxParams struct {
	OriginAccountID      int64           `json:"origin_account_id"`
	DestinationAccountID int64           `json:"destination_account_id"`
	Amount               decimal.Decimal `json:"amount"`
}

// TransferTxResult is the result of the transfer transaction
type TransferTxResult struct {
	Transfer           Transfer `json:"transfer"`
	OriginAccount      Account  `json:"origin_account"`
	DestinationAccount Account  `json:"destination_account"`
	OriginEntry        Entry    `json:"origin_entry"`
	DestinationEntry   Entry    `json:"destination_entry"`
}

// TransferTx performs a money transfer from one account to the other
// It creates a transfer record, add account entries, and update accounts balance within a single database transaction
func (store *Store) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult
	err := store.execTx(ctx, func(q *Queries) error {
		var err error
		// Create transfer
		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams(arg))
		if err != nil {
			return err
		}
		// Create entries
		result.OriginEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.OriginAccountID,
			Amount:    arg.Amount.Neg(),
		})
		if err != nil {
			return err
		}
		result.DestinationEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.DestinationAccountID,
			Amount:    arg.Amount,
		})
		if err != nil {
			return err
		}

		// Update accounts balance
		if arg.OriginAccountID < arg.DestinationAccountID {
			result.OriginAccount, result.DestinationAccount, err = addMoney(ctx, q, arg.OriginAccountID, arg.Amount.Neg(), arg.DestinationAccountID, arg.Amount)
		} else {
			result.DestinationAccount, result.OriginAccount, err = addMoney(ctx, q, arg.DestinationAccountID, arg.Amount, arg.OriginAccountID, arg.Amount.Neg())
		}
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return TransferTxResult{}, err
	}
	return result, nil
}

func addMoney(
	ctx context.Context,
	q *Queries,
	accountID1 int64,
	amount1 decimal.Decimal,
	accountID2 int64,
	amount2 decimal.Decimal,
) (account1 Account, account2 Account, err error) {
	account1, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     accountID1,
		Amount: amount1,
	})
	if err != nil {
		return
	}
	account2, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     accountID2,
		Amount: amount2,
	})
	if err != nil {
		return
	}
	return account1, account2, nil
}
