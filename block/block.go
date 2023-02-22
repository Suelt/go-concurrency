package block

import (
	"UndergraduateProj/Hasher"
	cfg "UndergraduateProj/config"
	"UndergraduateProj/receipt"
	"UndergraduateProj/tx"
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/gob"
	"fmt"
	"log"
	"math/big"
	"strconv"
	"time"
)

type Block struct {
	Header *Header

	Transactions []*tx.Transaction

	// caches
	Hash         [32]byte
	Ethash       [32]byte
	Position     *Blockposition
	Parenthash   [32]byte
	ReceivedAt   uint64
	ReceivedFrom [32]byte
	Tree         *Hasher.MerkleTree
}
type BlockNonce [8]byte

type Header struct {
	Blockroot [32]byte

	Blocknumber uint64
	Coinbase    [20]byte //miner address
	StateRoot   [32]byte
	TxRoot      [32]byte
	ReceiptRoot [32]byte

	Difficulty    *big.Int
	GasLimit      uint64
	GasUsed       uint64
	Time          int64
	ArbitraryData []byte
	MixDigest     [32]byte
	Nonce         BlockNonce
}

type Blockposition struct {
	Line int
	List int
}
type SerializedBlock struct {
	Header       *Header
	Transactions []*tx.SerializeTx
	Hash         [32]byte
	Receivedat   uint64
	Receivedfrom [32]byte
	Ethash       [32]byte
	Tree         *Hasher.MerkleTree
	Position     Blockposition
	Parenthash   [32]byte
}
type Blocks []*Block

func (Blocklist Blocks) Len() int {

	return len(Blocklist)

}

func (Blocklist Blocks) Less(i, j int) bool {
	if Blocklist[i].Header.Blocknumber < Blocklist[j].Header.Blocknumber {
		return true
	}
	return false
}

func (Blocklist Blocks) Swap(i, j int) {
	b := Blocklist[i]
	Blocklist[i] = Blocklist[j]
	Blocklist[j] = b

}

func (h *Header) EncodeNonce(nonce uint64) BlockNonce {
	var n BlockNonce
	binary.BigEndian.PutUint64(n[:], nonce)
	return n
}

func (h *Header) ConvertNoncetoUint64(nonce BlockNonce) uint64 {
	return binary.BigEndian.Uint64(nonce[:])
}

func (h *Header) CheckField() error {
	if h.Difficulty != nil {
		if difflen := h.Difficulty.BitLen(); difflen > 80 {
			return fmt.Errorf("too large difficulty,bitlen:%d,exceeds 80", difflen)
		}
	}
	//if datalen exceeds
	if datalen := len(h.ArbitraryData); datalen > cfg.MaxHeaderData {

		return fmt.Errorf("too large data,%d Byte", datalen)
	}
	return nil

}

func (h *Header) Serialize() ([]byte, error) {
	var encode bytes.Buffer
	enc := gob.NewEncoder(&encode)
	err := enc.Encode(h)

	return encode.Bytes(), err

}
func NewBlock(h *Header, txs tx.Transactions, receipts receipt.Receipts, hasher Hasher.Trie) *Block {
	//var hasher types.TrieHasher
	b := &Block{Header: h}
	if len(txs) == 0 {
		b.Header.TxRoot = [32]byte{}
	} else {
		//b.Header.TxRoot = types.DeriveSha(txs, trie.NewStackTrie(nil))
		txcontents := make([]Hasher.Content, len(txs))
		for index, value := range txs {
			txcontents[index] = *value

		}
		b.Header.TxRoot = hasher.DeriveHash(&txcontents)

		b.Transactions = make(tx.Transactions, len(txs))
		copy(b.Transactions, txs)
	}
	b.Header.Time = time.Now().Unix()
	hash, _ := b.HashBlock()
	b.Hash = hash
	if len(receipts) == 0 {
		b.Header.ReceiptRoot = [32]byte{}
	} else {

		for index, _ := range receipts {
			receipts[index].BlockHash = hash
		}
		receiptcontents := make([]Hasher.Content, len(receipts))
		for index, value := range receipts {
			receiptcontents[index] = *value

		}
		b.Header.ReceiptRoot = hasher.DeriveHash(&receiptcontents)
	}

	return b

}

