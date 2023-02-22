package blockchain

import (
	"UndergraduateProj/Hasher"
	"UndergraduateProj/block"
	cfg "UndergraduateProj/config"
	"UndergraduateProj/receipt"
	"fmt"
)

type BcValidator struct {
	bc     *Blockchain
	cfg    cfg.ValidatorConfig
	engine Engine
}

func NewBcValidator(bc *Blockchain) *BcValidator {
	return &BcValidator{
		bc,
		cfg.ValidatorConfig{},
		bc.Engine,
	}
}

func (validator *BcValidator) ValidateBody(block *block.Block) error {
	if leadingzeros := ComLeadingzeros(block.Ethash[:]); leadingzeros < cfg.LeadingZeros {
		return fmt.Errorf("leadingzeros mismatched,got %d,expected %d", leadingzeros, cfg.LeadingZeros)
	}
	fmt.Printf("Mining validation passed,got hash:%x,has %d leadingzeros\n", block.Ethash, 8*cfg.LeadingZeros)
	hash, _ := block.CalculateEthash()
	var hash32 [32]byte
	copy(hash32[:], hash)
	if hash32 != block.Ethash {
		return fmt.Errorf("blockhash unmatched")
	}
	//blockroot := validator.bc.Blockroot
	//if blockroot != block.Header.Blockroot {
	//	return fmt.Errorf("blockroot mismatched,got blockroot %x,expected %x", blockroot, block.Header.Blockroot)
	//}
	err := validator.ValidateBlockroot(block)
	if err != nil {
		return err
	}
	fmt.Printf("Blockroot validation passed\n")
	return nil
}

func (validator *BcValidator) ValidateState(b *block.Block, receipts receipt.Receipts, usedgas uint64) error {
	header := b.Header
	if header.GasUsed != usedgas {
		return fmt.Errorf("invalid gasused,locally processed %d,received %d", header.GasUsed, usedgas)
	}
	fmt.Printf("gasused Validation passed,totalgas used:%d\n", usedgas)
	if len(receipts) != len(b.Transactions) {
		return fmt.Errorf("got %d receipts but %d transactions", len(receipts), len(b.Transactions))
	}

	hasher := Hasher.NewNormalHasher()
	ReceiptContents := make([]Hasher.Content, len(receipts))
	//fmt.Println(len(ReceiptContents))
	for i := 0; i < len(receipts); i++ {
		ReceiptContents[i] = *receipts[i]
	}
	receiptroot := hasher.DeriveHash(&ReceiptContents)
	if receiptroot != b.Header.ReceiptRoot {
		return fmt.Errorf("Receiptroot mismatched,locally processed %x,received %x", receiptroot, b.Header.ReceiptRoot)
	}
	fmt.Printf("Receiptroot validation passed,Receiptroot:%x\n", receiptroot)
	stateroot := validator.bc.StateTrie.MerkleRoot()
	var root32 [32]byte
	copy(root32[:], stateroot)
	//fmt.Printf("%x\n", stateroot)
	if root32 != b.Header.StateRoot {
		return fmt.Errorf("Stateroot mismatched,locally processed %x,received %x", stateroot, b.Header.StateRoot)
	}

	fmt.Printf("Stateroot validation passed,Stateroot:%x\n", stateroot)

	return nil
}

func ComLeadingzeros(h []byte) uint64 {
	for i := 0; i < len(h); i++ {
		if h[i] != 0 {
			return uint64(i)
		}
	}
	return uint64(len(h))
}

func (Validator *BcValidator) ValidateBlockroot(b *block.Block) error {
	var (
		tree     = b.Tree
		position = b.Position
	)
	concatingblock := Validator.bc.Chains[position.Line].Currentblock
	var content Hasher.Content
	content = *concatingblock
	ifstillthatblock, _ := tree.VerifyContent(content)
	if ifstillthatblock {
		fmt.Printf("Current block at line %d is %d\n", position.Line,
			len(Validator.bc.Chains[position.Line].Chain))
		return nil
	} else {
		return fmt.Errorf("Current block at line %d is %d,but received %d\n", position.Line,
			len(Validator.bc.Chains[position.Line].Chain), position.List)
	}
}
