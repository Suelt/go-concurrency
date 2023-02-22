package blockchain

import (
	"UndergraduateProj/block"
	"UndergraduateProj/config"
	"UndergraduateProj/receipt"
	"UndergraduateProj/tx"
	"sync"
	"sync/atomic"
)

type ExecutionEnv struct {
	sb        tx.Smallbankop
	Tcount    int
	H         *block.Header
	Txs       tx.Transactions
	Receipts  receipt.Receipts
	Sets      tx.ConcurrentSet
	logs      []*receipt.Receiptlog
	Stateroot [32]byte
	mutex     sync.Mutex
}

func (env *ExecutionEnv) ExecuteTransaction(transaction *tx.Transaction, index uint64) error {

	//fmt.Println()

	rpt, set, err := ApplyTransaction(&env.mutex, env.sb, transaction, env.H.Blocknumber, env.H.Blockroot, &env.H.GasUsed, index)
	if err != nil {
		return err
	}
	//fmt.Println(rpt.TransactionIndex)
	//env.mutex.Lock()
	env.Txs[index] = transaction
	env.Receipts[index] = rpt
	env.Sets[index] = set
	//env.mutex.Unlock()
	return nil
}

func ApplyTransaction(mutex *sync.Mutex, sb tx.Smallbankop, transaction *tx.Transaction, blocknumber uint64, blockroot [32]byte, totalgas *uint64, index uint64) (*receipt.Receipt, *tx.Rwsets, error) {
	//	tx.Timewasting()
	//set, _ := transaction.Execute(env.sb)
	//fmt.Println(index)
	set, _ := transaction.Execute(sb)

	txhash := transaction.Hash()
	contractaddress := config.ContractAddress
	datalength := transaction.SerializeLength()
	//
	gasused := transaction.Gasprice() * datalength
	atomic.AddUint64(totalgas, gasused)
	//*totalgas += gasused

	//log := &receipt.Receiptlog{[]byte{1, 2}}
	//*logs = append(*logs, log)

	receipt := &receipt.Receipt{
		TxHash:          txhash,
		ContractAddress: tx.StringToBytearray(contractaddress),
		GasUsed:         gasused,

		TransactionIndex:  index,
		BlockHash:         blockroot,
		BlockNumber:       blocknumber,
		StateRoot:         [32]byte{},
		Cumulativegasused: *totalgas,
		//Logs:              *logs,
	}

	//fmt.Println(blockroot)
	return receipt, set, nil
}
