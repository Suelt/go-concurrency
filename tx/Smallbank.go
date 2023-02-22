package tx

import (
	"UndergraduateProj/Hasher"
	cfg "UndergraduateProj/config"
	"bytes"
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/binary"
	"encoding/gob"
	"fmt"
	"github.com/deckarep/golang-set"
	"github.com/pingcap/go-ycsb/pkg/generator"
	"log"
	"math"
	"math/rand"
	"sync"
	"time"
)

const Smallbanksize = 8192

const Txsperblock = 128

const Pw = 0.5

type Smallbank struct {
	Data    sync.Map
	Mp      map[string]int
	Dataset Dataset
	SCL     []SmallbankContractLog
}
type Smallbankdata struct {
	Address string

	Data float64
}

type SmallbankdataForDB struct {
	S     string
	Index int
	Data  Smallbankdata
}

func (d *SmallbankdataForDB) Serialize() ([]byte, error) {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(d)
	if err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

func float64ToByte(f float64) []byte {
	var buf [8]byte
	binary.BigEndian.PutUint64(buf[:], math.Float64bits(f))
	return buf[:]
}

type Dataset []*Smallbankdata

//
//func (d Dataset) Len() int {
//	return len(d)
//}
//func (d Dataset) EncodeIndex(i int, w *bytes.Buffer) {
//	smalldata := d[i]
//	Data := &logdata{
//		Address: smalldata.Address,
//	}
//	err := rlp.Encode(w, Data)
//	if err != nil {
//		return
//	}
//}

func (sb *Smallbank) GenerateStateRoot() [32]byte {
	list := sb.Dataset

	//fmt.Printf("%x\n", list[501].Address)
	//fmt.Printf("%x\n", list[500].Address)
	//return types.DeriveSha(list, trie.NewStackTrie(nil))
	Smallbankcontents := make([]Hasher.Content, Smallbanksize)
	//fmt.Println(list)
	for index, value := range list {
		Smallbankcontents[index] = *value
	}
	//fmt.Printf("%x\n", Hasher.NewNormalHasher().DeriveHash(&Smallbankcontent))
	return Hasher.NewNormalHasher().DeriveHash(&Smallbankcontents)
}

func (sb *Smallbank) GenerateStatetrie() *Hasher.MerkleTree {
	blockcontents := make([]Hasher.Content, Smallbanksize)
	list := sb.Dataset
	for index, value := range list {
		blockcontents[index] = *value

	}
	t, err := Hasher.NewTree(blockcontents)
	cfg.UseRw(err)
	return t
}
func (sb *Smallbank) WriteState(rw *Rwsets) {

	if rw.Wkv == nil {
		return
	}
	writeset := rw.Wkv
	for _, value := range writeset {
		index := sb.Mp[string(value.K[:])]
		sb.Dataset[index].Data = value.V
	}
}
func (sb *Smallbank) CommittoState(sets []*Rwsets) {
	for _, set := range sets {
		sb.WriteState(set)
	}
}

//return list

func (sbd Smallbankdata) Equals(other Hasher.Content) (bool, error) {
	return bytes.Equal(sbd.ConverttoByte(), other.(Smallbankdata).ConverttoByte()), nil
}
func (sbd Smallbankdata) ConverttoByte() []byte {
	s := [][]byte{[]byte(sbd.Address), float64ToByte(sbd.Data)}
	return bytes.Join(s, []byte{})
}

func (sbd Smallbankdata) CalculateHash() ([]byte, error) {
	h := sha256.New()
	datacopy := sbd.ConverttoByte()
	if _, err := h.Write(datacopy); err != nil {
		return nil, err
	}

	return h.Sum(nil), nil
}

type SmallbankContractLog struct {
	receiver [32]byte
	data     []byte
}

func (sb *Smallbank) SmallbankExecution(transaction *Transaction) (*Rwsets, error) {
	receiver := transaction.receiver
	var bytecopy []byte
	for _, value := range receiver {
		if value == 0 {
			break
		}
		bytecopy = append(bytecopy, value)
	}
	StringEquals := string(bytecopy)
	data := transaction.arbitraryData
	//fmt.Println(Data)
	var decoder bytes.Buffer
	dec := gob.NewDecoder(&decoder)
	decoder.Write(data)

	if StringEquals == "Transfer" {
		var transferlist [cfg.SmallbankAccessedkeys]*Transfer
		err := dec.Decode(&transferlist)
		if err != nil {
			log.Fatal("decode error:", err)
		}
		return sb.Transfer(transferlist)
	} else {
		var Depset [cfg.SmallbankAccessedkeys]*Key
		err := dec.Decode(&Depset)
		if err != nil {
			log.Fatal("decode error:", err)
		}

		if StringEquals == "Withdraw" {
			return sb.Withdraw(Depset)
		} else if StringEquals == "Query" {
			return sb.Query(Depset)
		} else {
			return sb.Deposit(Depset)
		}
	}

}

func Timewasting() {
	a := 0
	for j := 0; j < 310; j++ {

		for i := 0; i < math.MaxUint8; i++ {
			a = i
		}
	}
	cfg.UseRw(a)

	//time.Sleep(10 * time.Millisecond)
}
func (sb *Smallbank) Deposit(depset [cfg.SmallbankAccessedkeys]*Key) (*Rwsets, error) {
	//Todo Depset generation
	//depset := sb.NormalCaseConstruction()
	rkv := make([]Set, cfg.SmallbankAccessedkeys)
	wkv := make([]Set, cfg.SmallbankAccessedkeys)

	var v float64
	//var delta float64
	//cfg.UseRw(delta)
	//fmt.Println(1)
	for index, key := range depset {
		var sread, swrite Set
		target := sb.Mp[key.K]

		v = sb.Dataset[target].Data

		sread.K = key.K
		sread.V = v
		rkv[index] = sread
		Timewasting()

		value := sb.DepositComputation(v)
		//fmt.Printf("%f\n", value)

		//delta = v - s.V
		//fmt.Println(s.V, value)
		swrite.V = value
		swrite.K = key.K
		wkv[index] = swrite
		//	sb.Dataset[target].mutex.Lock()
		sb.Dataset[target].Data = value

		//d := &SmallbankdataForDB{
		//	s.K,
		//	target,
		//	*sb.Dataset[target],
		//}
		//serialized, err := d.Serialize()
		//if err != nil {
		//	return nil, err
		//}
		//err = db.Put([]byte(s.K), serialized, nil)
		//	sb.Dataset[target].mutex.Unlock()
		//		sb.Data.Store(string(s.K[:]), v)

		//fmt.Printf("%d:deposit on key %x,token increased by %f\n", index, s.K, delta)
	}

	rw := Rwsets{Rkv: rkv, Wkv: wkv}
	return &rw, nil
}

func (sb *Smallbank) Transfer(TransferList [cfg.SmallbankAccessedkeys]*Transfer) (*Rwsets, error) {

	//type TransferTransaction struct {
	//	Transaction
	//	TransferList []transfer
	//}
	//type transfer struct {
	//	from  []byte
	//	to    []byte
	//	value float64
	//}
	rkv := make([]Set, 2*cfg.SmallbankAccessedkeys)
	wkv := make([]Set, 2*cfg.SmallbankAccessedkeys)
	//var delta float64
	var keyfrom, keyto Key
	var setfrom, setto Set
	var vfrom, vto float64
	var epoch int
	//Todo TransferList generation
	//TransferList := sb.TransferConstruction()
	for index, trans := range TransferList {
		keyfrom.K = trans.From
		keyto.K = trans.To
		targetfrom := sb.Mp[keyfrom.K]

		vfrom = sb.Dataset[targetfrom].Data

		targetto := sb.Mp[keyto.K]

		vto = sb.Dataset[targetto].Data

		setfrom.K = trans.From
		setfrom.V = vfrom
		rkv[index+epoch] = setfrom
		setto.K = trans.To
		setto.V = vto
		rkv[index+epoch+1] = setto
		Timewasting()

		vfrom, vto = sb.TransferComputation(vfrom, vto)
		//Timewasting()
		//delta = vto - set.V
		setfrom.K = trans.From
		setfrom.V = vfrom
		wkv[index+epoch] = setfrom
		setto.K = trans.To
		setto.V = vto
		wkv[index+epoch+1] = setto

		sb.Dataset[targetfrom].Data = vfrom

		sb.Dataset[targetto].Data = vto
		//d := &SmallbankdataForDB{
		//	trans.From,
		//	targetfrom,
		//	*sb.Dataset[targetfrom],
		//}
		//serialized, err1 := d.Serialize()
		//if err1 != nil {
		//	return nil, err1
		//}
		//err1 = db.Put([]byte(trans.From), serialized, nil)
		//if err1 != nil {
		//	return nil, err1
		//}
		//dto := &SmallbankdataForDB{
		//	trans.To,
		//	targetto,
		//	*sb.Dataset[targetto],
		//}
		//serializedto, err2 := dto.Serialize()
		//if err2 != nil {
		//	return nil, err2
		//}
		//err2 = db.Put([]byte(trans.From), serializedto, nil)
		//if err2!=nil{
		//	return nil,err2
		//}
		//sb.Data.Store(string(trans.From[:]), vfrom)
		//sb.Data.Store(string(trans.To[:]), vto)
		epoch++
		//cfg.UseRw(delta)
		//fmt.Printf("%d:key %x transfers %f tokens to key %x\n", index, trans.From, delta, trans.To)
	}

	rw := Rwsets{Rkv: rkv[:], Wkv: wkv[:]}
	return &rw, nil
}

func (sb *Smallbank) Withdraw(depset [cfg.SmallbankAccessedkeys]*Key) (*Rwsets, error) {
	//Todo Depset generation
	//depset := sb.NormalCaseConstruction()

	rkv := make([]Set, cfg.SmallbankAccessedkeys)
	wkv := make([]Set, cfg.SmallbankAccessedkeys)
	s := Set{}
	//var delta float64
	//cfg.UseRw(delta)
	//var v float64
	var v float64
	for index, key := range depset {
		s = Set{}
		target := sb.Mp[key.K]

		v = sb.Dataset[target].Data

		//Dataload, _ := sb.Data.Load(string(key.K[:]))
		//fmt.Println(Dataload)
		//v = InterfaceToFloat64(Dataload)
		s.K = key.K
		s.V = v
		rkv[index] = s
		Timewasting()

		v = sb.WithdrawComputation(v)
		//delta = s.V - v
		s.V = v

		sb.Dataset[target].Data = v
		//d := &SmallbankdataForDB{
		//	s.K,
		//	target,
		//	*sb.Dataset[target],
		//}
		//serialized, err := d.Serialize()
		//if err != nil {
		//	return nil, err
		//}
		//err = db.Put([]byte(s.K), serialized, nil)
		//sb.Data.Store(string(s.K[:]), v)
		wkv[index] = s
		//fmt.Printf("%d:Withdraw on key %x,token decreased by %f\n", index, s.K, delta)

	}

	rw := Rwsets{Rkv: rkv, Wkv: wkv}
	return &rw, nil
}

func (sb *Smallbank) Query(depset [cfg.SmallbankAccessedkeys]*Key) (*Rwsets, error) {
	//Todo Depset generation
	//depset := sb.NormalCaseConstruction()

	rkv := make([]Set, cfg.SmallbankAccessedkeys)
	//wkv := make([]Set, 2)
	s := Set{}
	var v float64
	for index, key := range depset {
		s = Set{}
		//fmt.Printf("%x\n", key.K)
		target := sb.Mp[key.K]
		Timewasting()

		//Dataload, _ := sb.Data.Load(string(key.K[:]))
		v = sb.Dataset[target].Data
		//fmt.Printf("%x\n", string(key.K[:]))

		s.K = key.K
		s.V = v
		rkv[index] = s
		//fmt.Printf("%d:Query on key %x,value accessed:%f\n", index, s.K, v)
	}
	rw := Rwsets{Rkv: rkv, Wkv: nil}
	return &rw, nil
}

func (sb *Smallbank) Init() {

	var data sync.Map
	mp := make(map[string]int, Smallbanksize)
	n := 64
	rand.Seed(time.Now().Unix())
	i := 0
	for i < Smallbanksize {
		randBytes := make([]byte, n/2)
		rand.Read(randBytes)
		mystring := string(randBytes[:])
		//Dataload, _ := sb.Data.Load(mystring)
		//_, ok := InterfaceToFloat64(Dataload)
		//if ok {
		//
		//	fmt.Printf("Key Existed at %x!\n", mystring)
		//	continue
		//
		//} else {
		v := Randfloat(10, 100)
		//	mutex := make([]sync.RWMutex, Smallbanksize)
		d := &Smallbankdata{mystring, v}
		sb.Dataset = append(sb.Dataset, d)
		mp[mystring] = i
		data.Store(mystring, v)
		//fmt.Printf("%x,%f\n", randBytes, v)
		i++
		//fmt.Println(i)

		//}
	}
	sb.Data = data
	sb.Mp = mp
	//fmt.Println("Smallbank Account Generation Completed")

}

//var mutex sync.Mutex

func (sb *Smallbank) NormalCaseRandomization() *Key {
	rand.Seed(time.Now().Unix())
	iterator := rand.Intn(Smallbanksize)
	var keyaccessed Key
	f := func(key, value interface{}) bool {
		iterator--
		if iterator == 0 {
			k := key.(string)

			//fmt.Printf("%x\n", k)
			keyaccessed.K = k
			//fmt.Printf("%x\n", keyaccessed.K)

		}
		return true
	}
	sb.Data.Range(f)
	//fmt.Printf("%x\n", keyaccessed.K)
	return &keyaccessed
}
func (sb *Smallbank) NormalCaseConstruction() [cfg.SmallbankAccessedkeys]*Key {
	//var Mp map[string]string
	//Mp.
	var keys [cfg.SmallbankAccessedkeys]*Key
	for i := 0; i < cfg.SmallbankAccessedkeys; i++ {
		keys[i] = sb.NormalCaseRandomization()
	}
	//fmt.Printf("%x\n", keys[0].K)
	//fmt.Printf("%x\n", keys[1].K)
	return keys

}

func (sb *Smallbank) TransferRandomization() *Transfer {
	var transfer Transfer
	rand.Seed(time.Now().Unix())
	iterator := rand.Intn(Smallbanksize)
	j := iterator
	f1 := func(key, value interface{}) bool {
		iterator--
		if iterator == 0 {
			k := key.(string)

			transfer.From = k
			//keys = append(keys, keyaccessed)
		}
		return true
	}

	f2 := func(key, value interface{}) bool {
		j--
		if j == 0 {
			k := key.(string)

			transfer.To = k
			//keys = append(keys, keyaccessed)
		}
		return true
	}
	sb.Data.Range(f1)
	sb.Data.Range(f2)
	return &transfer

}

func (sb *Smallbank) TransferConstruction() [cfg.SmallbankAccessedkeys]*Transfer {
	var transferlist [cfg.SmallbankAccessedkeys]*Transfer

	for i := 0; i < cfg.SmallbankAccessedkeys; i++ {

		transferlist[i] = sb.TransferRandomization()
	}

	return transferlist
}

func (sb *Smallbank) WithdrawComputation(v float64) float64 {

	return 0.9 * v
}

func (sb *Smallbank) DepositComputation(v float64) float64 {
	//target := sb.Mp[str]

	//	sb.Dataset[target].mutex.RUnlock()
	return 1.1 * v
}
func (sb *Smallbank) TransferComputation(vfrom float64, vto float64) (float64, float64) {

	return 0.5 * vfrom, 0.5*vfrom + vto
}

func (sb *Smallbank) GenerateSmallbankTx(key *ecdsa.PrivateKey, set mapset.Set, zipfian *generator.Zipfian, r *rand.Rand) *Transaction {
	executiontype := [4]string{"Deposit", "Withdraw", "Transfer", "Query"}

	//x := 1
	//	x := 0
	var x int
	randomization := [2]int{Pw * 10000, 10000 - Pw*10000}
	m := r.Intn(10000)
	if m < randomization[0] {
		x = r.Intn(3)

	} else {
		x = 3

	}

	tx := InitTx("sender", executiontype[x])

	var encoder bytes.Buffer
	enc := gob.NewEncoder(&encoder)
	//dec := gob.NewDecoder(&encoder)
	if executiontype[x] == "Transfer" {
		var transferlist [cfg.SmallbankAccessedkeys]*Transfer
		for i := 0; i < cfg.SmallbankAccessedkeys; i++ {

			a := uint64(zipfian.Next(r))
			var b uint64
			for {
				b = uint64(zipfian.Next(r))
				if a != b {
					break
				}
			}

			transferlist[i] = &Transfer{
				sb.Dataset[a].Address,
				sb.Dataset[b].Address,
			}
		}
		//InputTypeA := sb.TransferConstruction()
		//if set.Contains(InputTypeA[0].From) {
		//	log.Panicln("K existed")
		//}
		//if set.Contains(InputTypeA[0].To) {
		//	log.Panicln("K existed")
		//}
		//set.Add(InputTypeA[0].From)
		//set.Add(InputTypeA[0].To)
		err := enc.Encode(transferlist)
		if err != nil {
			log.Fatal("encode error:", err)
		}
		//fmt.Println(encoder.Bytes())
		tx.arbitraryData = encoder.Bytes()
		//copy(tx.arbitraryData, encoder.Bytes())

	} else {
		//fmt.Println(1)
		var keyset [cfg.SmallbankAccessedkeys]*Key
		for i := 0; i < cfg.SmallbankAccessedkeys; i++ {
			a := uint64(zipfian.Next(r))
			//a := r.Intn(Smallbanksize)
			//rand.Seed(time.Now().Unix())

			fmt.Println(a)
			keyset[i] = &Key{
				sb.Dataset[a].Address,
			}

		}
		//InputTypeB := sb.NormalCaseConstruction()
		//if set.Contains(InputTypeB[0].K) {
		//	log.Panicln("K existed")
		//}
		//set.Add(InputTypeB[0].K)
		err := enc.Encode(keyset)
		if err != nil {
			log.Fatal("encode error:", err)
		}
		tx.arbitraryData = encoder.Bytes()

	}
	hash, _ := tx.CalculateHash()
	copy(tx.hash[:], hash)
	tx.Sign(key)
	//fmt.Println(tx.arbitraryData)

	return tx
}

func (sb *Smallbank) GenerateBatchSmallbankTx(zp *generator.Zipfian, cfg cfg.Config, key *ecdsa.PrivateKey) []*Transaction {
	var txs []*Transaction
	s1 := mapset.NewSet()
	rand.Seed(time.Now().UnixNano())
	//fmt.Println(1)
	var (
		src = rand.NewSource(time.Now().UnixNano())
		r   = rand.New(src)
	)

	for i := 0; i < Txsperblock; i++ {
		tx := sb.GenerateSmallbankTx(key, s1, zp, r)
		txs = append(txs, tx)
	}

	return txs

}

func (sb *Smallbank) ExecuteBatchSmallbankTx(txs []*Transaction) {
	//var Commitlog []string

	//for _, tx := range txs {
	//	rw, _ := sb.SmallbankExecution(tx)
	//	cfg.UseRw(rw)
	//	//Commitlog = append(Commitlog, RwsetsToString(rwsets))
	//}
	//err := Commit(Commitlog)
	//if err != nil {
	//panic("Commit failed")
	//}
}

func (sb *Smallbank) Duplicate() *Smallbank {
	var (
		//syncmap sync.Map
		mp      map[string]int
		dataset Dataset
		scl     []SmallbankContractLog
	)
	mp = make(map[string]int, Smallbanksize)
	dataset = make(Dataset, Smallbanksize)
	scl = make([]SmallbankContractLog, Smallbanksize)
	//f := func(key, value interface{}) bool {
	//	syncmap.Store(key, value)
	//	return true
	//}
	//sb.Data.Range(f)
	for key, value := range sb.Mp {
		mp[key] = value
	}
	for index, value := range sb.Dataset {
		dataset[index] = &Smallbankdata{
			Data:    value.Data,
			Address: value.Address,
		}
	}
	for index, value := range sb.SCL {
		scl[index] = value
	}
	return &Smallbank{

		Mp:      mp,
		Dataset: dataset,
		SCL:     scl,
	}
}
