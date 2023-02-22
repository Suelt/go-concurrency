package blockchain

import (
	"UndergraduateProj/block"
	"UndergraduateProj/receipt"
	"UndergraduateProj/tx"
	"math/big"
	"time"
)

//type Receipt struct {
//
//	// Implementation fields: These fields are added by geth when processing a transaction.
//	// They are stored in the chain database.
//	TxHash          [32]byte
//	ContractAddress [32]byte
//	GasUsed         uint64
//
//	// Inclusion information: These fields provide information about the inclusion of the
//	// transaction corresponding to this receipt.
//	BlockHash        [32]byte
//	BlockNumber      int64
//	TransactionIndex uint
//	StateRoot        [64]byte
//}
//type Receipts []Receipt

type Validator interface {
	ValidateBody(b *block.Block) error
	//ValidateOhie(bc *Blockchain) error
	ValidateState(b *block.Block, receipts receipt.Receipts, usedgas uint64) error
}

type Proposer interface {
	Mine(b *block.Block) (uint64, error)
	//GenerateReceiptsAndStateRoot(txs tx.Transactions) (receipt.Receipts, [32]byte, uint64)
	ProposeBlock() *block.Block
	//Multicast()
}

type Processor interface {
	Process(block *block.Block) (receipt.Receipts, uint64, error, time.Duration, int64)
}

type Engine interface {
	Verifyheader(b *block.Block) error
	Calcdifficulty(time int64, parent *block.Header) *big.Int
	GetBlock(hash [32]byte, position *block.Blockposition) *block.Block
	GetHeader(hash [32]byte, position *block.Blockposition) *block.Header
	GetTransactionCount(blockHash [32]byte) uint64
	GetTransactionInBlock(blockHash [32]byte, index uint64) *tx.Transaction
	GetBlockByHash(hash [32]byte) *block.Block

	GetBlockByPosition(position *block.Blockposition) *block.Block
	GetHeaderByHash(hash [32]byte) *block.Header

	GetHeaderByPosition(position *block.Blockposition) *block.Header
}
