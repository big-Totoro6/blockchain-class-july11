package database

import (
	"crypto/ecdsa"
	"errors"
	"fmt"
	"github.com/ardanlabs/blockchain/foundation/blockchain/signature"
	"math/big"
)

// ardanID is an arbitrary number for signing messages. This will make it
// clear that the signature comes from the Ardan blockchain.
// Ethereum and Bitcoin do this as well, but they use the value of 27.
const ardanID = 29

// Tx is the transactional information between two parties.
type Tx struct {
	ChainID uint16    `json:"chain_id"` // Ethereum: The chain id that is listed in the genesis file.
	Nonce   uint64    `json:"nonce"`    // Ethereum: Unique id for the transaction supplied by the user.
	FromID  AccountID `json:"from"`     // Ethereum: Account sending the transaction. Will be checked against signature.
	ToID    AccountID `json:"to"`       // Ethereum: Account receiving the benefit of the transaction.
	Value   uint64    `json:"value"`    // Ethereum: Monetary value received from this transaction.
	Tip     uint64    `json:"tip"`      // Ethereum: Tip offered by the sender as an incentive to mine this transaction.
	Data    []byte    `json:"data"`     // Ethereum: Extra data related to the transaction.
}

// SignedTx 这是个加过签名的交易类型 这个签名符合以太坊的规范==============================================================================================
// SignedTx is a signed version of the transaction. This is how clients like
// a wallet provide transactions for inclusion into the blockchain.
type SignedTx struct {
	Tx
	V *big.Int `json:"v"` // Ethereum: Recovery identifier, either 29 or 30 with ardanID.
	R *big.Int `json:"r"` // Ethereum: First coordinate of the ECDSA signature.
	S *big.Int `json:"s"` // Ethereum: Second coordinate of the ECDSA signature.
}

// NewTx constructs a new transaction.
func NewTx(chainID uint16, nonce uint64, fromID AccountID, toID AccountID, value uint64, tip uint64, data []byte) (Tx, error) {
	if !fromID.IsAccountID() {
		return Tx{}, errors.New("from account is not properly formatted")
	}
	if !toID.IsAccountID() {
		return Tx{}, errors.New("to account is not properly formatted")
	}

	tx := Tx{
		ChainID: chainID,
		Nonce:   nonce,
		FromID:  fromID,
		ToID:    toID,
		Value:   value,
		Tip:     tip,
		Data:    data,
	}

	return tx, nil
}

// Sign 现在我们需要一个方法 去把交易类型转化为 加签过的交易类型 （加签要给数据加入自己的stamp并且签名格式与RSV一致）
// Sign uses the specified private key to sign the transaction.
// 它接受一个私钥，使用 ECDSA 算法对交易进行签名，并将签名结果（v, r, s）与原始交易一起封装在 SignedTx 结构体中返回。
func (tx Tx) Sign(privateKey *ecdsa.PrivateKey) (SignedTx, error) {

	// Sign the transaction with the private key to produce a signature.
	v, r, s, err := signature.Sign(tx, privateKey)
	if err != nil {
		return SignedTx{}, err
	}

	// Construct the signed transaction by adding the signature
	// in the [R|S|V] format.
	signedTx := SignedTx{
		Tx: tx,
		V:  v,
		R:  r,
		S:  s,
	}

	return signedTx, nil
}

// 我需要一个验证的方法 很重要
func (tx SignedTx) Validate(chainID uint16) error {
	//首先就是判断chainid 这个是写在配置文件里面的 genesis里面的
	if tx.ChainID != chainID {
		return fmt.Errorf("invalid chain id, got[%d] exp[%d]", tx.ChainID, chainID)
	}

	//然后校验发送人 接收者
	if !tx.FromID.IsAccountID() {
		return errors.New("from account is not properly formatted")
	}

	if !tx.ToID.IsAccountID() {
		return errors.New("to account is not properly formatted")
	}

	//不能自己发给自己
	if tx.FromID == tx.ToID {
		return fmt.Errorf("transaction invalid, sending money to yourself, from %s, to %s", tx.FromID, tx.ToID)
	}
	//校验签名
	if err := signature.VerifySignature(tx.V, tx.R, tx.S); err != nil {
		return err
	}

	// 6. Extract address from signature and verify it matches fromID
	address, err := signature.FromAddress(tx.Tx, tx.V, tx.R, tx.S)
	if err != nil {
		return err
	}
	//校验地址 看你传过来的以太坊地址与我算出来的是否一致
	if address != string(tx.FromID) {
		return errors.New("signature address doesn't match from address")
	}
	return nil
}

// SignatureString returns the signature as a string.
func (tx SignedTx) SignatureString() string {
	return signature.SignatureString(tx.V, tx.R, tx.S)
}

// String implements the Stringer interface for logging.
func (tx SignedTx) String() string {
	return fmt.Sprintf("%s:%d", tx.FromID, tx.Nonce)
}
