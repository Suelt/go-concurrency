package blockchain

import (
	"UndergraduateProj/Hasher"
	"UndergraduateProj/block"
	"UndergraduateProj/config"
	"UndergraduateProj/receipt"
	"UndergraduateProj/tx"
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"github.com/pingcap/go-ycsb/pkg/generator"
	"log"
	"math"
	"strconv"
)

type BcProposer struct {
	bc     *Blockchain
	cfg    *config.ProposerConfig
	engine Engine
}

func NewBcProposer(bc *Blockchain) *BcProposer {
	return &BcProposer{
		bc:     bc,
		cfg:    &config.ProposerConfig{generator.NewZipfianWithRange(0, tx.Smallbanksize-1, tx.Zipfianconstant)},
		engine: bc.Engine,
	}
}
func (bp *BcProposer) Mine(b *block.Block) (uint64, error) {
	b.Header.Coinbase = [20]byte{192, 168, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}
	b.Header.Blockroot = bp.bc.Blockroot
	b.Header.MixDigest = [32]byte{}
	b.Tree = bp.bc.Tree
	var nonce uint64
	var diffbytes []byte
	if b.Header.Difficulty != nil {
		diffbytes = b.Header.Difficulty.Bytes()
	}

	s := [][]byte{
		b.Header.ArbitraryData[:], tx.Uint64toByte(b.Header.GasLimit), tx.Uint64toByte(b.Header.GasUsed),
		b.Header.TxRoot[:], b.Header.ReceiptRoot[:], b.Header.StateRoot[:], diffbytes,
		b.Header.Blockroot[:], []byte(strconv.FormatInt(b.Header.Time, 10)), b.Header.Coinbase[:],
		b.Header.MixDigest[:],
	}
	bytecopy := bytes.Join(s, []byte{})
	for i := 0; i < 8; i++ {
		bytecopy = append(bytecopy, 0)
	}
	length := len(bytecopy)
	//fmt.Println(length)
	//var maxsum uint64
	for nonce = 0; nonce < math.MaxUint64; nonce++ {
		n := b.Header.EncodeNonce(nonce)
		b.Header.Nonce = n
		for i := 0; i < 8; i++ {
			bytecopy[length-8+i] = n[i]
		}

		h := sha256.New()

		if _, err := h.Write(bytecopy); err != nil {
			log.Panic("Mining error:", err)
		}

		hash := h.Sum(nil)
		zeros := ComLeadingzeros(hash)
		//if zeros > maxsum {
		//	maxsum = zeros
		//	fmt.Printf("Now maxleadingzeros %d", maxsum)
		//	fmt.Printf("\n\n\n\n\n\n\n\n")
		//}
		if zeros >= config.LeadingZeros {

			//fmt.Printf("mining process succeeds,got nonce %d\n", nonce)
			hash32 := [32]byte{}
			copy(hash32[:], hash)
			b.Ethash = hash32
			//for i := 0; i < 32; i++ {
			//	fmt.Printf("%d", hash32[i])
			//
			//}
			//fmt.Printf("Ethash %x", hash32)
			//fmt.Println()
			return nonce, nil
		}
	}

	return 0, nil
}

func (bp *BcProposer) ProposeBlock() *block.Block {
	curve := elliptic.P256()
	privkey, _ := ecdsa.GenerateKey(curve, rand.Reader)
	txs := bp.bc.Sb.GenerateBatchSmallbankTx(bp.cfg.Zipfian, bp.cfg, privkey)
	//	fmt.Println("Transaction Generation Completed")
	//t1 := time.Now()
	env := bp.GenerateReceiptsAndStateRoot(txs)
	//t2 := time.Since(t1)
	//fmt.Println(t2)

	var (
		receipts  = env.Receipts
		stateroot = env.Stateroot
		gasused   = env.H.GasUsed
	)

	h := &block.Header{
		Difficulty:    bp.bc.Difficulty,
		GasLimit:      4200000 * config.Blocksize,
		ArbitraryData: []byte("Blockpacked"),
		Blocknumber:   receipts[0].BlockNumber,
		GasUsed:       gasused,
		StateRoot:     stateroot,
	}

	//fmt.Println(env.H.GasUsed)
	packBlock := bp.PackBlock(h, txs, receipts, Hasher.NewNormalHasher())
	//fmt.Println("Block packed")
	//packBlock.Header.Time = time.Now().Unix()
	return packBlock
}
func (bp *BcProposer) PackBlock(h *block.Header, txs tx.Transactions, receipts receipt.Receipts, trie Hasher.Trie) *block.Block {

	return block.NewBlock(h, txs, receipts, trie)
}

func (bp *BcProposer) GenerateReceiptsAndStateRoot(txs tx.Transactions) *ExecutionEnv {

	txsgas := ReArrangeTxsbyGas(txs)
	//Todo After Genesisblock
	blockhash := bp.bc.Blockroot
	blocknumber := bp.bc.FetchBlocknumberSum() + 1
	h := &block.Header{
		Blockroot:   blockhash,
		Blocknumber: blocknumber,
	}
	bc := bp.bc

	sbcopy := bc.Sb.Duplicate()
	env := bp.bc.MakeEnv(h, sbcopy)
	env.Txs = make(tx.Transactions, len(txs))
	env.Receipts = make(receipt.Receipts, len(txs))
	env.Sets = make(tx.ConcurrentSet, len(txs))
	//root := sbcopy.GenerateStateRoot()
	//env.Stateroot = root
	//fmt.Printf("%x\n", root)
	for index, transaction := range txsgas {
		err := env.ExecuteTransaction(transaction, uint64(index))
		//root := sbcopy.GenerateStateRoot()
		//env.Stateroot = root
		//fmt.Printf("%x\n", root)
		if err != nil {
			log.Panicln(err)
		}

	}
	root := sbcopy.GenerateStateRoot()
	env.Stateroot = root
	//sbcopy.CommittoState(env.Sets)

	//fmt.Println(receipts[1].GasUsed)
	return env
}

func ReArrangeTxsbyGas(txs []*tx.Transaction) []*tx.Transaction {
	//Todo
	return txs
}
