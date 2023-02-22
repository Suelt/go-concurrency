package blockchain

import (
	"UndergraduateProj/Hasher"
	"UndergraduateProj/block"
	cfg "UndergraduateProj/config"
	"UndergraduateProj/receipt"
	"UndergraduateProj/tx"
	"bytes"
	"encoding/gob"
	"fmt"
	mapset "github.com/deckarep/golang-set"
	"log"
	"runtime"
	"runtime/debug"
	"sync"
	"time"
)

type ConcurrentProcessor struct {
	bc     *Blockchain
	cfg    cfg.ConcurrentProcessorConfig
	engine Engine
}

type SequentialProcessor struct {
	bc     *Blockchain
	cfg    cfg.ProcessorConfig
	engine Engine
}

type PreProcessor struct {
	bc     *Blockchain
	cfg    cfg.ConcurrentProcessorConfig
	engine Engine
}

var tcons time.Duration

type Keyset struct {
	readkeys  []string
	writekeys []string
}

type ConcurrentUnit struct {
	transaction *tx.Transaction
	wait        []chan bool
	send        []chan bool
}

type Clue struct {
	lastWrite            int
	readlistBetweenWrite []int
}

func (Ccp *ConcurrentProcessor) Process(block *block.Block) (receipt.Receipts, uint64, error, time.Duration, int64) {
	//return Ccp.ProcessBlockWithoutDependencies(block)
	//return nil, 0, nil
	//var (
	//	t1 time.Time
	//	t2 time.Duration
	//)
	//t1 := time.Now()
	//rw := Ccp.FetchRwSetsByProcessing(block)
	//ke := Ccp.MapRwsetstoKeys(rw)
	//for _, value := range ke {
	//	fmt.Printf("%x           %x\n", value.readkeys[0], value.writekeys[0])
	//}
	//t2 := time.Since(t1)
	//fmt.Println(t2)
	//t1 := time.Now()

	env := Ccp.bc.MakeEnvWithBlock(block, Ccp.bc.Sb)
	//t3 := time.Now()
	//t3 := time.Now()
	ConcurrentUnits, tDetect := Ccp.DetectDependencies(block.Transactions)
	//t4 := time.Since(t3)
	//fmt.Println(t4)
	//t4 := time.Since(t3)
	//fmt.Println(t4)
	//cfg.UseRw(ConcurrentUnits)

	receipts, gasused, err := Ccp.ExecuteBatchTxsWithDependencies(env, ConcurrentUnits, block.Position)

	//cfg.UseRw(keysets)
	//t2 := time.Since(t1)
	//for _, value := range keysets {
	//	fmt.Printf("%x            %x\n", value.readkeys[0], value.writekeys[0])
	//}
	//fmt.Println(t2)
	//err = StoreWholeState(Ccp.bc.DBState, Ccp.bc.Sb.(*tx.Smallbank))
	//if err != nil {
	//	log.Panicln(err)
	//}
	return receipts, gasused, err, tcons, tDetect
}

//type ConcurrentUnit struct {
//transaction *tx.Transaction
//wait        []chan bool
//send        []chan bool
//}
//
//type Clue struct {
//	lastWrite            int
//	readlistBetweenWrite []int
//}

