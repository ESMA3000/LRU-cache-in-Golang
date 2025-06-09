package src

import "sync"

type Uints interface {
	~uint8 | ~uint16 | ~uint32 | ~uint64
}

type Node[U, K Uints, V any] struct {
	value   V
	key     K
	prevIdx U
	nextIdx U
}

type LRUMap[U, K Uints, V any] struct {
	nodes    []Node[U, K, V]
	freeList []U
	title    string
	keyToIdx map[K]U
	mutex    sync.RWMutex
	headIdx  U
	tailIdx  U
	NoIdx    U
	capacity U
}

type CacheManager[U, K Uints, V any] struct {
	caches map[K]*LRUMap[U, K, V]
}
