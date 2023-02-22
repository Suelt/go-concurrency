package blockchain

import (
	"UndergraduateProj/block"
	"UndergraduateProj/tx"
	"bytes"
	"crypto/rand"
	"encoding/gob"
	"github.com/syndtr/goleveldb/leveldb"
	"log"
	"math/big"
	"sort"
	"strconv"
)

func LoadDB() ([]*leveldb.DB, error) {
	DBs := make([]*leveldb.DB, 3)
	DBNormal, err := leveldb.OpenFile("DBNormal", nil)
	if err != nil {
		return nil, err
	}
	DBChain, err := leveldb.OpenFile("DBChain", nil)
	if err != nil {
		return nil, err
	}

	DBState, err := leveldb.OpenFile("DBState", nil)
	if err != nil {
		return nil, err
	}
	DBs[0] = DBNormal
	DBs[1] = DBChain
	DBs[2] = DBState
	return DBs, nil
}

func StoreTransactionByHash(db *leveldb.DB, transaction *tx.Transaction) error {
	data, err := transaction.Serialize()
	if err != nil {
		return err
	}
	hash, _ := transaction.CalculateHash()
	//s := [][]byte{[]byte("type:tx"), hash[:]}
	err = db.Put(hash[:], data, nil)
	if err != nil {
		return err
	}

	return nil
}

func StoreTransactionByPosition(db *leveldb.DB, transaction *tx.Transaction, pos *block.Blockposition, index uint64) error {
	posbytes := ConvertBlockpositionandIndextoBytes(pos, index)
	data, err := transaction.Serialize()
	if err != nil {
		return err
	}
	err = db.Put(posbytes, data, nil)

	return err
}
func ConvertBlockpositionandIndextoBytes(pos *block.Blockposition, index uint64) []byte {
	var (
		line = pos.Line
		list = pos.List
	)
	linestr := strconv.Itoa(line)
	retstr := strconv.Itoa(list) + linestr
	retstr = retstr + strconv.Itoa(int(index))
	return []byte(retstr)
}

func ConvertBlockpositiontoBytes(pos *block.Blockposition) []byte {
	var (
		line = pos.Line
		list = pos.List
	)
	linestr := strconv.Itoa(line)
	retstr := strconv.Itoa(list) + linestr
	return []byte(retstr)
}
func FetchTransactionByHash(db *leveldb.DB, txhash []byte) (*tx.Transaction, error) {
	data, err := db.Get(txhash, nil)
	if err != nil {
		return nil, err
	}
	return DeserializeTx(data)

}

func FetchTransactionByPosition(db *leveldb.DB, pos *block.Blockposition, index uint64) (*tx.Transaction, error) {
	data, err := db.Get(ConvertBlockpositionandIndextoBytes(pos, index), nil)
	if err != nil {
		return nil, err
	}
	return DeserializeTx(data)
}

func StoreBlockHeader(db *leveldb.DB, b *block.Block) error {
	data, err := b.Header.Serialize()
	if err != nil {
		return err
	}
	hash, errhashing := b.CalculateHash()
	if errhashing != nil {
		return errhashing
	}
	errdbstore := db.Put(hash[:], data, nil)
	return errdbstore
}

func FetchBlockHeaderByhash(db *leveldb.DB, headerhash []byte) (*block.Header, error) {
	data, err := db.Get(headerhash, nil)
	if err != nil {
		return nil, err
	}
	return DeserializeHeader(data)

}

func StoreBlock(db *leveldb.DB, b *block.Block) error {
	data, err := b.Serialize()
	if err != nil {
		return nil
	}
	posbytes := ConvertBlockpositiontoBytes(b.Position)
	err = db.Put(posbytes, data, nil)
	return err
}

func FetchBlock(db *leveldb.DB, position *block.Blockposition) (*block.Block, error) {
	posbytes := ConvertBlockpositiontoBytes(position)
	data, err := db.Get(posbytes, nil)
	if err != nil {
		return nil, err
	}
	b, errdeserialize := DeserializeBlock(data)
	return b, errdeserialize
}

