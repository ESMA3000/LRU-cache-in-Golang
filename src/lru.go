// Package src implements a fixed-size LRU (Least Recently Used) cache
package src

import "sync"

// InitLRUMap initializes a new LRU cache with given title and capacity
func InitLRUMap[U, K Uints, V any](title string, capacity U) *LRUMap[U, K, V] {
	if capacity >= ^U(0) {
		capacity = ^U(0) - 1
	}
	m := &LRUMap[U, K, V]{
		title:    title,
		capacity: capacity,
		nodes:    make([]Node[U, K, V], capacity),
		keyToIdx: make(map[K]U, capacity),
		freeList: make([]U, capacity),
		mutex:    sync.RWMutex{},
		NoIdx:    ^U(0),
		headIdx:  ^U(0),
		tailIdx:  ^U(0),
	}

	for i := range m.nodes {
		m.freeList[i] = U(i)
	}
	return m
}

// Internal node management methods
func (m *LRUMap[U, K, V]) newNode(key K, value V) Node[U, K, V] {
	return Node[U, K, V]{
		key:     key,
		value:   value,
		prevIdx: m.NoIdx,
		nextIdx: m.NoIdx,
	}
}

func (m *LRUMap[U, K, V]) getNodePtr(idx U) *Node[U, K, V] {
	return &m.nodes[idx]
}

func (m *LRUMap[U, K, V]) getFreeIndex() (U, bool) {
	if len(m.freeList) == 0 {
		return 0, false
	}
	idx := m.freeList[len(m.freeList)-1]
	m.freeList = m.freeList[:len(m.freeList)-1]
	return idx, true
}

func (m *LRUMap[U, K, V]) removeNode(node *Node[U, K, V]) {
	delete(m.keyToIdx, node.key)
	node.prevIdx = m.NoIdx
	node.nextIdx = m.NoIdx
}

func (m *LRUMap[U, K, V]) setHead(idx U) {
	if m.headIdx == idx {
		return
	}
	if m.headIdx == m.NoIdx {
		m.headIdx = idx
		m.tailIdx = idx
		return
	}
	node := m.getNodePtr(idx)
	if m.tailIdx == idx {
		m.tailIdx = node.prevIdx
	}

	m.unlinkNode(node)
	node.prevIdx = m.NoIdx
	node.nextIdx = m.headIdx
	m.nodes[m.headIdx].prevIdx = idx
	m.headIdx = idx
}

func (m *LRUMap[U, K, V]) removeTail() (U, bool) {
	if m.tailIdx == m.NoIdx {
		return m.NoIdx, false
	}

	oldTailIdx := m.tailIdx
	tailNode := m.getNodePtr(m.tailIdx)

	m.tailIdx = tailNode.prevIdx

	if m.headIdx == oldTailIdx {
		m.headIdx = m.NoIdx
	}

	if m.tailIdx != m.NoIdx {
		m.nodes[m.tailIdx].nextIdx = m.NoIdx
	}

	tailNode.prevIdx = m.NoIdx
	tailNode.nextIdx = m.NoIdx

	return oldTailIdx, true
}

func (m *LRUMap[U, K, V]) unlinkNode(node *Node[U, K, V]) {
	if node.prevIdx != m.NoIdx {
		m.nodes[node.prevIdx].nextIdx = node.nextIdx
	}
	if node.nextIdx != m.NoIdx {
		m.nodes[node.nextIdx].prevIdx = node.prevIdx
	}
}

// Public API methods

// Put adds or updates a key-value pair in the cache
func (m *LRUMap[U, K, V]) Put(key K, value V) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if existingIdx, ok := m.keyToIdx[key]; ok {
		m.nodes[existingIdx].value = value
		m.setHead(existingIdx)
		return
	}

	idx, ok := m.getFreeIndex()
	if !ok {
		if tailIdx, ok := m.removeTail(); ok {
			delete(m.keyToIdx, m.nodes[tailIdx].key)
			idx = tailIdx
		}
	}

	m.nodes[idx] = m.newNode(key, value)
	m.keyToIdx[key] = idx
	m.setHead(idx)
}

// Get retrieves a value from the cache by key
func (m *LRUMap[U, K, V]) Get(key K) V {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if idx, ok := m.keyToIdx[key]; ok {
		m.setHead(idx)
		return m.getNodePtr(idx).value
	}
	var zero V
	return zero
}

// Eject removes a key-value pair from the cache
func (m *LRUMap[U, K, V]) Eject(key K) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if idx, ok := m.keyToIdx[key]; ok {
		node := m.getNodePtr(idx)
		if idx == m.headIdx {
			m.headIdx = node.nextIdx
		}
		if idx == m.tailIdx {
			m.tailIdx = node.prevIdx
		}
		m.unlinkNode(node)
		m.removeNode(node)
		m.freeList = append(m.freeList, idx)
	}
}

// GetNode retrieves a node from the cache by key
func (m *LRUMap[U, K, V]) GetNode(key K) *Node[U, K, V] {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	if idx, ok := m.keyToIdx[key]; ok {
		m.setHead(idx)
		return m.getNodePtr(idx)
	}
	return nil
}

// Length returns the current number of items in the cache
func (m *LRUMap[U, K, V]) Length() U {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return U(len(m.keyToIdx))
}

// Clear removes all items from the cache
func (m *LRUMap[U, K, V]) Clear() {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	for i := range m.nodes {
		m.nodes[i] = m.newNode(K(0), *new(V))
	}
	for i := range m.freeList {
		m.freeList[i] = U(i)
	}
	for k := range m.keyToIdx {
		delete(m.keyToIdx, k)
	}
	m.headIdx = m.NoIdx
	m.tailIdx = m.NoIdx
}

// Iterator returns a slice of nodes in order (or reverse order)
func (m *LRUMap[U, K, V]) Iterator(rev bool) []*Node[U, K, V] {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	nodes := make([]*Node[U, K, V], 0, m.Length())
	if len(m.keyToIdx) == 0 || m.headIdx == m.NoIdx {
		return nodes
	}

	var curr U
	if rev {
		curr = m.tailIdx
	} else {
		curr = m.headIdx
	}
	for curr != m.NoIdx {
		nodes = append(nodes, m.getNodePtr(curr))
		if rev {
			curr = m.nodes[curr].prevIdx
		} else {
			curr = m.nodes[curr].nextIdx
		}
	}
	return nodes
}
