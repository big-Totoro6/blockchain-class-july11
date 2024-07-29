package state

import "github.com/ardanlabs/blockchain/foundation/blockchain/database"

// UpsertWalletTransaction accepts a transaction from a wallet for inclusion.
func (s *State) UpsertWalletTransaction(signedTx database.SignedTx) error {

	// CORE NOTE: It's up to the wallet to make sure the account has a proper
	// balance and this transaction has a proper nonce. Fees will be taken if
	// this transaction is mined into a block it doesn't have enough money to
	// pay or the nonce isn't the next expected nonce for the account.
	//钱包必须确保账户有足够的余额来支付交易费用。
	//交易的 nonce 必须是账户的下一个预期 nonce。
	//如果交易在被挖掘到区块中时，账户余额不足以支付费用，或者 nonce 不正确，交易将会失败。

	// Check the signed transaction has a proper signature, the from matches the
	// signature, and the from and to fields are properly formatted.
	if err := signedTx.Validate(s.genesis.ChainID); err != nil {
		return err
	}

	const oneUnitOfGas = 1
	tx := database.NewBlockTx(signedTx, s.genesis.GasPrice, oneUnitOfGas)
	if err := s.mempool.Upsert(tx); err != nil {
		return err
	}

	return nil
}
