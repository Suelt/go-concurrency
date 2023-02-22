package blockchain

import (
	"UndergraduateProj/Hasher"
	"UndergraduateProj/block"
	cfg "UndergraduateProj/config"
	"UndergraduateProj/receipt"
	"UndergraduateProj/tx"
	"crypto/rand"
	"fmt"
	lru "github.com/hashicorp/golang-lru"
	"github.com/syndtr/goleveldb/leveldb"
	"log"
	"math"
	"math/big"
	"time"
)

const chainnumber = 32

type Blockchain struct {
	Sb               tx.Smallbankop
	StateProcessor   Processor
	ChainValidator   Validator
	BlockProposer    Proposer
	Chains           [cfg.SinglechainSize]*Singlechain
	Difficulty       *big.Int
	Engine           Engine
	Blockcache       *lru.Cache
	Transactioncache *lru.Cache
	Headercache      *lru.Cache
	Blockroot        [32]byte
	ExecutionMode    bool
	Tree             *Hasher.MerkleTree
	DBNormal         *leveldb.DB
	DBChain          *leveldb.DB
	DBState          *leveldb.DB
	StateTrie        *Hasher.MerkleTree
}

type Singlechain struct {
	Genesisblock *block.Block
	Currentblock *block.Block
	Chain        []*block.Block
	cfg          cfg.ChainConfig
	isnil        bool
}

func (bc *Blockchain) ProcessBlock(block *block.Block) (receipt.Receipts, uint64, error, time.Duration, int64) {

	return bc.StateProcessor.Process(block)

}

func InitBlockchain() *Blockchain {
	DBlist, err := LoadDB()
	if err != nil {
		log.Panicln(err)
	}
	data, errnew := DBlist[0].Get([]byte{1, 2, 3, 4, 5, 6, 7, 8}, nil)
	cfg.UseRw(data)
	if errnew != nil {
		return NewBlockchain(DBlist)
	}
	cfg.UseRw(errnew)
	DBNormal := DBlist[0]
	DBChain := DBlist[1]
	DBState := DBlist[2]
	chain, errre := RegenerateBlockchainViaDB(DBChain)
	if errre != nil {
		log.Panicln(errre)
	}

	chain.DBChain = DBChain
	chain.DBNormal = DBNormal
	chain.DBState = DBState
	sb, errstate := RegenerateStateViaDB(DBState)
	if errstate != nil {
		log.Panicln(errstate)
	}

	chain.Sb = sb
	chain.StateTrie = chain.Sb.GenerateStatetrie()
	if err != nil {
		log.Panicln(err)
	}
	return chain
}
func NewBlockchain(DBlist []*leveldb.DB) *Blockchain {
	max := new(big.Int)
	max.Exp(big.NewInt(2), big.NewInt(70), nil).Sub(max, big.NewInt(1))
	n, _ := rand.Int(rand.Reader, max)

	var sb tx.Smallbank
	sb.Init()
	blockchain := &Blockchain{
		Sb: &sb,
	}

	blockchain.Difficulty = n
	blockchain.Engine = &Chainworker{Bc: blockchain}
	blockchain.BlockProposer = NewBcProposer(blockchain)
	blockchain.StateProcessor = NewConcurrentProcessor(blockchain)
	blockchain.ChainValidator = NewBcValidator(blockchain)
	position := &block.Blockposition{
		List: 0,
	}
	for i := 0; i < 32; i++ {
		blockchain.Chains[i] = NewSinglechain(uint64(i))
		blockchain.Chains[i].Chain[0].Header.Difficulty = blockchain.Difficulty
		position.Line = i
		blockchain.Chains[i].Chain[0].Position = position
		hash, _ := blockchain.Chains[i].Chain[0].HashBlock()
		blockchain.Chains[i].Chain[0].Hash = hash

	}
	blockchain.GenerateBlockRoot()
	blockchain.Tree = blockchain.GenerateBlockMerkletree()
	blockchain.Sb.GenerateStateRoot()
	blockchain.StateTrie = blockchain.Sb.GenerateStatetrie()
	blockchain.DBNormal = DBlist[0]
	blockchain.DBChain = DBlist[1]
	blockchain.DBState = DBlist[2]

	errput := blockchain.DBNormal.Put([]byte{1, 2, 3, 4, 5, 6, 7, 8}, []byte{1, 2}, nil)
	if errput != nil {
		log.Panicln(errput)
	}
	err := StoreWholeState(blockchain.DBState, blockchain.Sb.(*tx.Smallbank))
	if err != nil {
		log.Panicln(err)
	}
	return blockchain
}

func NewSinglechain(chainid uint64) *Singlechain {
	configuration := cfg.ChainConfig{
		ChainId: chainid,
	}
	single := NewScWithGenesis(block.NewGenesisBlock())
	single.isnil = false
	single.cfg = configuration
	return single
}

func NewScWithGenesis(genesis *block.Block) *Singlechain {

	single := Singlechain{}
	single.Genesisblock = genesis
	single.Currentblock = genesis
	single.Chain = append(single.Chain, genesis)
	return &single
}

