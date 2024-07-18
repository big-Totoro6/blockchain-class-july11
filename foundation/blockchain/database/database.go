// Package database handles all the lower level support for maintaining the
// blockchain in storage and maintaining an in-memory databse of account information.
// database 包处理所有维护区块链存储和内存数据库账户信息的底层支持。
// database 包提供以下功能：
// - 持久化存储区块链数据。
// - 在内存数据库中管理账户及其信息。
package database

import (
	"errors"
	"github.com/ardanlabs/blockchain/foundation/blockchain/genesis"
	"sync"
)

// Database 需要一个数据库的类型 里面有互斥锁   基础配置 以及账户信息 是一个map  通过账户的id 也就是以太坊地址 去关联 完整的账户信息
// Database manages data related to accounts who have transacted on the blockchain.
type Database struct {
	mu       sync.RWMutex
	genesis  genesis.Genesis
	accounts map[AccountID]Account
}

// New 需要一个工程函数 去构建这个数据库 他是指针传递，意味着我们不想它复制太多，有一个实例即可
// New evHandler 他是一个事件函数，因为我们不想把它跟这个工厂函数强绑定，通过这种方式，由调用者自己定义自己想要的事件处理函数（比如log 或者什么 自己自定义）这样更灵活
// New constructs a new database and applies account genesis information and
// reads/writes the blockchain database on disk if a dbPath is provided.
func New(genesis genesis.Genesis, evHandler func(v string, args ...any)) (*Database, error) {
	db := Database{
		genesis:  genesis,
		accounts: make(map[AccountID]Account),
	} //注意这里的make 不先开辟空间 切片会抛异常的

	// 把初始配置里的账户信息 更新到数据库里面
	// Update the database with account balance information from genesis.
	for accountStr, balance := range genesis.Balances {
		accountID, err := ToAccountID(accountStr)
		if err != nil {
			return nil, err
		}
		db.accounts[accountID] = newAccount(accountID, balance)
	}
	return &db, nil
}

// Remove deletes an account from the database.
func (db *Database) Remove(accountID AccountID) {
	db.mu.Lock()
	defer db.mu.Unlock()

	delete(db.accounts, accountID)
}

// Query retrieves an account from the database.
func (db *Database) Query(accountID AccountID) (Account, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	acount, exists := db.accounts[accountID]
	if !exists {
		return Account{}, errors.New("account does not exist")
	}

	return acount, nil
}

// Copy makes a copy of the current accounts in the database.
func (db *Database) Copy() map[AccountID]Account {
	db.mu.RLock()
	defer db.mu.RUnlock()

	accounts := make(map[AccountID]Account)
	for accountID, account := range db.accounts {
		accounts[accountID] = account
	}
	return accounts
}
