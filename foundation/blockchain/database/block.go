package database

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"math"
	"math/big"
	"time"

	"github.com/ardanlabs/blockchain/foundation/blockchain/merkle"
	"github.com/ardanlabs/blockchain/foundation/blockchain/signature"
)

// ErrChainForked is returned from validateNextBlock if another node's chain
// is two or more blocks ahead of ours.
var ErrChainForked = errors.New("blockchain forked, start resync")

// =============================================================================

// BlockData represents what can be serialized to disk and over the network.
type BlockData struct {
	Hash   string      `json:"hash"`
	Header BlockHeader `json:"block"`
	Trans  []BlockTx   `json:"trans"`
}

// NewBlockData constructs block data from a block.
func NewBlockData(block Block) BlockData {
	blockData := BlockData{
		Hash:   block.Hash(),
		Header: block.Header,
		Trans:  block.MerkleTree.Values(),
	}

	return blockData
}

// ToBlock converts a storage block into a database block.
func ToBlock(blockData BlockData) (Block, error) {
	tree, err := merkle.NewTree(blockData.Trans)
	if err != nil {
		return Block{}, err
	}

	block := Block{
		Header:     blockData.Header,
		MerkleTree: tree,
	}

	return block, nil
}

// =============================================================================

// BlockHeader represents common information required for each block.
type BlockHeader struct {
	Number        uint64    `json:"number"`          // Ethereum: Block number in the chain. 区块在链中的编号。
	PrevBlockHash string    `json:"prev_block_hash"` // Bitcoin: Hash of the previous block in the chain. 前一个区块的哈希值。
	TimeStamp     uint64    `json:"timestamp"`       // Bitcoin: Time the block was mined. 区块被挖掘的时间。
	BeneficiaryID AccountID `json:"beneficiary"`     // Ethereum: The account who is receiving fees and tips. 收取手续费和小费的账户。
	Difficulty    uint16    `json:"difficulty"`      // Ethereum: Number of 0's needed to solve the hash solution. 解决哈希函数所需的前导零的数量。
	MiningReward  uint64    `json:"mining_reward"`   // Ethereum: The reward for mining this block. 挖掘此区块的奖励。
	StateRoot     string    `json:"state_root"`      // Ethereum: Represents a hash of the accounts and their balances. 表示账户及其余额的哈希值。
	TransRoot     string    `json:"trans_root"`      // Both: Represents the merkle tree root hash for the transactions in this block. 表示此区块中交易的默克尔树根哈希。
	Nonce         uint64    `json:"nonce"`           // Both: Value identified to solve the hash solution. 解决哈希函数的数值标识。
}

// Block represents a group of transactions batched together.
type Block struct {
	Header     BlockHeader
	MerkleTree *merkle.Tree[BlockTx]
}

// POWArgs represents the set of arguments required to run POW.
type POWArgs struct {
	BeneficiaryID AccountID                   // 受益者ID，指定接收挖矿奖励的账户ID
	Difficulty    uint16                      // 难度，表示挖矿过程中要求的工作量难度
	MiningReward  uint64                      // 挖矿奖励，表示成功挖到新区块后给予矿工的奖励数量
	PrevBlock     Block                       // 前一个区块，包含了前一个区块的信息，用于构建新区块的前置条件
	StateRoot     string                      // 状态根，表示当前区块链状态的根哈希值，用于构建新区块时的状态校验
	Trans         []BlockTx                   // 交易列表，包含了新区块要包含的所有交易
	EvHandler     func(v string, args ...any) // 事件处理器，用于处理挖矿过程中的各种事件
}

