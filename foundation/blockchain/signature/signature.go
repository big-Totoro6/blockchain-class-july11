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
