package signature

import (
	"crypto/ecdsa"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
	"math/big"
)

const ardanID = 29

func Sign(value any, privateKey *ecdsa.PrivateKey) (v, r, s *big.Int, err error) {
	//需要数据 被盖章的数据 供签名使用
	data, err := stamp(value)
	if err != nil {
		return nil, nil, nil, err
	}
	//然后通过私钥与数据 得到签名
	sig, err := crypto.Sign(data, privateKey)
	if err != nil {
		return nil, nil, nil, err
	}

	// Extract the bytes for the original public key.
	publicKeyOrg := privateKey.Public()
	publicKeyECDSA, ok := publicKeyOrg.(*ecdsa.PublicKey)
	if !ok {
		return nil, nil, nil, errors.New("error casting public key to ECDSA")
	}
	publicKeyBytes := crypto.FromECDSAPub(publicKeyECDSA)

	// Check the public key validates the data and signature.
	rs := sig[:crypto.RecoveryIDOffset]                    //从签名 sig 中提取 R 和 S 的组合
	if !crypto.VerifySignature(publicKeyBytes, data, rs) { //验证公钥 publicKeyBytes 是否能验证签名 sig
		return nil, nil, nil, errors.New("invalid signature produced")
	}

	//把前面变成RSV格式返回
	v, r, s = toSignatureValues(sig)
	return v, r, s, err
}

// =============================================================================

// stamp returns a hash of 32 bytes that represents this data with
// the Ardan stamp embedded into the final hash.
func stamp(value any) ([]byte, error) {

	// Marshal the data.
	v, err := json.Marshal(value)
	if err != nil {
		return nil, err
	}

	// This stamp is used so signatures we produce when signing data
	// are always unique to the Ardan blockchain.
	stamp := []byte(fmt.Sprintf("\x19Ardan Signed Message:\n%d", len(v)))

	// Hash the stamp and txHash together in a final 32 byte array
	// that represents the data.
	data := crypto.Keccak256(stamp, v)

	return data, nil
}

// toSignatureValues converts the signature into the r, s, v values.
func toSignatureValues(sig []byte) (v, r, s *big.Int) {
	r = big.NewInt(0).SetBytes(sig[:32])
	s = big.NewInt(0).SetBytes(sig[32:64])
	v = big.NewInt(0).SetBytes([]byte{sig[64] + ardanID})

	return v, r, s
}

// VerifySignature verifies the signature conforms to our standards.
func VerifySignature(v, r, s *big.Int) error {

	// Check the recovery id is either 0 or 1.
	uintV := v.Uint64() - ardanID
	if uintV != 0 && uintV != 1 {
		return errors.New("invalid recovery id")
	}

	// Check the signature values are valid.
	if !crypto.ValidateSignatureValues(byte(uintV), r, s, false) {
		return errors.New("invalid signature values")
	}

	return nil
}

// ToSignatureBytes converts the r, s, v values into a slice of bytes
// with the removal of the ardanID.
func ToSignatureBytes(v, r, s *big.Int) []byte {
	sig := make([]byte, crypto.SignatureLength)

	rBytes := make([]byte, 32)
	r.FillBytes(rBytes)
	copy(sig, rBytes)

	sBytes := make([]byte, 32)
	s.FillBytes(sBytes)
	copy(sig[32:], sBytes)

	sig[64] = byte(v.Uint64() - ardanID)

	return sig
}

// FromAddress extracts the address for the account that signed the data.
func FromAddress(value any, v, r, s *big.Int) (string, error) {

	// Prepare the data for public key extraction.
	data, err := stamp(value)
	if err != nil {
		return "", err
	}

	// Convert the [R|S|V] format into the original 65 bytes.
	sig := ToSignatureBytes(v, r, s)

	// Capture the public key associated with this data and signature.
	publicKey, err := crypto.SigToPub(data, sig)
	if err != nil {
		return "", err
	}

	// Extract the account address from the public key.
	return crypto.PubkeyToAddress(*publicKey).String(), nil
}