// POW constructs a new Block and performs the work to find a nonce that
// solves the cryptographic POW puzzel.
func POW(ctx context.Context, args POWArgs) (Block, error) {

	// When mining the first block, the previous block's hash will be zero.
	prevBlockHash := signature.ZeroHash
	if args.PrevBlock.Header.Number > 0 {
		prevBlockHash = args.PrevBlock.Hash()
	}

	// Construct a merkle tree from the transaction for this block. The root
	// of this tree will be part of the block to be mined.
	tree, err := merkle.NewTree(args.Trans)
	if err != nil {
		return Block{}, err
	}

	// Construct the block to be mined.
	block := Block{
		Header: BlockHeader{
			Number:        args.PrevBlock.Header.Number + 1,
			PrevBlockHash: prevBlockHash,
			TimeStamp:     uint64(time.Now().UTC().UnixMilli()),
			BeneficiaryID: args.BeneficiaryID,
			Difficulty:    args.Difficulty,
			MiningReward:  args.MiningReward,
			StateRoot:     args.StateRoot,
			TransRoot:     tree.RootHex(), //
			Nonce:         0,              // Will be identified by the POW algorithm.
		},
		MerkleTree: tree,
	}

	// Peform the proof of work mining operation.
	if err := block.performPOW(ctx, args.EvHandler); err != nil {
		return Block{}, err
	}

	return block, nil
}

// performPOW does the work of mining to find a valid hash for a specified
// block. Pointer semantics are being used since a nonce is being discovered.
func (b *Block) performPOW(ctx context.Context, ev func(v string, args ...any)) error {
	ev("database: PerformPOW: MINING: started")
	defer ev("database: PerformPOW: MINING: completed")

	// Log the transactions that are a part of this potential block.
	for _, tx := range b.MerkleTree.Values() {
		ev("database: PerformPOW: MINING: tx[%s]", tx)
	}

	// Choose a random starting point for the nonce. After this, the nonce
	// will be incremented by 1 until a solution is found by us or another node.
	nBig, err := rand.Int(rand.Reader, big.NewInt(math.MaxInt64))
	if err != nil {
		return ctx.Err()
	}
	b.Header.Nonce = nBig.Uint64()

	ev("viewer: PerformPOW: MINING: running")

	// Loop until we or another node finds a solution for the next block.
	var attempts uint64
	for {
		attempts++
		if attempts%1_000_000 == 0 {
			ev("viewer: PerformPOW: MINING: running: attempts[%d]", attempts)
		}

		// Did we timeout trying to solve the problem.
		if ctx.Err() != nil {
			ev("database: PerformPOW: MINING: CANCELLED")
			return ctx.Err()
		}

		// Hash the block and check if we have solved the puzzle.
		hash := b.Hash()
		if !isHashSolved(b.Header.Difficulty, hash) {
			b.Header.Nonce++
			//记住这里是没找到合适的 他会回到上面for循环的开头 当contine之后
			continue
		}
		//执行到这里就代表找到了需要的hash
		//if you have solution; go for that you do not need sb tell you. you should cancel then you cancel
		ev("database: PerformPOW: MINING: SOLVED: prevBlk[%s]: newBlk[%s]", b.Header.PrevBlockHash, hash)
		ev("database: PerformPOW: MINING: attempts[%d]", attempts)

		return nil
	}
}

// Hash returns the unique hash for the Block.
func (b Block) Hash() string {
	if b.Header.Number == 0 {
		return signature.ZeroHash
	}

	// CORE NOTE: Hashing the block header and not the whole block so the blockchain
	// can be cryptographically checked by only needing block headers and not full
	// blocks with the transaction data. This will support the ability to have pruned
	// nodes and light clients in the future.
	// - A pruned node stores all the block headers, but only a small number of full
	//   blocks (maybe the last 1000 blocks). This allows for full cryptographic
	//   validation of blocks and transactions without all the extra storage.
	// - A light client keeps block headers and just enough sufficient information
	//   to follow the latest set of blocks being produced. The do not validate
	//   blocks, but can prove a transaction is in a block.

	return signature.Hash(b.Header)
}

