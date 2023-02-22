package tx

import (
	"encoding/binary"
	"fmt"
	"math/rand"
)

//Print transferinfo like Account A transfers 10.0 to account B
//a single transaction may access several keys and retStr concatenates the above transferring info and returns
func RwsetsToString(rw Rwsets) string {

	var value float64
	var retStr string
	rkv := rw.Rkv
	wkv := rw.Wkv

	retStr = ""
	var rk, wk string

	for index, log := range rkv {
		rk = log.K
		value = log.V
		format := "read%d:read key:%x,read value:%f"
		str := fmt.Sprintf(format, index, rk, value)
		retStr = retStr + str + "\n"
	}

	for index, log := range wkv {
		wk = log.K
		value = log.V
		format := "write%d:write key:%x,write value:%f"
		str := fmt.Sprintf(format, index, wk, value)
		retStr = retStr + str + "\n"
	}

	return retStr

}

func InterfaceToFloat64(unk interface{}) float64 {
	return unk.(float64)

}
func StringToBytearray(str string) [32]byte {
	b := []byte(str)
	bcopy32 := [32]byte{}
	copy(bcopy32[:], b)
	return bcopy32
}

func Randfloat(min float64, max float64) float64 {
	return min + rand.Float64()*(max-min)
}

func Uint64toByte(value uint64) []byte {
	n := make([]byte, 8)
	binary.BigEndian.PutUint64(n, value)
	return n
}

func Randint(min uint64, max uint64) uint64 {
	return uint64(uint64(rand.Intn(int(max-min))) + min)
}
