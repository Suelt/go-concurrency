package blockchain

import (
	"UndergraduateProj/block"
	"UndergraduateProj/tx"
	"errors"
	"fmt"
	"math/big"
)

type Chainworker struct {
	Bc *Blockchain
}

func (worker *Chainworker) Verifyheader(b *block.Block) error {
	var err error
	header := b.Header
	err = header.CheckField()
	if err != nil {
		return err
	}
	parentposition := &block.Blockposition{
		Line: b.Position.Line,
		List: b.Position.List - 1,
	}
	parent := worker.GetHeader(b.Parenthash, parentposition)
	if parent == nil {
		return errors.New("parent block invalid")
	}
	//fmt.Printf("%x\n%x\n", header.Blockroot, parent.Blockroot)
	//fmt.Println(header.Time, parent.Time)
	//if header.Time <= parent.Time {
	//	return errors.New("older Block(timestamp earlier than the parent block)")
	//}
	//difficulty := worker.Calcdifficulty(header.Time, parent)
	//if difficulty.Cmp(header.Difficulty) != 0 {
	//	return errors.New("difficulty Mismatched with the computed difficulty")
	//}
	cap := uint64(0x7fffffffffffffff)
	if header.GasLimit > cap {
		return fmt.Errorf("invalid gaslimit %d,exceeds 2^63-1", header.GasLimit)
	}
	if header.GasUsed > header.GasLimit {
		return fmt.Errorf("gasused %d Exceeds Gaslimit %d", header.GasUsed, header.GasLimit)
	}
	fmt.Println("Block header validation passed")
	return nil
}

func (worker *Chainworker) Calcdifficulty(time int64, header *block.Header) *big.Int {
	return worker.Bc.Difficulty
}
func (worker *Chainworker) GetBlock(hash [32]byte, position *block.Blockposition) *block.Block {
	BlockByhash := worker.GetBlockByHash(hash)
	BlockBypos := worker.GetBlockByPosition(position)

	if BlockByhash.Hash == BlockBypos.Hash {
		return BlockBypos
	}
	return nil
}
func (worker *Chainworker) GetHeader(hash [32]byte, position *block.Blockposition) *block.Header {
	if block := worker.GetBlock(hash, position); block != nil {
		return block.Header
	}
	return nil
}
func (worker *Chainworker) GetTransactionCount(Hash [32]byte) uint64 {
	return 0
}
func (worker *Chainworker) GetTransactionInBlock(Hash [32]byte, index uint64) *tx.Transaction {
	return nil
}
func (worker *Chainworker) GetBlockByHash(hash [32]byte) *block.Block {
	chains := worker.Bc.Chains
	for _, chain := range chains {
		for _, block := range chain.Chain {
			if hash == block.Hash {
				return block
			}
		}
	}
	return nil
}

func (worker *Chainworker) GetBlockByPosition(position *block.Blockposition) *block.Block {
	return worker.Bc.Chains[position.Line].Chain[position.List]
}
func (worker *Chainworker) GetHeaderByHash(hash [32]byte) *block.Header {
	if block := worker.GetBlockByHash(hash); block != nil {
		return block.Header
	}
	return nil
}

func (worker *Chainworker) GetHeaderByPosition(position *block.Blockposition) *block.Header {
	if block := worker.GetBlockByPosition(position); block != nil {
		return block.Header
	}
	return nil
}