func (Ccp *ConcurrentProcessor) ExecuteBatchTxsWithDependencies(env *ExecutionEnv, units []*ConcurrentUnit, blockposition *block.Blockposition) (receipt.Receipts, uint64, error) {
	//receipts := make(receipt.Receipts, len(units))
	var wg sync.WaitGroup
	//	var a int32
	//debug.SetMaxThreads(8)
	runtime.GOMAXPROCS(8)
	t1 := time.Now()
	for index, unit := range units {
		wg.Add(1)
		indexcopy := index
		unitcopy := unit
		go func() {
			var wgwait sync.WaitGroup
			wgwait.Add(len(unitcopy.wait))
			for i := 0; i < len(unitcopy.wait); i++ {
				icopy := i
				go func() {
					<-unitcopy.wait[icopy]

					wgwait.Done()
				}()
			}
			wgwait.Wait()
			//data, errnew := Ccp.bc.DBNormal.Get([]byte{1, 2, 3, 4, 5, 6, 7, 8}, nil)
			//cfg.UseRw(data, errnew)
			_ = env.ExecuteTransaction(units[indexcopy].transaction, uint64(indexcopy))
			//errhash := StoreTransactionByHash(Ccp.bc.DBNormal, unitcopy.transaction)
			//if errhash != nil {
			//	log.Println(errhash)
			//}
			//errpos := StoreTransactionByPosition(Ccp.bc.DBNormal, unitcopy.transaction, blockposition, uint64(indexcopy))
			//if errpos != nil {
			//	log.Println(errpos)
			//}
			////set := env.Sets[indexcopy]
			////rkv := set.Rkv
			//wkv := env.Sets[indexcopy].Wkv
			//for i := 0; i < len(wkv); i++ {
			//	//dread := &tx.Smallbankdata{
			//	//	rkv[i].K,
			//	//	rkv[i].V,
			//	//}
			//	//if len(wkv) == 0 {
			//	//	break
			//	//}
			//	dwrite := &tx.Smallbankdata{
			//		wkv[i].K,
			//		wkv[i].V,
			//	}
			//	if dwrite.Address == "" {
			//		continue
			//	}
			//	errstate := StoreState(Ccp.bc.DBState, dwrite.Address, Ccp.bc.Sb.(*tx.Smallbank).Mp[dwrite.Address], dwrite)
			//	if errstate != nil {
			//		log.Panicln(errstate)
			//	}
			//	//var cread Hasher.Content
			//	//cread = *dread
			//	//var cwrite Hasher.Content
			//	//cwrite = *dwrite
			//	//Ccp.bc.StateTrie.UpdateTree(cread, cwrite)
			//	//fmt.Printf("%x\n", Ccp.bc.StateTrie.Root.Hash)
			//}

			for i := 0; i < len(unitcopy.send); i++ {
				//fmt.Println(indexcopy)
				unitcopy.send[i] <- true
				close(unitcopy.send[i])
			}
			//atomic.AddInt32(&a, 1)
			//fmt.Println(a)
			wg.Done()
		}()
	}
	wg.Wait()
	t2 := time.Since(t1)
	tcons = t2
	//fmt.Println("time used: " + fmt.Sprint(t2))

	blockcontents := make([]Hasher.Content, tx.Smallbanksize)
	dataset := Ccp.bc.Sb.(*tx.Smallbank).Dataset
	for index, value := range dataset {
		blockcontents[index] = *value

	}
	errnew := Ccp.bc.StateTrie.RebuildTreeWith(blockcontents)
	cfg.UseRw(errnew)
	return env.Receipts, env.H.GasUsed, nil
}

