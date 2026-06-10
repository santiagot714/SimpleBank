package api

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/santiagot714/SimpleBank/db/sqlc"
	"github.com/santiagot714/SimpleBank/token"
	"github.com/shopspring/decimal"
)

type transferRequest struct {
	OriginAccountID      int64   `json:"origin_account_id" binding:"required,min=1"`
	DestinationAccountID int64   `json:"destination_account_id" binding:"required,min=1"`
	Amount               float64 `json:"amount" binding:"required,gt=0"`
	Currency             string  `json:"currency" binding:"required,currency"`
}

// transfer handles the transfer of money between two accounts.
// TODO: Create tests for this function
func (server *Server) createTransfer(ctx *gin.Context) {
	var req transferRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	originAccount, valid := server.validAccount(ctx, req.OriginAccountID, req.Currency)
	if !valid {
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	if originAccount.Owner != authPayload.Username {
		err := errors.New("origin account does not belong to the authenticated user")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}
	_, valid = server.validAccount(ctx, req.DestinationAccountID, req.Currency)
	if !valid {
		return
	}

	arg := db.TransferTxParams{
		OriginAccountID:      req.OriginAccountID,
		DestinationAccountID: req.DestinationAccountID,
		Amount:               decimal.NewFromFloat(req.Amount),
	}
	result, err := server.store.TransferTx(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, result)
}

// validAccount checks if the account is valid and has the correct currency.
func (server *Server) validAccount(ctx *gin.Context, accountID int64, currency string) (db.Account, bool) {
	account, err := server.store.GetAccount(ctx, accountID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return db.Account{}, false
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return db.Account{}, false
	}
	if account.Currency != currency {
		err := fmt.Errorf("account [%d] currency mismatch: %s vs %s", accountID, account.Currency, currency)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return db.Account{}, false
	}
	return account, true
}
