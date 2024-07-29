package public

import "github.com/ardanlabs/blockchain/foundation/blockchain/database"

// 应用层的model 类似于dto 不想把业务层的model字段原原本本全部返回暴露
type tx struct {
	FromAccount database.AccountID `json:"from"`
	To          database.AccountID `json:"to"`
	FromName    string             `json:"from_name"`
	ToName      string             `json:"to_name"`
	ChainID     uint16             `json:"chain_id"`
	Nonce       uint64             `json:"nonce"`
	Value       uint64             `json:"value"`
	Tip         uint64             `json:"tip"`
	Data        []byte             `json:"data"`
	TimeStamp   uint64             `json:"timestamp"`
	GasPrice    uint64             `json:"gas_price"`
	GasUnits    uint64             `json:"gas_units"`
	Sig         string             `json:"sig"`
}
