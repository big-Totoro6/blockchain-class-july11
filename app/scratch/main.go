package main

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"log"
	"math/big"
	"strings"
)

// ZeroHash represents a hash code of zeros.
const ZeroHash string = "0x0000000000000000000000000000000000000000000000000000000000000000"

// ardanID is an arbitrary number for signing messages. This will make it
// clear that the signature comes from the Ardan blockchain.
// Ethereum and Bitcoin do this as well, but they use the value of 27.
const ardanID = 29

// =============================================================================

// Tx is the transactional information between two parties. 我们想查看交易时候发生的，我们给出的签名是啥样，计算出的公钥长啥样
type Tx struct {
	FromID string `json:"from"`  // Ethereum: Account sending the transaction. Will be checked against signature.
	ToID   string `json:"to"`    // Ethereum: Account receiving the benefit of the transaction.
	Value  uint64 `json:"value"` // Ethereum: Monetary value received from this transaction.
}

// 让我们展开一下，假设 run 函数可能包含了一系列的操作或者逻辑，如果其中的任何一步出现问题，就会返回一个非 nil 的错误。在这种情况下，main 函数负责捕获这个错误并输出到日志，然后退出程序，避免继续执行可能已经无法继续的任务。
func main() {
	if err := run(); err != nil {
		//日志输出: 使用 log.Fatalln(err) 来打印错误信息，它会将错误信息输出到标准错误并调用 os.Exit(1) 终止程序。这种方法适合于简单的命令行应用或者服务启动脚本。
		log.Fatalln(err)
	}

	//if err := publicKeyJudge(); err != nil {
	//	//日志输出: 使用 log.Fatalln(err) 来打印错误信息，它会将错误信息输出到标准错误并调用 os.Exit(1) 终止程序。这种方法适合于简单的命令行应用或者服务启动脚本。
	//	log.Fatalln(err)
	//}
}

func run() error {
	tx := Tx{
		FromID: "0xF01813E4B85e178A83e29B8E7bF26BD830a25f32", //交易发起人
		ToID:   "Aaron",                                      //交易的对方
		Value:  1000,                                         //交易的金额
	}
	//我们需要一个私钥 路径是我们的账户的路径
	path := "zblock/accounts/kennedy.ecdsa"
	privateKey, err := crypto.LoadECDSA(path)

	if err != nil {
		return fmt.Errorf("unable to load private key for node: %w", err)
	}

	//then 我们需要对这个数据tx进行一个编码 总不能原原本本发过去
	//JSON 编码: 使用 json.Marshal(tx) 将交易对象 tx 编码为 JSON 格式的字节流，这是为了确保交易数据能够被签名，而不是直接发送原始结构体。
	data, err := json.Marshal(tx)
	if err != nil {
		return fmt.Errorf("unable to Marshal: %w", err)
	}

	// 对数据进行处理=============================================================================

	//way1 对数据进行SHA-256哈希
	//hash := sha256.Sum256(data)
	////然后我们需要签名 用私钥去加密数据 得到签名
	////hash[:] 是一个 32 字节的哈希值的切片，确保它是正确计算和格式化的 SHA-256 哈希值。
	//sig, err := crypto.Sign(hash[:], privateKey)
	//if err != nil {
	//	return fmt.Errorf("unable to Sign: %w", err)
	//}

	// =============================================================================
	//way 2 以太坊自带的哈希函数
	//给数据盖章 加自己的身份标识
	// This stamp is used so signatures we produce when signing data
	// are always unique to the Ardan blockchain.
	stamp := []byte(fmt.Sprintf("\x19Ardan Signed Message:\n%d", len(data)))
	v := crypto.Keccak256(stamp, data)
	sig, err := crypto.Sign(v, privateKey)
	if err != nil {
		return fmt.Errorf("unable to Sign: %w", err)
	}

	// =============================================================================
	//需要变成16进制的
	fmt.Println(string(hexutil.Encode(sig)))

	//解析公钥
	publicKey, err := crypto.SigToPub(v, sig)
	if err != nil {
		return fmt.Errorf("unable to pub: %w", err)
	}
	//使用crypto.PubkeyToAddress将推导出的公钥转换为以太坊地址，并打印出来。
	fmt.Println(crypto.PubkeyToAddress(*publicKey).String())

	//打印一下被转换成VRS的签名
	//vv, r, s, err := ToVRSFromHexSig(string(hexutil.Encode(sig))) //传入的值0xd7f668a655b0d0e53f39ba1774fa748ea8458c544ee9df685f557f2d03a19dcb2d57caf751363dc5ab7ad37ec4b2645330db601aade65aec6b0787296362fcbc00
	vv, r, s, err := ToVRSFromHexSig(hex.EncodeToString(sig)) //      传入的值  d7f668a655b0d0e53f39ba1774fa748ea8458c544ee9df685f557f2d03a19dcb2d57caf751363dc5ab7ad37ec4b2645330db601aade65aec6b0787296362fcbc00

	if err != nil {
		fmt.Println("Error:", err)
		return nil
	}

	// Print v, r, s (just for demonstration)
	fmt.Println("v:", vv.String())
	fmt.Println("r:", r.String())
	fmt.Println("s:", s.String())

	return nil

}

