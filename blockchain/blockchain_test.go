package blockchain

import (
	"testing"
)

func TestBlockchainFormation(t *testing.T) {
	blockchain := InitBlockchain()
	//Concurrent := NewConcurrentProcessor(blockchain)
	//Sequential := NewSequentialProcessor(blockchain)
	blockchain.StateProcessor = NewSequentialProcessor(blockchain)
	//blockchain.StateProcessor = NewConcurrentProcessor(blockchain)
	blockchain.GenerateMainWorkFlow()
	//fmt.Println(b.Header.Position)
	//position := &block.Blockposition{
	//	b.Header.Position.Line,
	//	b.Header.Position.List - 1,
	//}
	//parent := blockchain.engine.GetBlock(b.Header.Parenthash, position)
	//
	//fmt.Println(parent.Hash)
	//fmt.Println(b.Header.Parenthash)

	//bc := NewBlockchain()
	//var blocks block.Blocks
	//for _, chain := range bc.Chains {
	//	blocks = append(blocks, chain.Currentblock)
	//}
	//bc.Generate(blocks)

}

//
//func TestDetectDependency(t *testing.T) {
//
//	blockchain := InitBlockchain()
//	blockchain.StateProcessor = NewConcurrentProcessor(blockchain)
//	curve := elliptic.P256()
//	privkey, _ := ecdsa.GenerateKey(curve, rand.Reader)
//	bp := NewBcProposer(blockchain)
//	var txslist []tx.Transactions
//	var keysetslist [][]*Keyset
//	Ccp := NewConcurrentProcessor(blockchain)
//	for i := 0; i < 100; i++ {
//		txs := blockchain.Sb.GenerateBatchSmallbankTx(bp.cfg.Zipfian, bp.cfg, privkey)
//		txslist = append(txslist, txs)
//		rwsets := Ccp.FetchRwSetsByProcessing(txs)
//		keysets := Ccp.MapRwsetstoKeys(rwsets)
//		keysetslist = append(keysetslist, keysets)
//
//	}
//	t1 := time.Now()
//	for m := 0; m < 100; m++ {
//		keysets := keysetslist[m]
//		txs := txslist[m]
//		keys := Ccp.GetAccessedkeys(keysets)
//		config.UseRw(keys, txs)
//		dependencymap := make(map[string]*Clue, len(keys))
//		ConcurrentUnits := make([]*ConcurrentUnit, len(keysets))
//		//t1 := time.Now().UnixNano()
//		for i := 0; i < len(keys); i++ {
//			dependencymap[keys[i]] = &Clue{
//				-1,
//				[]int{},
//			}
//
//		}
//		for i := 0; i < len(keysets); i++ {
//			ConcurrentUnits[i] = &ConcurrentUnit{
//				transaction: txs[i],
//			}
//		}
//		for i := 0; i < len(txs); i++ {
//			readkeys := keysets[i].readkeys
//			writekeys := keysets[i].writekeys
//
//			for j := 0; j < len(writekeys); j++ {
//				if writekeys[j] == "" {
//					continue
//				}
//				readlist := dependencymap[writekeys[j]].readlistBetweenWrite
//				for k := 0; k < len(readlist); k++ {
//
//					chread_write_dependency := make(chan bool)
//
//					ConcurrentUnits[readlist[k]].send = append(ConcurrentUnits[readlist[k]].send, chread_write_dependency)
//					ConcurrentUnits[i].wait = append(ConcurrentUnits[i].wait, chread_write_dependency)
//				}
//				dependencymap[writekeys[j]].readlistBetweenWrite = []int{}
//				last := dependencymap[writekeys[j]].lastWrite
//				if last == -1 {
//					dependencymap[writekeys[j]].lastWrite = i
//					continue
//				}
//				if last == i {
//					continue
//				}
//				chwrite_write_dependency := make(chan bool)
//				ConcurrentUnits[last].send = append(ConcurrentUnits[last].send, chwrite_write_dependency)
//				ConcurrentUnits[i].wait = append(ConcurrentUnits[i].wait, chwrite_write_dependency)
//				dependencymap[writekeys[j]].lastWrite = i
//			}
//			for j := 0; j < len(readkeys); j++ {
//
//				lastwrite := dependencymap[readkeys[j]].lastWrite
//				if lastwrite == -1 {
//					dependencymap[readkeys[j]].readlistBetweenWrite = append(dependencymap[readkeys[j]].readlistBetweenWrite, i)
//					continue
//				}
//				if lastwrite == i {
//					continue
//				}
//				ch_write_read_dependency := make(chan bool)
//
//				ConcurrentUnits[lastwrite].send = append(ConcurrentUnits[lastwrite].send, ch_write_read_dependency)
//				ConcurrentUnits[i].wait = append(ConcurrentUnits[i].wait, ch_write_read_dependency)
//			}
//
//		}
//		//t2 := time.Now().UnixNano()
//
//	}
//
//	t2 := time.Since(t1)
//	log.Println(t2)
//
//	t3 := time.Now()
//
//	for m := 0; m < 100; m++ {
//		keysets := keysetslist[m]
//		txs := txslist[m]
//		keys := Ccp.GetAccessedkeys(keysets)
//		dependencymap := make(map[string]*Clue, len(keys))
//		ConcurrentUnits := make([]*ConcurrentUnit, len(keysets))
//		config.UseRw(keys, txs, dependencymap, ConcurrentUnits)
//	}
//
//	t4 := time.Since(t3)
//	log.Println(t4)
//}
