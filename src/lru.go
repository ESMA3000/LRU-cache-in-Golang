package src

import (
	"fmt"
	"strings"
)

type Node struct {
	key   string
	value []byte
	prev  *Node
	next  *Node
}

type LRUCache struct {
	capacity uint8
	head     *Node
	tail     *Node
	nodes    map[string]*Node
}

func newNode(key string, value []byte) *Node {
	return &Node{
		key:   key,
		value: value,
		prev:  nil,
		next:  nil,
	}
}

func InitLRU(capacity uint8) LRUCache {
	return LRUCache{
		capacity: capacity,
		head:     nil,
		tail:     nil,
		nodes:    make(map[string]*Node, capacity),
	}
}

func (c *LRUCache) setHead(node *Node) {
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

func (c *LRUCache) removeTail() {
	var newTail *Node = c.tail.prev
	newTail.next = nil
	c.removeNode(c.tail)
	c.tail = newTail
}

func (c *LRUCache) addNode(key string, value []byte) {
	c.nodes[key] = newNode(key, value)
	c.setHead(c.nodes[key])

	if uint8(len(c.nodes)) > c.capacity {
		c.removeTail()
	}
}

func (c *LRUCache) removeNode(node *Node) {
	delete(c.nodes, node.key)
}

func (c LRUCache) GetNode(key string) *Node {
	if node, ok := c.nodes[key]; ok {
		c.setHead(node)
		return node
	}
	return nil
}

func (c LRUCache) Get(key string) []byte {
	if node, ok := c.nodes[key]; ok {
		c.setHead(node)
		return node.value
	}
	return nil
}

func (c *LRUCache) Put(key string, value []byte) {
	if node, ok := c.nodes[key]; ok {
		node.value = value
		c.setHead(node)
	} else {
		c.addNode(key, value)
	}
}

func (c *LRUCache) Eject(key string) {
	if node, ok := c.nodes[key]; ok {
		c.removeNode(node)
	}
}

func (c *LRUCache) Clear() {
	for key := range c.nodes {
		delete(c.nodes, key)
	}
}

func (c LRUCache) Print() string {
	c.PrintNodes()
	var builder strings.Builder
	for _, node := range c.nodes {
		builder.WriteString(fmt.Sprintf("Key: %s, Value: %v\n",
			node.key, node.value))
	}
	return builder.String()
}

func (c LRUCache) PrintNodes() {
	for _, node := range c.nodes {
		fmt.Println(node)
	}
}
