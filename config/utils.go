package config

func UseRw(vals ...interface{}) {
	for _, val := range vals {
		_ = val
	}
}

func MapUinttoBitlength(value uint64) uint64 {
	var (
		res   uint64 = 0
		multi uint64 = 1
	)
	for i := 0; i < 64; i++ {
		if value == multi {
			res = uint64(i)
			break
		}
		multi *= 2
	}
	return res
}