func RegenerateBlockchainViaDB(db *leveldb.DB) (*Blockchain, error) {

	max := new(big.Int)
	max.Exp(big.NewInt(2), big.NewInt(70), nil).Sub(max, big.NewInt(1))
	n, _ := rand.Int(rand.Reader, max)
	bc := &Blockchain{
		Difficulty: n,
	}

	bc.Engine = &Chainworker{Bc: bc}
	bc.BlockProposer = NewBcProposer(bc)
	bc.StateProcessor = NewConcurrentProcessor(bc)
	bc.ChainValidator = NewBcValidator(bc)
	position := &block.Blockposition{
		List: 0,
	}
	for i := 0; i < 32; i++ {
		bc.Chains[i] = NewSinglechain(uint64(i))
		bc.Chains[i].Chain[0].Header.Difficulty = bc.Difficulty
		position.Line = i
		bc.Chains[i].Chain[0].Position = position
		hash, _ := bc.Chains[i].Chain[0].HashBlock()
		bc.Chains[i].Chain[0].Hash = hash

	}
	iter := db.NewIterator(nil, nil)
	var blocklist block.Blocks
	for iter.Next() {

		//key := iter.Key()
		value := iter.Value()
		b, err := DeserializeBlock(value)
		if err != nil {
			return nil, err
		}
		blocklist = append(blocklist, b)
		//pos := b.Position
		//log.Println(pos.Line)
		//bc.Chains[pos.Line].Chain = append(bc.Chains[pos.Line].Chain, b)
		//bc.Chains[pos.Line].Currentblock = b

	}
	sort.Sort(blocklist)
	for i := 0; i < len(blocklist); i++ {
		pos := blocklist[i].Position
		log.Println(pos.Line)
		bc.Chains[pos.Line].Chain = append(bc.Chains[pos.Line].Chain, blocklist[i])
		bc.Chains[pos.Line].Currentblock = blocklist[i]
	}
	iter.Release()
	err := iter.Error()
	if err != nil {
		log.Panicln(err)
	}
	bc.GenerateBlockRoot()
	bc.Tree = bc.GenerateBlockMerkletree()

	return bc, nil
}

func RegenerateStateViaDB(db *leveldb.DB) (tx.Smallbankop, error) {
	var (
		m  map[string]int
		sb tx.Smallbank
	)
	dataset := make([]*tx.Smallbankdata, tx.Smallbanksize)
	m = make(map[string]int, tx.Smallbanksize)
	iter := db.NewIterator(nil, nil)
	for iter.Next() {
		value := iter.Value()
		d, err := DeserializeSmallbankData(value)
		if err != nil {
			return nil, err
		}
		m[d.S] = d.Index
		dataset[d.Index] = &d.Data
	}
	sb.Dataset = dataset
	sb.Mp = m
	return &sb, nil
}

func StoreState(db *leveldb.DB, s string, index int, data *tx.Smallbankdata) error {
	d := &tx.SmallbankdataForDB{
		s,
		index,
		*data,
	}
	serialized, err := d.Serialize()
	if err != nil {
		return err
	}
	err = db.Put([]byte(s), serialized, nil)
	return err
}
func FetchState(db *leveldb.DB, s string) (*tx.SmallbankdataForDB, error) {
	data, err := db.Get([]byte(s), nil)
	if err != nil {
		return nil, err
	}
	return DeserializeSmallbankData(data)

}
func StoreWholeState(db *leveldb.DB, sb *tx.Smallbank) error {
	for index, value := range sb.Mp {
		err := StoreState(db, index, value, sb.Dataset[value])
		if err != nil {
			return err
		}
	}
	return nil
}
func DeserializeTx(data []byte) (*tx.Transaction, error) {
	var transaction tx.SerializeTx
	reader := bytes.NewReader(data)
	decode := gob.NewDecoder(reader)
	err := decode.Decode(&transaction)
	tx := tx.NewTxWithSerializedTx(&transaction)
	return tx, err

}

func DeserializeHeader(data []byte) (*block.Header, error) {
	var header block.Header
	reader := bytes.NewReader(data)
	decode := gob.NewDecoder(reader)
	err := decode.Decode(&header)

	return &header, err
}

func DeserializeBlock(data []byte) (*block.Block, error) {
	var b block.SerializedBlock
	reader := bytes.NewReader(data)
	decode := gob.NewDecoder(reader)
	err := decode.Decode(&b)
	if err != nil {
		return nil, err
	}
	txs := tx.ConvertToNormalTx(b.Transactions)
	return &block.Block{
		Header:       b.Header,
		Transactions: txs,
		Hash:         b.Hash,
		Ethash:       b.Ethash,
		Position:     &b.Position,
		Tree:         b.Tree,
		Parenthash:   b.Parenthash,
		ReceivedFrom: b.Receivedfrom,
		ReceivedAt:   b.Receivedat,
	}, nil

}

func DeserializeSmallbankData(data []byte) (*tx.SmallbankdataForDB, error) {
	var sb tx.SmallbankdataForDB
	reader := bytes.NewReader(data)
	decode := gob.NewDecoder(reader)
	err := decode.Decode(&sb)
	return &sb, err
}
