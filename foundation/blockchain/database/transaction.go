package database

// Tx is the transactional information between two parties.
//type Tx struct {
//	ChainID uint16    `json:"chain_id"` // Ethereum: The chain id that is listed in the genesis file.
//	Nonce   uint64    `json:"nonce"`    // Ethereum: Unique id for the transaction supplied by the user.
//	FromID  AccountID `json:"from"`     // Ethereum: Account sending the transaction. Will be checked against signature.
//	ToID    AccountID `json:"to"`       // Ethereum: Account receiving the benefit of the transaction.
//	Value   uint64    `json:"value"`    // Ethereum: Monetary value received from this transaction.
//	Tip     uint64    `json:"tip"`      // Ethereum: Tip offered by the sender as an incentive to mine this transaction.
//	Data    []byte    `json:"data"`     // Ethereum: Extra data related to the transaction.
//}
//
//// SignedTx 这是个加过签名的交易类型 这个签名符合以太坊的规范==============================================================================================
//// SignedTx is a signed version of the transaction. This is how clients like
//// a wallet provide transactions for inclusion into the blockchain.
//type SignedTx struct {
//	Tx
//	V *big.Int `json:"v"` // Ethereum: Recovery identifier, either 29 or 30 with ardanID.
//	R *big.Int `json:"r"` // Ethereum: First coordinate of the ECDSA signature.
//	S *big.Int `json:"s"` // Ethereum: Second coordinate of the ECDSA signature.
//}
