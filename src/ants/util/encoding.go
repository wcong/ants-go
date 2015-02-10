package util

import (
	"hash/fnv"
)

func HashString(encode string) uint64 {
	hash := fnv.New64()
	hash.Write([]byte(encode))
	return hash.Sum64()
}
