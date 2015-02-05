package util

import (
	"ants/node"
	"hash/fnv"
)

func HashNode(nodeInfo *node.Node) uint64 {
	hash := fnv.New64()
	hash.Write([]byte(nodeInfo.Ip + string(nodeInfo.Port)))
	return hash.Sum64()
}
