package tx

import (
	"UndergraduateProj/Hasher"
	cfg "UndergraduateProj/config"
	"crypto/ecdsa"
	mapset "github.com/deckarep/golang-set"
	"github.com/pingcap/go-ycsb/pkg/generator"
	"math/rand"
)

//type SmallbankInterface struct {
//	smallbankop
//}

type Smallbankop interface {
	Init()
	Deposit(depset [cfg.SmallbankAccessedkeys]*Key) (*Rwsets, error)
	Withdraw(depset [cfg.SmallbankAccessedkeys]*Key) (*Rwsets, error)
	Query(depset [cfg.SmallbankAccessedkeys]*Key) (*Rwsets, error)
	Transfer(transferlist [cfg.SmallbankAccessedkeys]*Transfer) (*Rwsets, error)
	SmallbankExecution(transaction *Transaction) (*Rwsets, error)
	GenerateSmallbankTx(*ecdsa.PrivateKey, mapset.Set, *generator.Zipfian, *rand.Rand) *Transaction

	GenerateBatchSmallbankTx(zipfian *generator.Zipfian, cfg cfg.Config, key *ecdsa.PrivateKey) []*Transaction
	GenerateStateRoot() [32]byte
	WriteState(rw *Rwsets)
	CommittoState(sets []*Rwsets)
	Duplicate() *Smallbank
	GenerateStatetrie() *Hasher.MerkleTree
	//	ReverseState(rw []*Rwsets)
}

//func (SbI *SmallbankInterface) SmallbankExecution(transaction *Transaction) {
//
//}
//
//func (SbI *SmallbankInterface) Init() {
//	SbI.Init()
//}
//
//func (SbI *SmallbankInterface) Deposit(transaction *Transaction) (Rwsets, error) {
//	return SbI.Deposit(transaction)
//}
//func (SbI *SmallbankInterface) Withdraw(transaction *Transaction) (Rwsets, error) {
//	return SbI.Withdraw(transaction)
//}
//func (SbI *SmallbankInterface) Query(transaction *Transaction) (Rwsets, error) {
//	return SbI.Query(transaction)
//}
//func (SbI *SmallbankInterface) Transfer(transaction *Transaction) (Rwsets, error) {
//	return SbI.Transfer(transaction)
//}
