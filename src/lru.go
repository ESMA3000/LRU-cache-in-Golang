// Package src implements a fixed-size LRU (Least Recently Used) cache
package src

// InitLRUMap initializes a new LRU cache with given title and capacity
func InitLRUMap(title string, capacity uint8) *LRUMap {
	if capacity >= NoIdx {
		capacity = NoIdx - 1
	}
	m := &LRUMap{
		title:    title,
		capacity: capacity,
		nodes:    make([]Node, capacity),
		keyToIdx: make(map[uint64]uint8, capacity),
		freeList: make([]uint8, capacity),
		headIdx:  NoIdx,
		tailIdx:  NoIdx,
	}

	for i := range m.nodes {
		m.freeList[i] = uint8(i)
	}
	return m
}

// newNode creates a new cache node with the given key and value
func newNode(key uint64, value []byte) Node {
	return Node{
		key:     key,
		value:   value,
		prevIdx: NoIdx,
		nextIdx: NoIdx,
	}
}

// Internal node management methods

func (m *LRUMap) getNodePtr(idx uint8) *Node {
	return &m.nodes[idx]
}

func (m *LRUMap) getFreeIndex() (uint8, bool) {
	if len(m.freeList) == 0 {
		return 0, false
	}
	idx := m.freeList[len(m.freeList)-1]
	m.freeList = m.freeList[:len(m.freeList)-1]
	return idx, true
}

func (m *LRUMap) removeNode(node *Node) {
	delete(m.keyToIdx, node.key)
	node.prevIdx = NoIdx
	node.nextIdx = NoIdx
}

func (m *LRUMap) setHead(idx uint8) {
	if m.headIdx == idx {
		return
	}
	if m.headIdx == NoIdx {
		m.headIdx = idx
		m.tailIdx = idx
		return
	}
	node := m.getNodePtr(idx)
	if m.tailIdx == idx {
		m.tailIdx = node.prevIdx
	}

	m.unlinkNode(node)
	node.prevIdx = NoIdx
	node.nextIdx = m.headIdx
	m.nodes[m.headIdx].prevIdx = idx
	m.headIdx = idx
}

func (m *LRUMap) removeTail() (uint8, bool) {
	if m.tailIdx == NoIdx {
		return NoIdx, false
	}

	oldTailIdx := m.tailIdx
	tailNode := m.getNodePtr(m.tailIdx)

	m.tailIdx = tailNode.prevIdx

	if m.headIdx == oldTailIdx {
		m.headIdx = NoIdx
	}

	if m.tailIdx != NoIdx {
		m.nodes[m.tailIdx].nextIdx = NoIdx
	}

	tailNode.prevIdx = NoIdx
	tailNode.nextIdx = NoIdx

	return oldTailIdx, true
}

func (m *LRUMap) unlinkNode(node *Node) {
	if node.prevIdx != NoIdx {
		m.nodes[node.prevIdx].nextIdx = node.nextIdx
	}
	if node.nextIdx != NoIdx {
		m.nodes[node.nextIdx].prevIdx = node.prevIdx
	}
}

// Public API methods

// Put adds or updates a key-value pair in the cache
func (m *LRUMap) Put(key uint64, value []byte) {
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

	m.nodes[idx] = newNode(key, value)
	m.keyToIdx[key] = idx
	m.setHead(idx)
}

// Get retrieves a value from the cache by key
func (m *LRUMap) Get(key uint64) []byte {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if idx, ok := m.keyToIdx[key]; ok {
		m.setHead(idx)
		return m.getNodePtr(idx).value
	}
	return nil
}

// Eject removes a key-value pair from the cache
func (m *LRUMap) Eject(key uint64) {
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
func (m *LRUMap) GetNode(key uint64) *Node {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	if idx, ok := m.keyToIdx[key]; ok {
		m.setHead(idx)
		return m.getNodePtr(idx)
	}
	return nil
}

// Length returns the current number of items in the cache
func (m *LRUMap) Length() uint8 {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return uint8(len(m.keyToIdx))
}

// Clear removes all items from the cache
func (m *LRUMap) Clear() {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	for i := range m.nodes {
		m.nodes[i] = newNode(0, nil)
	}
	for i := range m.freeList {
		m.freeList[i] = uint8(i)
	}
	for k := range m.keyToIdx {
		delete(m.keyToIdx, k)
	}
	m.headIdx = NoIdx
	m.tailIdx = NoIdx
}

// Iterator returns a slice of nodes in order (or reverse order)
func (m *LRUMap) Iterator(rev bool) []*Node {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	nodes := make([]*Node, 0, m.Length())
	if len(m.keyToIdx) == 0 || m.headIdx == NoIdx {
		return nodes
	}

	var curr uint8
	if rev {
		curr = m.tailIdx
	} else {
		curr = m.headIdx
	}
	for curr != NoIdx {
		nodes = append(nodes, m.getNodePtr(curr))
		if rev {
			curr = m.nodes[curr].prevIdx
		} else {
			curr = m.nodes[curr].nextIdx
		}
	}
	return nodes
}