// ValidateBlock takes a block and validates it to be included into the blockchain.
func (b Block) ValidateBlock(previousBlock Block, stateRoot string, evHandler func(v string, args ...any)) error {
	evHandler("database: ValidateBlock: validate: blk[%d]: check: chain is not forked", b.Header.Number)

	// The node who sent this block has a chain that is two or more blocks ahead
	// of ours. This means there has been a fork and we are on the wrong side.
	nextNumber := previousBlock.Header.Number + 1
	if b.Header.Number >= (nextNumber + 2) {
		return ErrChainForked
	}

	evHandler("database: ValidateBlock: validate: blk[%d]: check: block difficulty is the same or greater than parent block difficulty", b.Header.Number)

	if b.Header.Difficulty < previousBlock.Header.Difficulty {
		return fmt.Errorf("block difficulty is less than previous block difficulty, parent %d, block %d", previousBlock.Header.Difficulty, b.Header.Difficulty)
	}

	evHandler("database: ValidateBlock: validate: blk[%d]: check: block hash has been solved", b.Header.Number)

	hash := b.Hash()
	if !isHashSolved(b.Header.Difficulty, hash) {
		return fmt.Errorf("%s invalid block hash", hash)
	}

	evHandler("database: ValidateBlock: validate: blk[%d]: check: block number is the next number", b.Header.Number)

	if b.Header.Number != nextNumber {
		return fmt.Errorf("this block is not the next number, got %d, exp %d", b.Header.Number, nextNumber)
	}

	evHandler("database: ValidateBlock: validate: blk[%d]: check: parent hash does match parent block", b.Header.Number)

	if b.Header.PrevBlockHash != previousBlock.Hash() {
		return fmt.Errorf("parent block hash doesn't match our known parent, got %s, exp %s", b.Header.PrevBlockHash, previousBlock.Hash())
	}

	if previousBlock.Header.TimeStamp > 0 {
		evHandler("database: ValidateBlock: validate: blk[%d]: check: block's timestamp is greater than parent block's timestamp", b.Header.Number)

		parentTime := time.Unix(int64(previousBlock.Header.TimeStamp), 0)
		blockTime := time.Unix(int64(b.Header.TimeStamp), 0)
		if blockTime.Before(parentTime) {
			return fmt.Errorf("block timestamp is before parent block, parent %s, block %s", parentTime, blockTime)
		}

		// This is a check that Ethereum does but we can't because we don't run all the time.

		// evHandler("database: ValidateBlock: validate: blk[%d]: check: block is less than 15 minutes apart from parent block", b.Header.Number)

		// dur := blockTime.Sub(parentTime)
		// if dur.Seconds() > time.Duration(15*time.Second).Seconds() {
		// 	return fmt.Errorf("block is older than 15 minutes, duration %v", dur)
		// }
	}

	evHandler("database: ValidateBlock: validate: blk[%d]: check: state root hash does match current database", b.Header.Number)

	if b.Header.StateRoot != stateRoot {
		return fmt.Errorf("state of the accounts are wrong, current %s, expected %s", stateRoot, b.Header.StateRoot)
	}

	evHandler("database: ValidateBlock: validate: blk[%d]: check: merkle root does match transactions", b.Header.Number)

	if b.Header.TransRoot != b.MerkleTree.RootHex() {
		return fmt.Errorf("merkle root does not match transactions, got %s, exp %s", b.MerkleTree.RootHex(), b.Header.TransRoot)
	}

	return nil
}

// isHashSolved checks the hash to make sure it complies with
// the POW rules. We need to match a difficulty number of 0's.
func isHashSolved(difficulty uint16, hash string) bool {
	const match = "0x00000000000000000"

	if len(hash) != 66 {
		return false
	}

	//也就是说 假如传入的difficulty是4  不管执行多少遍这个代码 最终比较的都是hash[:6] 比较前六位 相当于这是我们自己手动给传进来的难度值去调整难度
	difficulty += 2
	return hash[:difficulty] == match[:difficulty]
}
