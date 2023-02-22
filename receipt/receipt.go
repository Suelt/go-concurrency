package receipt

import (
	"UndergraduateProj/Hasher"
	"UndergraduateProj/tx"
	"bytes"
	"crypto/sha256"
)

type Receipt struct {

	// Implementation fields: These fields are added by geth when processing a transaction.
	// They are stored in the chain database.
	TxHash            [32]byte
	ContractAddress   [32]byte
	GasUsed           uint64
	Cumulativegasused uint64
	// Inclusion information: These fields provide information about the inclusion of the
	// transaction corresponding to this receipt.
	BlockHash        [32]byte
	BlockNumber      uint64
	TransactionIndex uint64
	StateRoot        [32]byte
	Logs             []*Receiptlog
}

type Receiptlog struct {
	Data []byte
}

func (Rpt Receipt) CalculateHash() ([]byte, error) {
	h := sha256.New()
	ReceiptByteCopy := Rpt.ConverttoByte()
	if _, err := h.Write(ReceiptByteCopy); err != nil {
		return nil, err
	}

	return h.Sum(nil), nil
}
func (Rpt Receipt) Equals(other Hasher.Content) (bool, error) {
	return bytes.Equal(Rpt.ConverttoByte(), other.(Receipt).ConverttoByte()), nil
}

func (Rpt Receipt) ConverttoByte() []byte {
	s := [][]byte{Rpt.TxHash[:], Rpt.ContractAddress[:], Rpt.BlockHash[:], Rpt.StateRoot[:],
		tx.Uint64toByte(Rpt.GasUsed),
		tx.Uint64toByte(Rpt.BlockNumber),
		tx.Uint64toByte(Rpt.TransactionIndex)}
	return bytes.Join(s, []byte{})
}

func NewReceipt(root [32]byte) *Receipt {
	return &Receipt{StateRoot: root}
}

type ReceiptRlP struct {
	State         []byte
	cumulativegas uint64
	logs          []*Receiptlog
}
type Receipts []*Receipt

func (rs Receipts) Len() int {
	return len(rs)
}

// EncodeIndex encodes the i'th receipt to w.
//func (rs Receipts) EncodeIndex(i int, w *bytes.Buffer) {
//	r := rs[i]
//	data := &ReceiptRlP{
//		r.StateRoot[:], r.Cumulativegasused, r.Logs}
//	err := rlp.Encode(w, data)
//	if err != nil {
//		return
//	}
//}
