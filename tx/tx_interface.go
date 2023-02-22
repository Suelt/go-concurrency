package tx

import (
	"crypto/ecdsa"
	"github.com/cbergoon/merkletree"
)

type TxInterface struct {
	txnormal
	merkletree.Content
}

type txnormal interface {
	Serialize() ([32]byte, error)
	Hash() ([32]byte, error)
	Verify() bool
	Sign(ecdsa.PrivateKey)
	DetectDependencies() []chan bool
	AddIntoChanList(ch chan bool)
	Execute() error
}

//type TrieNodeContent interface {
//	CalculateHash() ([]byte, error)
//	Equals(other merkletree.Content) (bool, error)
//}