func (Ccp *ConcurrentProcessor) DetectDependencies(txs tx.Transactions) ([]*ConcurrentUnit, int64) {
	//t1 := time.Now()

	rwsets := Ccp.FetchRwSetsByProcessing(txs)
	keysets := Ccp.MapRwsetstoKeys(rwsets)
	//	t2 := time.Since(t1)
	//fmt.Println(t2)

	keys := Ccp.GetAccessedkeys(keysets)
	dependencymap := make(map[string]*Clue, len(keys))
	ConcurrentUnits := make([]*ConcurrentUnit, len(keysets))
	t1 := time.Now().UnixNano()
	for i := 0; i < len(keys); i++ {
		dependencymap[keys[i]] = &Clue{
			-1,
			[]int{},
		}

	}
	for i := 0; i < len(keysets); i++ {
		ConcurrentUnits[i] = &ConcurrentUnit{
			transaction: txs[i],
		}
	}
	for i := 0; i < len(txs); i++ {
		readkeys := keysets[i].readkeys
		writekeys := keysets[i].writekeys
		for j := 0; j < len(writekeys); j++ {
			if writekeys[j] == "" {
				continue
			}
			readlist := dependencymap[writekeys[j]].readlistBetweenWrite
			for k := 0; k < len(readlist); k++ {
				chread_write_dependency := make(chan bool)

				ConcurrentUnits[readlist[k]].send = append(ConcurrentUnits[readlist[k]].send, chread_write_dependency)
				ConcurrentUnits[i].wait = append(ConcurrentUnits[i].wait, chread_write_dependency)
			}
			dependencymap[writekeys[j]].readlistBetweenWrite = []int{}
			last := dependencymap[writekeys[j]].lastWrite
			if last == -1 {
				dependencymap[writekeys[j]].lastWrite = i
				continue
			}
			if last == i {
				continue
			}
			chwrite_write_dependency := make(chan bool)
			ConcurrentUnits[last].send = append(ConcurrentUnits[last].send, chwrite_write_dependency)
			ConcurrentUnits[i].wait = append(ConcurrentUnits[i].wait, chwrite_write_dependency)
			dependencymap[writekeys[j]].lastWrite = i
		}
		for j := 0; j < len(readkeys); j++ {

			lastwrite := dependencymap[readkeys[j]].lastWrite
			if lastwrite == -1 {
				dependencymap[readkeys[j]].readlistBetweenWrite = append(dependencymap[readkeys[j]].readlistBetweenWrite, i)
				continue
			}
			if lastwrite == i {
				continue
			}
			dependencymap[readkeys[j]].readlistBetweenWrite = append(dependencymap[readkeys[j]].readlistBetweenWrite, i)
			ch_write_read_dependency := make(chan bool)
			ConcurrentUnits[lastwrite].send = append(ConcurrentUnits[lastwrite].send, ch_write_read_dependency)
			ConcurrentUnits[i].wait = append(ConcurrentUnits[i].wait, ch_write_read_dependency)
		}

	}
	t2 := time.Now().UnixNano()

	return ConcurrentUnits, t2 - t1
}

func (Ccp *ConcurrentProcessor) GetAccessedkeys(keysets []*Keyset) []string {
	s := mapset.NewSet()
	for _, set := range keysets {
		for i := 0; i < len(set.readkeys); i++ {
			s.Add(set.readkeys[i])
		}
		for i := 0; i < len(set.writekeys); i++ {
			s.Add(set.writekeys[i])
		}
	}
	strings := s.ToSlice()
	stringscopy := make([]string, len(strings))
	for i := 0; i < len(strings); i++ {
		stringscopy[i] = strings[i].(string)
	}
	return stringscopy

}