func (bc *Blockchain) Mine() *block.Block {
	proposedblock := bc.BlockProposer.ProposeBlock()
	nonce, err := bc.BlockProposer.Mine(proposedblock)

	if err != nil {
		log.Panic("Mining process error!", err)
	}

	encodednonce := proposedblock.Header.EncodeNonce(nonce)
	proposedblock.Header.Nonce = encodednonce
	hash, _ := proposedblock.CalculateEthash()
	hash32 := [32]byte{}
	copy(hash32[:], hash)
	//copy(proposedblock.Hash[:], hash)
	line := GetBlockLine(hash32)
	list := len(bc.Chains[line].Chain)
	position := &block.Blockposition{
		Line: int(line),
		List: list,
	}
	//position and parenthash needs to be removed
	proposedblock.Position = position
	proposedblock.Parenthash = bc.Chains[line].Chain[list-1].Hash

	return proposedblock
}

func (bc *Blockchain) VerifyBlock(UnvalidatedBlock *block.Block) error {
	err := bc.Engine.Verifyheader(UnvalidatedBlock)
	if err == nil {
		err = bc.ChainValidator.ValidateBody(UnvalidatedBlock)
	}
	if err != nil {
		return err
	}
	return nil
}
func (bc *Blockchain) PendBlock(b *block.Block) (error, time.Duration, int64) {

	err := bc.VerifyBlock(b)
	if err != nil {
		return err, 0, 0
	}

	receipts, gasused, err, tret, tDete := bc.ProcessBlock(b)
	if err != nil {
		return err, 0, 0
	}
	cfg.UseRw(receipts, gasused)
	fmt.Println(len(receipts))
	err = bc.ChainValidator.ValidateState(b, receipts, gasused)
	if err != nil {
		return err, 0, 0
	}

	var (
		line = b.Position.Line
	)

	bc.Chains[line].Chain = append(bc.Chains[line].Chain, b)
	bc.Chains[line].Currentblock = b
	fmt.Printf("Block maps to chain %d,position %d\n", line, len(bc.Chains[line].Chain))
	bc.GenerateBlockRoot()
	bc.Tree = bc.GenerateBlockMerkletree()

	//errstoreblock := StoreBlock(bc.DBChain, b)
	//if errstoreblock != nil {
	//	log.Println(errstoreblock)
	//}
	//
	//errstoreheader := StoreBlockHeader(bc.DBNormal, b)
	//if errstoreheader != nil {
	//	log.Println(errstoreheader)
	//}
	return nil, tret, tDete
}

//assemble block
func (bc *Blockchain) GenerateMainWorkFlow() {
	a := bc.FetchBlocknumberSum()
	var tsum time.Duration
	var tdetesum int64
	for i := 0; i < 50; i++ {

		fmt.Printf("block%d\n", a)
		//time.Sleep(1000 * time.Millisecond)
		b := bc.Mine()
		err, tret, tDete := bc.PendBlock(b)
		tsum += tret
		tdetesum += tDete
		log.Println(tdetesum)
		log.Println("end")
		//time.Sleep(1000 * time.Millisecond)
		a++

		if err != nil {
			log.Panic(err)
		}
	}
	fmt.Println()
	fmt.Println()
	fmt.Println()
	fmt.Println("total time consumption:")
	fmt.Println(tsum)
	//fmt.Println(tdetesum)
}

func GetBlockLine(bytes [32]byte) uint64 {
	b := bytes[31]
	bytearray := new([8]byte)
	var res uint64
	var multi uint64
	for i := 5; i > 0; i-- {
		bytearray[i] = b & 1
		b = b >> 1
		multi = uint64(bytearray[i])
		res += multi * uint64(math.Pow(float64(2), float64(5-i)))
	}
	return res
}

//Generate

func (bc *Blockchain) GenerateBlockRoot() {

	blockcontents := make([]Hasher.Content, cfg.SinglechainSize)
	for index, value := range bc.Chains {
		blockcontents[index] = *value.Currentblock

	}
	bc.Blockroot = Hasher.NewNormalHasher().DeriveHash(&blockcontents)

}

func (bc *Blockchain) GenerateBlockMerkletree() *Hasher.MerkleTree {
	blockcontents := make([]Hasher.Content, cfg.SinglechainSize)
	for index, value := range bc.Chains {
		blockcontents[index] = *value.Currentblock

	}
	t, err := Hasher.NewTree(blockcontents)
	if err != nil {
		log.Fatal(err)
	}
	return t
}
func (bc *Blockchain) FetchBlocknumberSum() uint64 {
	var Sum uint64
	for _, chain := range bc.Chains {
		if chain.isnil == false {
			Sum += uint64(len(chain.Chain))
		}
	}
	return Sum - 32
}

func (bc *Blockchain) MakeEnv(h *block.Header, sb tx.Smallbankop) *ExecutionEnv {

	env := &ExecutionEnv{
		sb: sb,
		H:  h,
	}
	return env
}

func (bc *Blockchain) MakeEnvWithBlock(b *block.Block, sb tx.Smallbankop) *ExecutionEnv {
	h := &block.Header{
		Blockroot:   b.Hash,
		Blocknumber: b.Header.Blocknumber,
	}
	env := bc.MakeEnv(h, sb)
	txslength := len(b.Transactions)
	env.Txs = make(tx.Transactions, txslength)
	env.Receipts = make(receipt.Receipts, txslength)
	env.Sets = make(tx.ConcurrentSet, txslength)
	return env
}

func (bc *Blockchain) BatchBlockProcessingWorkflow() {

}
