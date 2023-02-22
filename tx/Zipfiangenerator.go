package tx

import (
	"github.com/pingcap/go-ycsb/pkg/generator"
	"math/rand"
	"time"
)

const Zipfianconstant = 0.999

func Generateindex(zipfian *generator.ScrambledZipfian) uint64 {

	var ret uint64
	var (
		src = rand.NewSource(time.Now().Unix())
		r   = rand.New(src)
	)

	ret = uint64(zipfian.Next(r))
	//fmt.Println(ret)
	return ret
}
