package main

import (
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"log"
)

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
}

func run() error {
	tx := Tx{
		FromID: "Bill",  //交易发起人
		ToID:   "Aaron", //交易的对方
		Value:  1000,    //交易的金额
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

	//输出
	//	�Օ��Su&V��
	//�F����s�����4�ք%$$o��,�������b��Hi����Ҕj����

	// =============================================================================
	//way 2 以太坊自带的哈希函数
	v := crypto.Keccak256(data)
	sig, err := crypto.Sign(v, privateKey)
	if err != nil {
		return fmt.Errorf("unable to Sign: %w", err)
	}
	//d%��#D3�Y��K�`�С�|Y1�]�F{�8���w�a�Cj��ؼ@)0@_H���ͯ��L

	// =============================================================================
	//然后打印一下看看签名
	fmt.Println(string(sig))
	//事实上这很难看懂 需要变成16进制的
	fmt.Println(string(hexutil.Encode(sig)))

	return nil

}
