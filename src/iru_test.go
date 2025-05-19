package src

import (
	"testing"
)

func TestLRUMap(t *testing.T) {
	t.Run("Basic Operations", func(t *testing.T) {
		cache := InitLRUMap(2)

		// Test Put and Get
		cache.Put(1, []byte("one"))
		if val := cache.Get(1); string(val) != "one" {
			t.Errorf("Expected 'one', got %s", string(val))
		}

		// Test capacity and eviction
		cache.Put(2, []byte("two"))
		cache.Put(3, []byte("three"))
		if val := cache.Get(1); val != nil {
			t.Error("Expected key 1 to be evicted")
		}
	})

	t.Run("Update Existing", func(t *testing.T) {
		cache := InitLRUMap(2)
		cache.Put(1, []byte("one"))
		cache.Put(1, []byte("new-one"))

		if val := cache.Get(1); string(val) != "new-one" {
			t.Errorf("Expected 'new-one', got %s", string(val))
		}
	})

	t.Run("Clear Cache", func(t *testing.T) {
		cache := InitLRUMap(2)
		cache.Put(1, []byte("one"))
		cache.Put(2, []byte("two"))
		cache.Clear()

		if len(cache.nodes) != 0 {
			t.Errorf("Expected empty cache, got size %d", len(cache.nodes))
		}
	})

	t.Run("LRU Order", func(t *testing.T) {
		cache := InitLRUMap(2)
		cache.Put(1, []byte("one"))
		cache.Put(2, []byte("two"))

		// Access 1, making it most recently used
		cache.Get(1)

		// Add new item, should evict 2
		cache.Put(3, []byte("three"))

		if val := cache.Get(2); val != nil {
			t.Error("Expected key 2 to be evicted")
		}
		if val := cache.Get(1); val == nil {
			t.Error("Expected key 1 to be present")
		}
	})

	t.Run("Eject", func(t *testing.T) {
		cache := InitLRUMap(2)
		cache.Put(1, []byte("one"))
		cache.Eject(1)

		if val := cache.Get(1); val != nil {
			t.Error("Expected key 1 to be ejected")
		}
	})
}

func BenchmarkLRUMap(b *testing.B) {
	cache := InitLRUMap(100)
	data := []byte("test-value")

	b.Run("Put", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			cache.Put(uint64(i), data)
		}
	})

	b.Run("Get", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			cache.Get(uint64(i % 100))
		}
	})

	b.Run("Mixed", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			if i%2 == 0 {
				cache.Put(uint64(i), data)
			} else {
				cache.Get(uint64(i - 1))
			}
		}
	})
}
