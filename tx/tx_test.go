package tx

import (
	"testing"
)

//func TestSerialize(t *testing.T) {
//	tx := InitTx("123123131231312312", "Deposit", "12312312312312")
//
//	Bytes, _ := tx.Serialize()
//	fmt.Println(len(Bytes))
//	Hash, _ := tx.HashTx()
//	fmt.Println(len(Hash))
//
//}
//
//func TestSmallbank(t *testing.T) {
//	var sb Smallbank
//	sb.Init()
//	tx := InitTx("123123131231312312", "Deposit", "12312312312312")
//	hash, _ := tx.HashTx()
//	tx.hash = hash
//
//}

func TestSmallbankExecution(t *testing.T) {
	//var sb Smallbank
	//y0 := time.Now()
	//sb.Init()
	//curve := elliptic.P256()
	//privkey, _ := ecdsa.GenerateKey(curve, rand.Reader)
	//y1 := time.Since(y0)
	//fmt.Printf("time cost=%v\n", y1)
	//x0 := time.Now()
	//txs := sb.GenerateBatchSmallbankTx(cfg.ConcurrentProcessorConfig{}, privkey)
	//
	//x1 := time.Since(x0)
	//fmt.Printf("time cost=%v\n", x1)
	//
	//fmt.Println("generation txs success")
	//
	//x2 := time.Now()
	//////fmt.Println(txs)
	//sb.ExecuteBatchSmallbankTx(txs)
	//x3 := time.Since(x2)
	//fmt.Println(x3)

}

//func TestConcurrentMap(t *testing.T) {
//	var syncmap sync.Map
//	syncmap.Store("qcrao", 18)
//	syncmap.Store("stefno", 20)
//	go func() {
//		for {
//			time.Sleep(100 * time.Millisecond)
//			cock, _ := syncmap.Load("qcrao")
//			fmt.Println(cock)
//		}
//	}()
//	go func() {
//		for {
//			time.Sleep(100 * time.Millisecond)
//			syncmap.Store("qcrao", 20)
//
//		}
//	}()
//	select {}
//}

//func TestStr(t *testing.T) {
//
//	tx := InitTx("sender", "Deposit", "123")
//	receiver := tx.receiver
//	var bytecopy []byte
//	for _, value := range receiver {
//		if value == 0 {
//			break
//		}
//		bytecopy = append(bytecopy, value)
//	}
//	var str string
//	str = string(bytecopy)
//	fmt.Println(len(bytecopy))
//	fmt.Println(len(str))
//	//receiver = bytes.TrimSuffix(receiver, []byte{10})[:]
//	fmt.Println(len(str))
//	if strings.Compare(str, "Deposit") == 0 {
//		fmt.Println("Equals")
//	}
//}

func TestRandTx(t *testing.T) {
	//var mm map[int]map[int]string
	//config.UseRw()

	//rand.Seed(86)
	//iterator := rand.Intn(100000)
	//iterator2 := rand.Intn(100000)
	//fmt.Println(iterator, iterator2)

	//var sb Smallbank
	//sb.Init()
	//fmt.Printf("%x\n", sb.GenerateStateRoot())
	//fmt.Printf("%x\n", sb.GenerateStateRoot())
	//Generateindex()
}
