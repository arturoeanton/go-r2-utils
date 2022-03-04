package hash

import "hash/fnv"

func HashStringUint64(s string) uint64 {
	h := fnv.New64()
	h.Write([]byte(s))
	return h.Sum64()
}