//type Block struct {
//	Header *Header
//
//	Transactions []*tx.Transaction
//
//	// caches
//	hash [32]byte
//
//	ReceivedAt   int64
//	ReceivedFrom [32]byte
//}
//type BlockNonce [8]byte
//
//type Header struct {
//	ParentHash [32]byte
//
//	Coinbase    [20]byte //miner address
//	StateRoot   [32]byte
//	TxRoot      [32]byte
//	ReceiptRoot [32]byte
//
//	Difficulty    *big.Int
//	Number        *big.Int
//	GasLimit      uint64
//	GasUsed       uint64
//	Time          uint64
//	ArbitraryData []byte
//	MixDigest     [32]byte
//	Nonce         BlockNonce
//}
//func MakeBenchBlock(sb tx.Smallbankop, config cfg.Config) *Block {
//	curve := elliptic.P256()
//
//	max := new(big.Int)
//	max.Exp(big.NewInt(2), big.NewInt(130), nil).Sub(max, big.NewInt(1))
//	n, err := rand.Int(rand.Reader, max)
//	if err != nil {
//		//error handling
//	}
//	var (
//		privkey, _ = ecdsa.GenerateKey(curve, rand.Reader)
//		txs        = make([]*tx.Transaction, 100)
//	)
//	h := &Header{
//		Difficulty:    n,
//		GasLimit:      12345678,
//		Time:          time.Now().Unix(),
//		ArbitraryData: []byte("Benchblock"),
//	}
//
//	txs = sb.GenerateBatchSmallbankTx(config, privkey)
//
//	return &Block{Header: h, Transactions: txs}
//}

func (B Block) ConverttoByte() []byte {
	var diffbytes []byte
	if B.Header.Difficulty != nil {
		diffbytes = B.Header.Difficulty.Bytes()
	}

	s := [][]byte{
		B.Header.ArbitraryData[:], tx.Uint64toByte(B.Header.GasLimit), tx.Uint64toByte(B.Header.GasUsed),
		B.Header.TxRoot[:], B.Header.ReceiptRoot[:], B.Header.StateRoot[:], diffbytes,
		B.Header.Blockroot[:], []byte(strconv.FormatInt(B.Header.Time, 10)), B.Header.Coinbase[:],
		B.Header.MixDigest[:], B.Header.Nonce[:],
	}
	return bytes.Join(s, []byte{})
}
func (B *Block) Serialize() ([]byte, error) {
	var encode bytes.Buffer
	serializedTxs := tx.ConvertToSerializeTx(B.Transactions)
	//transactions := make([]tx.Transaction, len(B.Transactions))

	sb := SerializedBlock{
		Header:       B.Header,
		Transactions: serializedTxs,
		Hash:         B.Hash,
		Receivedat:   B.ReceivedAt,
		Receivedfrom: B.ReceivedFrom,
		Position:     *B.Position,
		Ethash:       B.Ethash,
		//Tree:         B.Tree,
		Parenthash: B.Parenthash,
	}

	enc := gob.NewEncoder(&encode)
	err := enc.Encode(sb)

	if err != nil {
		log.Panic("block encode fail:", err)
	}

	return encode.Bytes(), nil
}

func (B Block) CalculateHash() ([]byte, error) {
	h := sha256.New()
	BlockByteCopy := B.ConverttoByte()
	if _, err := h.Write(BlockByteCopy); err != nil {
		return nil, err
	}

	return h.Sum(nil), nil
}
func (B Block) Equals(other Hasher.Content) (bool, error) {
	return bytes.Equal(B.ConverttoByte(), other.(Block).ConverttoByte()), nil
}

func (B *Block) HashBlock() ([32]byte, error) {
	var diffbytes []byte
	if B.Header.Difficulty != nil {
		diffbytes = B.Header.Difficulty.Bytes()
	}
	s := [][]byte{
		B.Header.ArbitraryData[:],
		tx.Uint64toByte(B.Header.GasLimit),
		tx.Uint64toByte(B.Header.GasUsed),
		B.Header.TxRoot[:],
		B.Header.StateRoot[:],
		diffbytes,
		[]byte(strconv.FormatInt(B.Header.Time, 10)),
		tx.Uint64toByte(B.Header.Blocknumber),
		B.Header.Blockroot[:],
	}
	//Difficulty:    bp.bc.Difficulty,
	//GasLimit:      42000 * config.Blocksize,
	//	ArbitraryData: []byte("Blockpacked"),
	//	Blocknumber:   receipts[0].BlockNumber,
	//	GasUsed:       gasused,
	//	StateRoot:     stateroot,
	//	txroot
	//time
	bytes := bytes.Join(s, []byte{})
	bytecopy := [32]byte{}
	copy(bytecopy[:], bytes)
	return bytecopy, nil
}

