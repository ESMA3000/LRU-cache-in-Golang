package src

import (
	"testing"
)

func TestLRUMap(t *testing.T) {
	t.Run("InitLRUMap", func(t *testing.T) {
		cache := InitLRUMap(2)
		if cache.capacity != 2 {
			t.Errorf("Expected capacity 2, got %d", cache.capacity)
		}
		if cache.head != nil || cache.tail != nil {
			t.Error("Expected empty head and tail")
		}
		if len(cache.nodes) != 0 {
			t.Error("Expected empty nodes map")
		}
	})

	t.Run("newNode", func(t *testing.T) {
		node := newNode(1, []byte("test"))
		if node.key != 1 {
			t.Errorf("Expected key 1, got %d", node.key)
		}
		if string(node.value) != "test" {
			t.Errorf("Expected value 'test', got %s", string(node.value))
		}
		if node.prev != nil || node.next != nil {
			t.Error("Expected nil prev and next pointers")
		}
	})

	t.Run("removeNode", func(t *testing.T) {
		cache := InitLRUMap(2)
		cache.Put(1, []byte("one"))
		node := cache.nodes[1]
		cache.removeNode(node)
		if _, exists := cache.nodes[1]; exists {
			t.Error("Expected node to be removed from map")
		}
	})

	t.Run("setHead", func(t *testing.T) {
		cache := InitLRUMap(2)
		node := newNode(1, []byte("one"))

		// Test empty cache
		cache.setHead(node)
		if cache.head != node || cache.tail != node {
			t.Error("Expected node to be both head and tail")
		}

		// Test adding second node
		node2 := newNode(2, []byte("two"))
		cache.setHead(node2)
		if cache.head != node2 || cache.head.next != node {
			t.Error("Expected node2 to be head and linked to node1")
		}
	})

	t.Run("removeTail", func(t *testing.T) {
		cache := InitLRUMap(2)
		cache.Put(1, []byte("one"))
		cache.Put(2, []byte("two"))
		oldTail := cache.tail
		cache.removeTail()
		if cache.tail == oldTail {
			t.Error("Expected tail to change")
		}
	})

	t.Run("addNode", func(t *testing.T) {
		cache := InitLRUMap(2)
		cache.addNode(1, []byte("one"))
		if cache.Length() != 1 {
			t.Error("Expected length 1")
		}
		if cache.head.key != 1 {
			t.Error("Expected key 1 at head")
		}
	})

	t.Run("GetNode", func(t *testing.T) {
		cache := InitLRUMap(2)
		cache.Put(1, []byte("one"))
		node := cache.GetNode(1)
		if node == nil || node.key != 1 {
			t.Error("Expected to get node with key 1")
		}
	})

	t.Run("Get", func(t *testing.T) {
		cache := InitLRUMap(2)
		cache.Put(1, []byte("one"))
		val := cache.Get(1)
		if string(val) != "one" {
			t.Errorf("Expected 'one', got %s", string(val))
		}
	})

	t.Run("Put Update", func(t *testing.T) {
		cache := InitLRUMap(2)
		cache.Put(1, []byte("one"))
		cache.Put(1, []byte("new"))
		if string(cache.Get(1)) != "new" {
			t.Error("Expected value to be updated")
		}
	})

	t.Run("Eject", func(t *testing.T) {
		cache := InitLRUMap(2)
		cache.Put(1, []byte("one"))
		cache.Eject(1)
		if cache.Get(1) != nil {
			t.Error("Expected key to be ejected")
		}
	})

	t.Run("Length", func(t *testing.T) {
		cache := InitLRUMap(2)
		cache.Put(1, []byte("one"))
		if cache.Length() != 1 {
			t.Errorf("Expected length 1, got %d", cache.Length())
		}
	})

	t.Run("Clear", func(t *testing.T) {
		cache := InitLRUMap(2)
		cache.Put(1, []byte("one"))
		cache.Put(2, []byte("two"))
		cache.Clear()
		if cache.Length() != 0 {
			t.Error("Expected empty cache")
		}
	})

}

func BenchmarkNodePool(b *testing.B) {
	b.Run("NodeAllocation", func(b *testing.B) {
		for i := 0; b.Loop(); i++ {
			node := newNode(uint64(i), []byte("test"))
			nodePool.Put(node)
		}
	})
}
