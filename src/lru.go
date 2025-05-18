package src

import (
	"fmt"
	"strings"
)

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
}

func newNode(key uint64, value []byte) *Node {
	return &Node{
		key:   key,
		value: value,
		prev:  nil,
		next:  nil,
	}
}

func InitLRUMap(capacity uint8) LRUMap {
	return LRUMap{
		capacity: capacity,
		head:     nil,
		tail:     nil,
		nodes:    make(map[uint64]*Node, capacity),
	}
}

func (c *LRUMap) setHead(node *Node) {
	if c.head == nil {
		c.head = node
		c.tail = node
		return
	}
	if c.head == node {
		return
	}

	if node.prev != nil {
		node.prev.next = node.next
	}
	if node.next != nil {
		node.next.prev = node.prev
	}

	node.prev = nil
	node.next = c.head
	c.head.prev = node
	c.head = node
}

func (c *LRUMap) removeTail() {
	var newTail *Node = c.tail.prev
	newTail.next = nil
	c.removeNode(c.tail)
	c.tail = newTail
}

func (c *LRUMap) addNode(key uint64, value []byte) {
	c.nodes[key] = newNode(key, value)
	c.setHead(c.nodes[key])

	if uint8(len(c.nodes)) > c.capacity {
		c.removeTail()
	}
}

func (c *LRUMap) removeNode(node *Node) {
	delete(c.nodes, node.key)
}

func (c LRUMap) GetNode(key uint64) *Node {
	if node, ok := c.nodes[key]; ok {
		c.setHead(node)
		return node
	}
	return nil
}

func (c LRUMap) Get(key uint64) []byte {
	if node, ok := c.nodes[key]; ok {
		c.setHead(node)
		return node.value
	}
	return nil
}

func (c *LRUMap) Put(key uint64, value []byte) {
	if node, ok := c.nodes[key]; ok {
		node.value = value
		c.setHead(node)
	} else {
		c.addNode(key, value)
	}
}

func (c *LRUMap) Eject(key uint64) {
	if node, ok := c.nodes[key]; ok {
		c.removeNode(node)
	}
}

func (c *LRUMap) Clear() {
	for key := range c.nodes {
		delete(c.nodes, key)
	}
}

func (c LRUMap) Print() string {
	c.PrintNodes()
	var builder strings.Builder
	for _, node := range c.nodes {
		builder.WriteString(fmt.Sprintf("Key: %d, Value: %v\n",
			node.key, node.value))
	}
	return builder.String()
}

func (c LRUMap) PrintNodes() {
	for _, node := range c.nodes {
		fmt.Println(node)
	}
}