func (B *Block) CalculateEthash() ([]byte, error) {
	h := sha256.New()
	BlockByteCopy := B.ConverttoByte()

	if _, err := h.Write(BlockByteCopy); err != nil {
		return nil, err
	}

	return h.Sum(nil), nil
}
func NewGenesisBlock() *Block {
	header := &Header{
		Nonce: [8]byte{6, 6, 6, 6, 6, 6, 6, 6},
	}
	block := NewBlock(header, []*tx.Transaction{}, []*receipt.Receipt{}, Hasher.NewNormalHasher())
	//block.Header.Blockroot = [32]byte{}
	block.Header.Coinbase = [20]byte{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}
	block.Header.ArbitraryData = []byte("This is genesis block of a single chain")
	return block
}

//func (b Blocks) Len() int {
//	return len(b)
//}
//
//type h struct {
//	number uint64
//	nonce  [8]byte
//}
//
//func (b Blocks) EncodeIndex(i int, w *bytes.Buffer) {
//	selected := b[i]
//	data := &h{
//		selected.Header.Blocknumber,
//		selected.Header.Nonce,
//	}
//
//	rlp.Encode(w, data)
//
//}

//func makeBenchBlock() *Block {
//	var (
//		key, _   = crypto.GenerateKey()
//		txs      = make([]*Transaction, 70)
//		receipts = make([]*Receipt, len(txs))
//		signer   = LatestSigner(params.TestChainConfig)
//		uncles   = make([]*Header, 3)
//	)
//	Header := &Header{
//		Difficulty: math.BigPow(11, 11),
//		Number:     math.BigPow(2, 9),
//		GasLimit:   12345678,
//		GasUsed:    1476322,
//		Time:       9876543,
//		Extra:      []byte("coolest block on chain"),
//	}
//	for i := range txs {
//		amount := math.BigPow(2, int64(i))
//		price := big.NewInt(300000)
//		data := make([]byte, 100)
//		tx := NewTransaction(uint64(i), common.Address{}, amount, 123457, price, data)
//		signedTx, err := SignTx(tx, signer, key)
//		if err != nil {
//			panic(err)
//		}

//		txs[i] = signedTx
//		receipts[i] = NewReceipt(make([]byte, 32), false, tx.Gas())
//	}
//	for i := range uncles {
//		uncles[i] = &Header{
//			Difficulty: math.BigPow(11, 11),
//			Number:     math.BigPow(2, 9),
//			GasLimit:   12345678,
//			GasUsed:    1476322,
//			Time:       9876543,
//			Extra:      []byte("benchmark uncle"),
//		}
//	}
//	return NewBlock(Header, txs, uncles, receipts, newHasher())
//}

//func (ethash *Ethash) mine(block *types.Block, id int, seed uint64, abort chan struct{}, found chan *types.Block) {
//	// Extract some data from the Header
//	var (
//		Header  = block.Header()
//		hash    = ethash.SealHash(Header).Bytes()
//		target  = new(big.Int).Div(two256, Header.Difficulty)
//		number  = Header.Number.Uint64()
//		dataset = ethash.dataset(number, false)
//	)
//	// Start generating random nonces until we abort or find a good one
//	var (
//		attempts  = int64(0)
//		nonce     = seed
//		powBuffer = new(big.Int)
//	)
//	logger := ethash.config.Log.New("miner", id)
//	logger.Trace("Started ethash search for new nonces", "seed", seed)
//search:
//	for {
//		select {
//		case <-abort:
//			// Mining terminated, update stats and abort
//			logger.Trace("Ethash nonce search aborted", "attempts", nonce-seed)
//			ethash.hashrate.Mark(attempts)
//			break search
//
//		default:
//			// We don't have to update hash rate on every nonce, so update after after 2^X nonces
//			attempts++
//			if (attempts % (1 << 15)) == 0 {
//				ethash.hashrate.Mark(attempts)
//				attempts = 0
//			}
//			// Compute the PoW value of this nonce
//			digest, result := hashimotoFull(dataset.dataset, hash, nonce)
//			if powBuffer.SetBytes(result).Cmp(target) <= 0 {
//				// Correct nonce found, create a new Header with it
//				Header = types.CopyHeader(Header)
//				Header.Nonce = types.EncodeNonce(nonce)
//				Header.MixDigest = common.BytesToHash(digest)
//
//				// Seal and return a block (if still needed)
//				select {
//				case found <- block.WithSeal(Header):
//					logger.Trace("Ethash nonce found and reported", "attempts", nonce-seed, "nonce", nonce)
//				case <-abort:
//					logger.Trace("Ethash nonce found but discarded", "attempts", nonce-seed, "nonce", nonce)
//				}
//				break search
//			}
//			nonce++
//		}
//	}
//	// Datasets are unmapped in a finalizer. Ensure that the dataset stays live
//	// during sealing so it's not unmapped while being read.
//	runtime.KeepAlive(dataset)
//}