func (Stp *SequentialProcessor) Process(b *block.Block) (receipt.Receipts, uint64, error, time.Duration, int64) {
	//blockroot := block.Header.Blockroot
	//blocknumber := block.Header.

	env := Stp.bc.MakeEnvWithBlock(b, Stp.bc.Sb)
	t1 := time.Now()
	for index, transaction := range b.Transactions {
		//data, errnew := Stp.bc.DBNormal.Get([]byte{1, 2, 3, 4, 5, 6, 7, 8}, nil)
		//cfg.UseRw(data, errnew)
		_ = env.ExecuteTransaction(transaction, uint64(index))
		//err := StoreTransactionByHash(Stp.bc.DBNormal, transaction)
		//if err != nil {
		//	log.Panicln(err)
		//}
		//errpos := StoreTransactionByPosition(Stp.bc.DBNormal, transaction, b.Position, uint64(index))
		//if errpos != nil {
		//	log.Println(errpos)
		//}
		//set := env.Sets[index]
		////rkv := set.Rkv
		//wkv := set.Wkv
		//
		//for i := 0; i < len(wkv); i++ {
		//	//dread := &tx.Smallbankdata{
		//	//	rkv[i].K,
		//	//	rkv[i].V,
		//	//}
		//	if len(wkv) == 0 {
		//		break
		//	}
		//	dwrite := &tx.Smallbankdata{
		//		wkv[i].K,
		//		wkv[i].V,
		//	}
		//	if dwrite.Address == "" {
		//		continue
		//	}
		//	errstate := StoreState(Stp.bc.DBState, dwrite.Address, Stp.bc.Sb.(*tx.Smallbank).Mp[dwrite.Address], dwrite)
		//	if errstate != nil {
		//		log.Panicln(errstate)
		//	}
		//	//var cread Hasher.Content
		//	//cread = *dread
		//	//var cwrite Hasher.Content
		//	//cwrite = *dwrite
		//	//Stp.bc.StateTrie.UpdateTree(cread, cwrite, Stp.bc.Sb.(*tx.Smallbank).Mp[dwrite.Address])
		//	//fmt.Printf("%x\n", Stp.bc.StateTrie.Root.Hash)
		//}

	}
	t2 := time.Since(t1)
	cfg.UseRw(t2)
	//fmt
	fmt.Println(t2)
	blockcontents := make([]Hasher.Content, tx.Smallbanksize)
	dataset := Stp.bc.Sb.(*tx.Smallbank).Dataset
	for index, value := range dataset {
		blockcontents[index] = *value

	}
	errnew := Stp.bc.StateTrie.RebuildTreeWith(blockcontents)
	cfg.UseRw(errnew)
	//blockcontents := make([]Hasher.Content, tx.Smallbanksize)
	//dataset := Stp.bc.Sb.(*tx.Smallbank).Dataset
	//for index, value := range dataset {
	//	blockcontents[index] = *value
	//
	//}
	//err := Stp.bc.StateTrie.RebuildTreeWith(blockcontents)
	//t2 = time.Since(t1)

	//Stp.bc.Sb.CommittoState(env.Sets)
	//fmt.Println(env.H.GasUsed)
	return env.Receipts, env.H.GasUsed, nil, t2, 0

}

func (Ccp *ConcurrentProcessor) ProcessBatchBlocks(chain []*block.Block) ([]receipt.Receipts, []uint64, error) {
	return nil, nil, nil
}

func (Ccp *ConcurrentProcessor) ProcessBlockWithoutDependencies(b *block.Block) (receipt.Receipts, uint64, error) {
	env := Ccp.bc.MakeEnvWithBlock(b, Ccp.bc.Sb)

	var wg sync.WaitGroup
	transactions := b.Transactions
	debug.SetMaxThreads(16)
	t1 := time.Now()

	for index, tx := range transactions {
		wg.Add(1)
		indexcopy := index
		txcopy := tx
		go Ccp.f(txcopy, indexcopy, &wg, env)
	}
	wg.Wait()
	t2 := time.Since(t1)
	cfg.UseRw(t2)
	//cfg.UseRw(t2)
	//fmt
	//fmt.Println(t2)
	fmt.Println(len(env.Sets))
	//Ccp.bc.Sb.CommittoState(env.Sets)
	//fmt.Println(env.H.GasUsed)

	return env.Receipts, env.H.GasUsed, nil
	//return nil, 0, nil
}
func (Ccp *ConcurrentProcessor) f(transaction *tx.Transaction, index int, wg *sync.WaitGroup, env *ExecutionEnv) {
	_ = env.ExecuteTransaction(transaction, uint64(index))

	wg.Done()
}

