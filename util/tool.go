package util

import "hash/fnv"

// String2Uint32 ...
func String2Uint32(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}
