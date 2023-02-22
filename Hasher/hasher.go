package Hasher

import (
	"log"
)

type Trie interface {
	DeriveHash(Contents *[]Content) [32]byte
}

type NormalHasher struct {
	flag uint64
}

func NewNormalHasher() *NormalHasher {
	n := NormalHasher{1}
	return &n
}

func (hasher *NormalHasher) DeriveHash(Contents *[]Content) [32]byte {
	//fmt.Println(*Content)
	t, err := NewTree(*Contents)
	if err != nil {
		log.Fatal(err)
	}
	root := t.MerkleRoot()
	var root32 [32]byte
	copy(root32[:], root)
	return root32
}