func (Ccp *ConcurrentProcessor) FetchRwSetsByProcessing(txs tx.Transactions) []*tx.Rwsets {
	b := &block.Block{
		Transactions: txs,
		Hash:         [32]byte{},
		Header:       &block.Header{Blocknumber: 0},
	}
	env := Ccp.bc.MakeEnvWithBlock(b, Ccp.bc.Sb.Duplicate())
	var wg sync.WaitGroup
	transactions := txs
	for index, tx := range transactions {
		wg.Add(1)
		indexcopy := index
		txcopy := tx
		go Ccp.f(txcopy, indexcopy, &wg, env)
	}
	wg.Wait()

	//Ccp.bc.Sb.CommittoState(env.Sets)
	//fmt.Println(env.H.GasUsed)

	return env.Sets
	//sbcopy := Ccp.bc.Sb.Duplicate()
	//env := Ccp.bc.MakeEnvWithBlock(b, sbcopy)

}

func (Ccp *ConcurrentProcessor) FetchRwsetsByStaticAnalysis(txs tx.Transactions) []*tx.Rwsets {
	transactions := txs
	sets := make([]*tx.Rwsets, len(transactions))
	for index, transaction := range transactions {
		receiver := transaction.Receiver()
		var bytecopy []byte
		for _, value := range receiver {
			if value == 0 {
				break
			}
			bytecopy = append(bytecopy, value)
		}
		StringEquals := string(bytecopy)
		data := transaction.ArbitraryData()
		var decoder bytes.Buffer
		dec := gob.NewDecoder(&decoder)
		decoder.Write(data)
		if StringEquals == "Transfer" {
			var transferlist [cfg.SmallbankAccessedkeys]*tx.Transfer
			err := dec.Decode(&transferlist)
			if err != nil {
				log.Fatal("decode error:", err)
			}
			epoch := 0
			var (
				rkv = make([]tx.Set, 2*cfg.SmallbankAccessedkeys)
				wkv = make([]tx.Set, 2*cfg.SmallbankAccessedkeys)
			)
			for i := 0; i < cfg.SmallbankAccessedkeys; i++ {

				rkv[i+epoch].K = transferlist[i].From
				rkv[i+epoch+1].K = transferlist[i].To

				wkv[i+epoch].K = transferlist[i].From
				wkv[i+epoch+1].K = transferlist[i].To
				epoch++
			}
			sets[index] = &tx.Rwsets{
				rkv,
				wkv,
			}
		} else {
			var Depset [cfg.SmallbankAccessedkeys]*tx.Key
			err := dec.Decode(&Depset)
			if err != nil {
				log.Fatal("decode error:", err)
			}
			var (
				rkv = make([]tx.Set, cfg.SmallbankAccessedkeys)
				wkv = make([]tx.Set, cfg.SmallbankAccessedkeys)
			)
			for i := 0; i < cfg.SmallbankAccessedkeys; i++ {

				rkv[i].K = Depset[i].K
				wkv[i].K = Depset[i].K
			}
			sets[index] = &tx.Rwsets{
				rkv,
				wkv,
			}
		}
	}
	return sets
}
func (Ccp *ConcurrentProcessor) MapRwsetstoKeys(sets []*tx.Rwsets) []*Keyset {
	retsets := make([]*Keyset, len(sets))
	var (
		readkeys  []string
		writekeys []string
	)
	for index, rwset := range sets {
		readkeys = []string{}
		writekeys = []string{}
		rkv := rwset.Rkv
		wkv := rwset.Wkv
		for _, keyvalue := range rkv {
			readkeys = append(readkeys, keyvalue.K)
		}
		for _, keyvalue := range wkv {
			writekeys = append(writekeys, keyvalue.K)
		}
		retsets[index] = &Keyset{
			readkeys:  readkeys,
			writekeys: writekeys,
		}
	}
	return retsets
}
func NewConcurrentProcessor(bc *Blockchain) *ConcurrentProcessor {
	return &ConcurrentProcessor{
		bc:     bc,
		cfg:    cfg.ConcurrentProcessorConfig{},
		engine: bc.Engine,
	}
}

func NewSequentialProcessor(bc *Blockchain) *SequentialProcessor {
	return &SequentialProcessor{
		bc:     bc,
		cfg:    cfg.ProcessorConfig{},
		engine: bc.Engine,
	}
}
