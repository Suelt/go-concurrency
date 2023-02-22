package tx

import (
	"UndergraduateProj/Hasher"
	"bytes"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"log"
	"time"
)

type Transaction struct {
	sender        [32]byte
	receiver      [32]byte //four receivers,each denoting a different smallbank transaction
	signature     [64]byte //signature of tx
	hash          [32]byte
	arbitraryData []byte //Data
	//sets          Rwsets      //Read-write sets  -  kv
	gaslimit    uint64
	gasprice    uint64
	timestamp   uint64
	nonce       uint64
	channellist []chan bool // list of channels to listen and send
	//TxInterface             //tx_interface.go
}

type SerializeTx struct {
	Sender        [32]byte
	Receiver      [32]byte
	ArbitraryData []byte
	//sets          Rwsets      //Read-write sets  -  kv
	Gaslimit  uint64
	Gasprice  uint64
	Timestamp uint64 // Timestamp of tx
	Nonce     uint64
	Signature [64]byte //signature of tx
}

//NormalTransaction Key accessed
type Key struct {
	K string
}

type Transfer struct {
	From string
	To   string
}

type ConcurrentSet []*Rwsets

type Rwsets struct {
	Rkv []Set
	Wkv []Set
}

type Set struct {
	K string
	V float64
}

type Transactions []*Transaction

func (Tx *Transaction) Serialize() ([]byte, error) {
	//Todo
	//Encode and return byte array
	var encode bytes.Buffer

	// construct the pure version of tx
	st := SerializeTx{
		Sender:        Tx.sender,
		Receiver:      Tx.receiver,
		ArbitraryData: Tx.arbitraryData,
		Nonce:         Tx.nonce,
		Timestamp:     Tx.timestamp,
		Signature:     Tx.signature,
		Gaslimit:      Tx.gaslimit,
		Gasprice:      Tx.gasprice,
	}

	enc := gob.NewEncoder(&encode)
	err := enc.Encode(st)

	if err != nil {
		log.Panic("tx encode fail:", err)
	}

	return encode.Bytes(), nil

}

func (Tx *Transaction) SerializeLength() uint64 {
	var encode bytes.Buffer

	// construct the pure version of tx
	st := SerializeTx{
		Sender:        Tx.sender,
		Receiver:      Tx.receiver,
		ArbitraryData: Tx.arbitraryData,
		Nonce:         Tx.nonce,
		Timestamp:     Tx.timestamp,
		Signature:     Tx.signature,
		Gaslimit:      Tx.gaslimit,
		Gasprice:      Tx.gasprice,
	}

	enc := gob.NewEncoder(&encode)
	err := enc.Encode(st)
	if err != nil {
		log.Panic("tx encode fail:", err)
	}

	return uint64(len(encode.Bytes()))
}
func (Tx *Transaction) Verify() bool {

	return true
}
func (Tx *Transaction) Sign(privKey *ecdsa.PrivateKey) {

	txCopy := *Tx
	txCopy.hash = [32]byte{}
	txCopy.signature = [64]byte{}
	dataToSign := fmt.Sprintf("%x\n", txCopy)

	r, s, err := ecdsa.Sign(rand.Reader, privKey, []byte(dataToSign))
	if err != nil {
		log.Panic(err)
	}

	signature := append(r.Bytes(), s.Bytes()...)

	var sign64 [64]byte
	copy(sign64[:], signature)

	Tx.signature = sign64
}
func (Tx Transaction) DetectDependencies() []chan bool {
	var list []chan bool
	//Todo Detect
	return list
}

func (Tx *Transaction) AddIntoChanList(ch chan bool) {
	Tx.SetChannellist(append(Tx.Channellist(), ch))

}

func (Tx *Transaction) Execute(sb Smallbankop) (*Rwsets, error) {

	rw, err := sb.SmallbankExecution(Tx)
	if err != nil {
		log.Panicln("Execution failed:", err)
	}
	return rw, nil
}

func Commit(str []string) error {
	for index, single := range str {
		fmt.Printf("transaction%d\n%s", index, single)
	}
	return nil
}

func InitTx(sender string, receiver string) *Transaction {

	//type Transaction struct {
	//	sender        [32]byte //transaction proposer
	//	receiver      [32]byte //four receivers,each denoting a different smallbank transaction
	//	signature     [32]byte //signature of tx
	//	hash          [32]byte
	//	arbitraryData [256]byte //Data
	//	//sets          Rwsets      //Read-write sets  -  kv

	//	timestamp   int64       // Timestamp of tx
	//	nonce       int64       // nonce of mining
	//	channellist []chan bool // list of channels to listen and send
	//	TxInterface             //tx_interface.go
	//}
	var Tx Transaction
	senderbyte := StringToBytearray(sender)
	receiverbyte := StringToBytearray(receiver)
	//arbitrarydatabyte := StringToBytearray(arbitrarydata)
	Tx.gaslimit = 42000
	Tx.gasprice = Randint(50, 70)
	Tx.SetSender(senderbyte)
	Tx.SetReceiver(receiverbyte)
	//Tx.SetArbitraryData(arbitrarydatabyte)
	Tx.SetTimestamp(uint64(time.Now().Unix()))
	return &Tx
}

func (Tx *Transaction) Channellist() []chan bool {
	return Tx.channellist
}

func (Tx *Transaction) SetChannellist(channellist []chan bool) {
	Tx.channellist = channellist
}

//func (Tx *Transaction) Sets() Rwsets {
//	return Tx.sets
//}
//
//func (Tx *Transaction) SetSets(sets Rwsets) {
//	Tx.sets = sets
//}

