// Package public maintains the group of handlers for public access.
package public

import (
	"context"
	"github.com/ardanlabs/blockchain/foundation/blockchain/database"
	"github.com/ardanlabs/blockchain/foundation/blockchain/state"
	"net/http"

	"github.com/ardanlabs/blockchain/foundation/web"
	"go.uber.org/zap"
)

// Handlers manages the set of bar ledger endpoints.
type Handlers struct {
	Log   *zap.SugaredLogger
	State *state.State
}

// Sample just provides a starting point for the class.
func (h Handlers) Sample(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	resp := struct {
		Status string
	}{
		Status: "OK",
	}

	return web.Respond(ctx, w, resp, http.StatusOK)
}

// Genesis returns the genesis information.
func (h Handlers) Genesis(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	gen := h.State.Genesis()
	return web.Respond(ctx, w, gen, http.StatusOK)
}

// Accounts returns the current balances for all users.
func (h Handlers) Accounts(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	accountStr := web.Param(r, "account")

	var accounts map[database.AccountID]database.Account
	switch accountStr {
	case "":
		accounts = h.State.Accounts()

	default:
		accountID, err := database.ToAccountID(accountStr)
		if err != nil {
			return err
		}
		account, err := h.State.QueryAccount(accountID)
		if err != nil {
			return err
		}
		accounts = map[database.AccountID]database.Account{accountID: account}
	}

	return web.Respond(ctx, w, accounts, http.StatusOK)
}

// Mempool returns the set of uncommitted transactions.
func (h Handlers) Mempool(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	acct := web.Param(r, "account")

	mempool := h.State.Mempool()

	//保证即使内存池空的也会返回一个空的交易列表
	trans := []tx{}
	//遍历 mempool 中的所有交易，并筛选出与特定账户（acct）相关的交易。如果交易的发送方或接收方与指定的账户匹配，则继续处理该交易；否则，跳过该交易
	for _, tran := range mempool {
		//把与账户匹配的交易赛选
		if acct != "" && ((acct != string(tran.FromID)) && (acct != string(tran.ToID))) {
			continue
		}

		trans = append(trans, tx{
			FromAccount: tran.FromID,
			To:          tran.ToID,
			ChainID:     tran.ChainID,
			Nonce:       tran.Nonce,
			Value:       tran.Value,
			Tip:         tran.Tip,
			Data:        tran.Data,
			TimeStamp:   tran.TimeStamp,
			GasPrice:    tran.GasPrice,
			GasUnits:    tran.GasUnits,
			Sig:         tran.SignatureString(),
		})
	}

	return web.Respond(ctx, w, trans, http.StatusOK)
}
