package state

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ardanlabs/blockchain/foundation/blockchain/database"
)

// =============================================================================

// The set of different consensus protocols that can be used.
const (
	ConsensusPOW = "POW"
	ConsensusPOA = "POA"
)

// ErrNoTransactions is returned when a block is requested to be created
// and there are not enough transactions.
var ErrNoTransactions = errors.New("no transactions in mempool")

// =============================================================================

// MineNewBlock attempts to create a new block with a proper hash that can become
// the next block in the chain.
func (s *State) MineNewBlock(ctx context.Context) (database.Block, error) {
	defer s.evHandler("viewer: MineNewBlock: MINING: completed")

	s.evHandler("state: MineNewBlock: MINING: check mempool count")

	// Are there enough transactions in the pool.
	if s.mempool.Count() == 0 {
		return database.Block{}, ErrNoTransactions
	}

	// Pick the best transactions from the mempool.
	trans := s.mempool.PickBest(s.genesis.TransPerBlock)

	difficulty := s.genesis.Difficulty

	// Attempt to create a new block by solving the POW puzzle. This can be cancelled.
	block, err := database.POW(ctx, database.POWArgs{
		BeneficiaryID: s.beneficiaryID,
		Difficulty:    difficulty,
		MiningReward:  s.genesis.MiningReward,
		PrevBlock:     s.db.LatestBlock(),
		StateRoot:     s.db.HashState(),
		Trans:         trans,
		EvHandler:     s.evHandler,
	})
	if err != nil {
		return database.Block{}, err
	}

	// Just check one more time we were not cancelled.
	if ctx.Err() != nil {
		return database.Block{}, ctx.Err()
	}

	s.evHandler("state: MineNewBlock: MINING: validate and update database")

	// Validate the block and then update the blockchain database.
	if err := s.validateUpdateDatabase(block); err != nil {
		return database.Block{}, err
	}

	return block, nil
}

// =============================================================================

// validateUpdateDatabase takes the block and validates the block against the
// consensus rules. If the block passes, then the state of the node is updated
// including adding the block to disk.
func (s *State) validateUpdateDatabase(block database.Block) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.evHandler("state: validateUpdateDatabase: validate block")

	// CORE NOTE: I could add logic to determine if this block was mined by this
	// node or a peer. If the block is mined by this node, even if a peer beat
	// me to this function for the same block number, I could replace the peer
	// block with my own and attempt to have other peers accept my block instead.

	if err := block.ValidateBlock(s.db.LatestBlock(), s.db.HashState(), s.evHandler); err != nil {
		return err
	}

	s.evHandler("state: validateUpdateDatabase: write to disk")

	// Write the new block to the chain on disk.
	if err := s.db.Write(block); err != nil {
		return err
	}
	s.db.UpdateLatestBlock(block)

	s.evHandler("state: validateUpdateDatabase: update accounts and remove from mempool")

	//// Process the transactions and update the accounts.
	//for _, tx := range block.MerkleTree.Values() {
	//	s.evHandler("state: validateUpdateDatabase: tx[%s] update and remove", tx)
	//
	//	// Remove this transaction from the mempool.
	//	s.mempool.Delete(tx)
	//
	//	// Apply the balance changes based on this transaction.
	//	if err := s.db.ApplyTransaction(block, tx); err != nil {
	//		s.evHandler("state: validateUpdateDatabase: WARNING : %s", err)
	//		continue
	//	}
	//}
	//
	//s.evHandler("state: validateUpdateDatabase: apply mining reward")
	//
	//// Apply the mining reward for this block.
	//s.db.ApplyMiningReward(block)
	//
	//// Send an event about this new block.
	//s.blockEvent(block)

	return nil
}

// blockEvent provides a specific event about a new block in the chain for
// application specific support.
func (s *State) blockEvent(block database.Block) {
	data, err := json.Marshal(database.NewBlockData(block))
	if err != nil {
		data = []byte(fmt.Sprintf("{error: %q}", err.Error()))
	}

	s.evHandler("viewer: block: %s", string(data))
}
