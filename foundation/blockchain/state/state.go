// Package state is the core API for the blockchain and implements all the
// business rules and processing.
package state

import (
	"github.com/ardanlabs/blockchain/foundation/blockchain/database"
	"github.com/ardanlabs/blockchain/foundation/blockchain/genesis"
	"sync"
)

// EventHandler defines a function that is called when events 之前说的事件处理函数 用作日志记录的 在接口里面作为函数参数
// occur in the processing of persisting blocks.
type EventHandler func(v string, args ...any)

// =============================================================================

// Config represents the configuration required to start
// the blockchain node.
type Config struct {
	BeneficiaryID database.AccountID //受益人以太坊地址
	Genesis       genesis.Genesis
	EvHandler     EventHandler
}

// State manages the blockchain database.
type State struct {
	mu          sync.RWMutex
	resyncWG    sync.WaitGroup
	allowMining bool

	beneficiaryID database.AccountID
	evHandler     EventHandler

	genesis genesis.Genesis
	db      *database.Database
}

// New constructs a new blockchain for data management.
func New(cfg Config) (*State, error) {

	// Build a safe event handler function for use.
	ev := func(v string, args ...any) {
		if cfg.EvHandler != nil {
			cfg.EvHandler(v, args...)
		}
	}

	// Access the storage for the blockchain.
	db, err := database.New(cfg.Genesis, ev)
	if err != nil {
		return nil, err
	}

	// Create the State to provide support for managing the blockchain.
	state := State{
		beneficiaryID: cfg.BeneficiaryID,
		evHandler:     ev,
		allowMining:   true,

		genesis: cfg.Genesis,
		db:      db,
	}
	return &state, nil
}

// Shutdown cleanly brings the node down.
func (s *State) Shutdown() error {
	s.evHandler("state: shutdown: started")
	defer s.evHandler("state: shutdown: completed")

	// Wait for any resync to finish.
	s.resyncWG.Wait()
	return nil
}
