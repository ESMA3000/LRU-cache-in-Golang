package src

import (
	"fmt"
	"strings"
	"sync"
)

var nodePool = sync.Pool{
	New: func() any {
		return &Node{}
	},
}

type Node struct {
	key   uint64
	value []byte
	prev  *Node
	next  *Node
}

type LRUMap struct {
	capacity uint8
	head     *Node
	tail     *Node
	nodes    map[uint64]*Node
	mutex    sync.Mutex
}

func newNode(key uint64, value []byte) *Node {
	node := nodePool.Get().(*Node)
	node.key = key
	node.value = value
	node.prev = nil
	node.next = nil
	return node
}

func (m *LRUMap) removeNode(node *Node) {
	if node == nil {
		return
	}
	delete(m.nodes, node.key)
	*node = Node{}
	nodePool.Put(node)
}

func (m *LRUMap) unlinkNode(node *Node) {
	if node.prev != nil {
		node.prev.next = node.next
	}
	if node.next != nil {
		node.next.prev = node.prev
	}
}

func InitLRUMap(capacity uint8) *LRUMap {
	for range int(capacity) {
		nodePool.Put(&Node{})
	}
	return &LRUMap{
		capacity: capacity,
		nodes:    make(map[uint64]*Node, capacity),
	}
}

func (m *LRUMap) setHead(node *Node) {
	if m.head == nil {
		m.head = node
		m.tail = node
		return
	} else if m.head == node || node == nil {
		return
	}

	if m.tail == node {
		m.tail = node.prev
	}

	m.unlinkNode(node)
	node.prev = nil
	node.next = m.head
	m.head.prev = node
	m.head = node
}

func (m *LRUMap) removeTail() {
	if m.tail == nil {
		return
	}
	var newTail *Node = m.tail.prev
	if newTail != nil {
		newTail.next = nil
	}
	m.removeNode(m.tail)
	m.tail = newTail
}

func (m *LRUMap) addNode(key uint64, value []byte) {
	m.nodes[key] = newNode(key, value)
	m.setHead(m.nodes[key])

	if uint8(len(m.nodes)) > m.capacity {
		m.removeTail()
	}
}

func (m *LRUMap) GetNode(key uint64) *Node {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	if node, ok := m.nodes[key]; ok {
		m.setHead(node)
		return node
	}
	return nil
}

func (m *LRUMap) Get(key uint64) []byte {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	if node, ok := m.nodes[key]; ok {
		m.setHead(node)
		return node.value
	}
	return nil
}

func (m *LRUMap) Put(key uint64, value []byte) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	if node, ok := m.nodes[key]; ok {
		node.value = value
		m.setHead(node)
	} else {
		m.addNode(key, value)
	}
}

func (m *LRUMap) Eject(key uint64) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	if node, ok := m.nodes[key]; ok {
		m.removeNode(node)
	}
}

func (m *LRUMap) Length() uint8 {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	return uint8(len(m.nodes))
}

func (m *LRUMap) Clear() {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	for key := range m.nodes {
		m.removeNode(m.nodes[key])
	}
}

func (m *LRUMap) Print() string {
	m.PrintNodes()
	var builder strings.Builder
	for _, node := range m.nodes {
		builder.WriteString(fmt.Sprintf("Key: %d, Value: %v\n",
			node.key, node.value))
	}
	return builder.String()
}

func (m *LRUMap) PrintNodes() {
	fmt.Println(m.head, m.tail)
	for _, node := range m.nodes {
		fmt.Println(node)
	}
}
