package src

import (
	"sync"
)

const NoIdx uint8 = 255

type Node struct {
	value   []byte
	key     uint64
	prevIdx uint8
	nextIdx uint8
}

type LRUMap struct {
	nodes    []Node
	freeList []uint8
	title    string
	keyToIdx map[uint64]uint8
	mutex    sync.RWMutex
	headIdx  uint8
	tailIdx  uint8
	capacity uint8
}
