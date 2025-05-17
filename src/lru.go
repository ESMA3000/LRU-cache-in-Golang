package src

import (
	"fmt"
)

type Node struct {
	key   string
	value any
	next  *Node
	prev  *Node
}

type LRUCache struct {
	capacity int
	nodes    map[string]*Node
}

func newNode(key string, value any) *Node {
	return &Node{
		key:   key,
		value: value,
		next:  nil,
		prev:  nil,
	}
}

func InitLRU(capacity int) LRUCache {
	return LRUCache{
		capacity: capacity,
		nodes:    make(map[string]*Node, capacity),
	}
}

func (c LRUCache) findHead() *Node {
	if len(c.nodes) == 0 {
		return nil
	}

	for _, node := range c.nodes {
		if node.prev == nil && node.next != nil {
			return node
		}
	}

	if len(c.nodes) == 1 {
		for _, node := range c.nodes {
			return node
		}
	}

	for _, node := range c.nodes {
		node.prev = nil
		return node
	}

	return nil
}

func (c LRUCache) findTail() *Node {
	for _, node := range c.nodes {
		if node.next == nil {
			return node
		}
	}
	return nil
}

func (c LRUCache) setHead(node *Node) {
	var currHead *Node = c.findHead()
	if len(c.nodes) == 1 {
		if currHead != node {
			node.next = currHead
			currHead.prev = node
			node.prev = nil
		}
		return
	}

	if currHead == node {
		return
	}

	if node.prev != nil {
		node.prev.next = node.next
	}
	if node.next != nil {
		node.next.prev = node.prev
	}

	node.prev = nil
	node.next = currHead
	currHead.prev = node
}

func (c LRUCache) removeTail() {
	var currTail *Node = c.findTail()
	var newTail *Node = currTail.prev
	if newTail != nil {
		newTail.next = nil
	}
	c.removeNode(currTail)
}

func (c LRUCache) addNode(key string, value any) {
	c.nodes[key] = newNode(key, value)
	c.setHead(c.nodes[key])

	if len(c.nodes) > c.capacity {
		c.removeTail()
	}
}

func (c LRUCache) removeNode(node *Node) {
	delete(c.nodes, node.key)
}

func (c LRUCache) GetNode(key string) *Node {
	if node, ok := c.nodes[key]; ok {
		c.setHead(node)
		return node
	}
	return nil
}

func (c LRUCache) Get(key string) any {
	if node, ok := c.nodes[key]; ok {
		c.setHead(node)
		return node.value
	}
	return nil
}

func (c LRUCache) Put(key string, value any) {
	if node, ok := c.nodes[key]; ok {
		node.value = value
	} else {
		c.addNode(key, value)
	}
}

func (c LRUCache) Eject(key string) {
	if node, ok := c.nodes[key]; ok {
		c.removeNode(node)
	}
}

func (c LRUCache) Clear() {
	for key := range c.nodes {
		delete(c.nodes, key)
	}
}

func (c LRUCache) Print() {
	for _, node := range c.nodes {
		fmt.Printf("Key: %s, Value: %v\n", node.key, node.value)
	}
}