// 验证一下 被同一个私钥加密的不同数据签名 解析出来是否是同一个公钥
func publicKeyJudge() error {
	tx1 := Tx{
		FromID: "Bill",  //交易发起人
		ToID:   "Aaron", //交易的对方
		Value:  1000,    //交易的金额
	}
	tx2 := Tx{
		FromID: "Bill",  //交易发起人
		ToID:   "Aaron", //交易的对方
		Value:  2000,    //交易的金额
	}
	//我们需要一个私钥 路径是我们的账户的路径
	path := "zblock/accounts/kennedy.ecdsa"
	privateKey, err := crypto.LoadECDSA(path)
	if err != nil {
		return fmt.Errorf("unable to load private key for node: %w", err)
	}
	data1, err := json.Marshal(tx1)
	if err != nil {
		return fmt.Errorf("unable to Marshal: %w", err)
	}
	data2, err := json.Marshal(tx2)
	if err != nil {
		return fmt.Errorf("unable to Marshal: %w", err)
	}

	v1 := crypto.Keccak256(data1)
	sig1, err := crypto.Sign(v1, privateKey)
	if err != nil {
		return fmt.Errorf("unable to Sign: %w", err)
	}

	v2 := crypto.Keccak256(data2)
	sig2, err := crypto.Sign(v2, privateKey)
	if err != nil {
		return fmt.Errorf("unable to Sign: %w", err)
	}
	fmt.Println("这是tx1 的签名" + string(hexutil.Encode(sig1)))
	fmt.Println("这是tx2 的签名" + string(hexutil.Encode(sig2)))

	// 提取公钥 使用你的私钥对数据进行签名时，签名操作实际上已经使用了私钥的椭圆曲线参数。从签名中可以提取出公钥。
	publicKey1, err := crypto.Ecrecover(v1, sig1)
	if err != nil {
		return fmt.Errorf("unable to recover public key for tx1: %w", err)
	}

	publicKey2, err := crypto.Ecrecover(v2, sig2)
	if err != nil {
		return fmt.Errorf("unable to recover public key for tx2: %w", err)
	}

	fmt.Println("这是tx1 的公钥" + string(hexutil.Encode(publicKey1)))
	fmt.Println("这是tx2 的公钥" + string(hexutil.Encode(publicKey2)))

	if bytes.Compare(publicKey1, publicKey2) == 0 {
		fmt.Println("两次签名使用的公钥是相同的。")
	} else {
		fmt.Println("两次签名使用的公钥不同。")
	}

	return nil
}

// 写一个 把正常的16进制签名转换成符合以太坊的签名的method
func ToVRSFromHexSig(sigStr string) (v, r, s *big.Int, err error) {
	// 首先解码十六进制签名
	// 去除开头的 "0x" 前缀（如果有的话）
	hexStr := strings.TrimPrefix(sigStr, "0x")
	sig, err := hex.DecodeString(hexStr)
	//sig, err := hex.DecodeString(sigStr[2:]) // 如果用这个的话，如果我传入的是没有带0x的16进制数 他就会去除掉有效数据了
	if err != nil {
		return nil, nil, nil, err
	}
	// 确保解码后的签名长度为 65 字节
	if len(sig) != 65 {
		return nil, nil, nil, fmt.Errorf("签名长度错误: %d", len(sig))
	}

	// 提取 r 和 s 部分
	r = new(big.Int).SetBytes(sig[0:32])
	s = new(big.Int).SetBytes(sig[32:64])

	// 提取 v 部分（最后一个字节）
	v = new(big.Int).SetUint64(uint64(sig[64]))

	return v, r, s, nil
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

// ToSignatureBytesWithArdanID converts the r, s, v values into a slice of bytes
// keeping the Ardan id.
func ToSignatureBytesWithArdanID(v, r, s *big.Int) []byte {
	sig := ToSignatureBytes(v, r, s)
	sig[64] = byte(v.Uint64())

	return sig
}