func (Tx *Transaction) Sender() [32]byte {
	return Tx.sender
}

func (Tx *Transaction) SetSender(sender [32]byte) {
	Tx.sender = sender
}

func (Tx *Transaction) Receiver() [32]byte {
	return Tx.receiver
}

func (Tx *Transaction) SetReceiver(receiver [32]byte) {
	Tx.receiver = receiver
}

func (Tx *Transaction) Signature() [64]byte {
	return Tx.signature
}

func (Tx *Transaction) SetSignature(signature [64]byte) {
	Tx.signature = signature
}

func (Tx *Transaction) ArbitraryData() []byte {
	return Tx.arbitraryData
}

func (Tx *Transaction) SetArbitraryData(arbitraryData []byte) {
	Tx.arbitraryData = arbitraryData
}

func (Tx *Transaction) Timestamp() uint64 {
	return Tx.timestamp
}

func (Tx *Transaction) SetTimestamp(timestamp uint64) {
	Tx.timestamp = timestamp
}

func (Tx *Transaction) Nonce() uint64 {
	return Tx.nonce
}

func (Tx *Transaction) SetNonce(nonce uint64) {
	Tx.nonce = nonce
}

func (Tx *Transaction) Hash() [32]byte {
	return Tx.hash
}

func (Tx *Transaction) SetHash(hash [32]byte) {
	Tx.hash = hash
}

func (Tx *Transaction) Gaslimit() uint64 {
	return Tx.gaslimit
}

func (Tx *Transaction) SetGaslimit(gaslimit uint64) {
	Tx.gaslimit = gaslimit
}

func (Tx *Transaction) Gasprice() uint64 {
	return Tx.gasprice
}

func (Tx *Transaction) SetGasprice(gasprice uint64) {
	Tx.gasprice = gasprice
}

//func NewTransaction(TI TxInterface) *Transaction {
//	return &Transaction{TI: TI}
//}

func (Tx Transaction) ConverttoByte() []byte {
	s := [][]byte{Tx.sender[:], Tx.receiver[:], Tx.signature[:], Tx.hash[:], Tx.arbitraryData[:],
		Uint64toByte(Tx.timestamp),
		Uint64toByte(Tx.nonce),
		Uint64toByte(Tx.gasprice),
		Uint64toByte(Tx.gaslimit)}
	return bytes.Join(s, []byte{})
}

func (Tx Transaction) CalculateHash() ([]byte, error) {
	h := sha256.New()
	TxByteCopy := Tx.ConverttoByte()
	if _, err := h.Write(TxByteCopy); err != nil {
		return nil, err
	}

	return h.Sum(nil), nil
}

func (Tx Transaction) Equals(other Hasher.Content) (bool, error) {
	return bytes.Equal(Tx.ConverttoByte(), other.(Transaction).ConverttoByte()), nil
}

func ConvertToSerializeTx(txs Transactions) []*SerializeTx {
	serializedList := make([]*SerializeTx, len(txs))
	for i := 0; i < len(txs); i++ {
		serializedList[i] = &SerializeTx{
			Sender:        txs[i].sender,
			Receiver:      txs[i].receiver,
			Gaslimit:      txs[i].gaslimit,
			Gasprice:      txs[i].gasprice,
			Signature:     txs[i].signature,
			Nonce:         txs[i].nonce,
			Timestamp:     txs[i].timestamp,
			ArbitraryData: txs[i].arbitraryData,
		}
	}
	return serializedList
}

func ConvertToNormalTx(txs []*SerializeTx) Transactions {
	normalTxList := make([]*Transaction, len(txs))
	for i := 0; i < len(txs); i++ {
		normalTxList[i] = &Transaction{
			sender:        txs[i].Sender,
			receiver:      txs[i].Receiver,
			gaslimit:      txs[i].Gaslimit,
			gasprice:      txs[i].Gasprice,
			signature:     txs[i].Signature,
			nonce:         txs[i].Nonce,
			timestamp:     txs[i].Timestamp,
			arbitraryData: txs[i].ArbitraryData,
		}
	}
	return normalTxList
}
func NewTxWithSerializedTx(serialized *SerializeTx) *Transaction {
	t := &Transaction{
		serialized.Sender,
		serialized.Receiver,
		serialized.Signature,
		[32]byte{},
		serialized.ArbitraryData,
		serialized.Gaslimit,
		serialized.Gasprice,
		serialized.Timestamp,
		serialized.Nonce,
		nil,
	}
	hash, _ := t.CalculateHash()
	copy(t.hash[:], hash)
	return t
}

//func (s Transactions) Len() int {
//	return len(s)
//}
//

//func (s Transactions) EncodeIndex(i int, w *bytes.Buffer) {
//	tx := s[i]
//	Data := &SerializeTx{
//		Sender:        tx.sender,
//		Receiver:      tx.receiver,
//		Gaslimit:      tx.gaslimit,
//		Gasprice:      tx.gasprice,
//		Signature:     tx.signature,
//		Nonce:         tx.nonce,
//		Timestamp:     tx.timestamp,
//		ArbitraryData: tx.arbitraryData,
//	}
//
//	err := rlp.Encode(w, Data)
//	if err != nil {
//		return
//	}
//
//}

//sender        [32]byte //transaction proposer
//receiver      [32]byte //four receivers,each denoting a different smallbank transaction
//signature     [64]byte //signature of tx
//hash          [32]byte
//arbitraryData [32]byte //Data
////sets          Rwsets      //Read-write sets  -  kv
//
//timestamp   uint64 // Timestamp of tx
//nonce       uint64
