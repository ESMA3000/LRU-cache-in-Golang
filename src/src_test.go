package src

import (
	"fmt"
	"runtime"
	"sync"
	"testing"
)

func TestSrc(t *testing.T) {
	t.Run("LRUMap", func(t *testing.T) {
		t.Run("InitLRUMap", func(t *testing.T) {
			cache := InitLRUMap("test", 2)
			if cache.capacity != 2 {
				t.Errorf("Expected capacity 2, got %d", cache.capacity)
			}
			if cache.headIdx != NoIdx || cache.tailIdx != NoIdx {
				t.Error("Expected NoIdx for head and tail")
			}
			if len(cache.nodes) != 2 {
				t.Errorf("Expected nodes array length 2, got %d", len(cache.nodes))
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
			if node.prevIdx != NoIdx || node.nextIdx != NoIdx {
				t.Error("Expected NoIdx for prev and next indices")
			}
		})

		t.Run("removeNode", func(t *testing.T) {
			cache := InitLRUMap("test", 2)
			cache.Put(1, []byte("one"))
			if idx, ok := cache.keyToIdx[1]; ok {
				node := cache.getNodePtr(idx)
				cache.removeNode(node)
				if _, exists := cache.keyToIdx[1]; exists {
					t.Error("Expected key to be removed from keyToIdx")
				}
			} else {
				t.Error("Expected key 1 to exist in keyToIdx")
			}
		})

		t.Run("setHead", func(t *testing.T) {
			cache := InitLRUMap("test", 2)
			cache.Put(1, []byte("one"))
			cache.Put(2, []byte("two"))

			firstIdx := cache.headIdx
			secondIdx := cache.nodes[firstIdx].nextIdx

			if firstIdx == NoIdx || secondIdx == NoIdx {
				t.Error("Expected valid indices for head and next node")
			}

			if string(cache.nodes[firstIdx].value) != "two" {
				t.Error("Expected 'two' at head")
			}

			if string(cache.nodes[secondIdx].value) != "one" {
				t.Error("Expected 'one' as second node")
			}
		})

		t.Run("removeTail", func(t *testing.T) {
			cache := InitLRUMap("test", 2)
			cache.Put(1, []byte("one"))
			cache.Put(2, []byte("two"))
			oldTailIdx, ok := cache.removeTail()
			if !ok || cache.tailIdx == oldTailIdx {
				t.Error("Expected successful tail removal")
			}
		})

		t.Run("Put", func(t *testing.T) {
			cache := InitLRUMap("test", 2)
			cache.Put(1, []byte("one"))
			if cache.Length() != 1 {
				t.Error("Expected length 1")
			}
			if cache.nodes[cache.headIdx].key != 1 {
				t.Error("Expected key 1 at head")
			}
		})

		t.Run("GetNode", func(t *testing.T) {
			cache := InitLRUMap("test", 2)
			cache.Put(1, []byte("one"))
			node := cache.GetNode(1)
			if node == nil || node.key != 1 {
				t.Error("Expected to get node with key 1")
			}
		})

		t.Run("Get", func(t *testing.T) {
			cache := InitLRUMap("test", 2)
			cache.Put(1, []byte("one"))
			val := cache.Get(1)
			if string(val) != "one" {
				t.Errorf("Expected 'one', got %s", string(val))
			}
		})

		t.Run("Put Update", func(t *testing.T) {
			cache := InitLRUMap("test", 2)
			cache.Put(1, []byte("one"))
			cache.Put(1, []byte("new"))
			if string(cache.Get(1)) != "new" {
				t.Error("Expected value to be updated")
			}
		})

		t.Run("Eject", func(t *testing.T) {
			cache := InitLRUMap("test", 2)
			cache.Put(1, []byte("one"))
			cache.Eject(1)
			if cache.Get(1) != nil {
				t.Error("Expected key to be ejected")
			}
		})

		t.Run("Length", func(t *testing.T) {
			cache := InitLRUMap("test", 2)
			cache.Put(1, []byte("one"))
			if cache.Length() != 1 {
				t.Errorf("Expected length 1, got %d", cache.Length())
			}
		})

		t.Run("Clear", func(t *testing.T) {
			cache := InitLRUMap("test", 2)
			cache.Put(1, []byte("one"))
			cache.Put(2, []byte("two"))
			cache.Clear()
			if cache.Length() != 0 {
				t.Error("Expected empty cache")
			}
		})
	})

	t.Run("LRUMap", func(t *testing.T) {
		t.Run("TestLRUMapConcurrency", func(t *testing.T) {
			cache := InitLRUMap("test", 16)
			var wg sync.WaitGroup
			iterations := 10000

			wg.Add(2)

			// Writer goroutine
			go func() {
				defer wg.Done()
				for i := 0; i < iterations; i++ {
					cache.Put(uint64(i), []byte("test"))
				}
			}()

			// Reader goroutine
			go func() {
				defer wg.Done()
				for i := 0; i < iterations; i++ {
					cache.Get(uint64(i))
				}
			}()

			wg.Wait()

			if cache.Length() > 16 {
				t.Error("Cache exceeded capacity")
			}
		})

	})
}

func BenchmarkLRUMap(b *testing.B) {
	// Test different cache sizes
	sizes := []uint8{4, 8, 16, 64, 256 - 1}
	operations := []int{100, 1000, 10000}

	for _, size := range sizes {
		cache := InitLRUMap("test", size)
		data := []byte("test-value-that-is-larger-than-32-bytes-to-simulate-real-world-data")

		// Basic operations
		b.Run(fmt.Sprintf("Put/Size-%d", size), func(b *testing.B) {
			for i := 0; b.Loop(); i++ {
				cache.Put(uint64(i), data)
			}
		})

		b.Run(fmt.Sprintf("Get/Size-%d", size), func(b *testing.B) {
			for i := 0; b.Loop(); i++ {
				cache.Get(uint64(i % int(size)))
			}
		})

		// Concurrent operations
		for _, numOps := range operations {
			b.Run(fmt.Sprintf("Concurrent/Size-%d/Ops-%d", size, numOps), func(b *testing.B) {
				var wg sync.WaitGroup
				workers := runtime.GOMAXPROCS(0)      // Use system CPU count
				workChan := make(chan int, workers*4) // Increased buffer

				// Pre-warm worker pool
				for w := 0; w < workers*2; w++ {
					wg.Add(1)
					go func(workerID int) {
						defer wg.Done()
						for j := range workChan {
							key := uint64((workerID * numOps) + j)
							if workerID < workers {
								cache.Put(key, data)
							} else {
								cache.Get(key)
							}
						}
					}(w)
				}

				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					for j := 0; j < numOps; j++ {
						workChan <- j
					}
				}
				close(workChan)
				wg.Wait()
			})
		}

		// Mixed workload
		b.Run(fmt.Sprintf("MixedWorkload/Size-%d", size), func(b *testing.B) {
			ops := []string{"put", "get", "get", "get", "put"} // 40% writes, 60% reads
			for i := 0; b.Loop(); i++ {
				op := ops[i%len(ops)]
				key := uint64(i % int(size))

				switch op {
				case "put":
					cache.Put(key, data)
				case "get":
					cache.Get(key)
				}
			}
		})

		// Cache eviction
		b.Run(fmt.Sprintf("Eviction/Size-%d", size), func(b *testing.B) {
			for b.Loop() {
				// Force eviction by inserting size+1 items
				for j := 0; j <= int(size); j++ {
					cache.Put(uint64(j), data)
				}
			}
		})
		workers := 4
		b.Run("ConcurrentMixed", func(b *testing.B) {
			var wg sync.WaitGroup
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				wg.Add(workers)
				for w := 0; w < workers; w++ {
					go func(id int) {
						defer wg.Done()
						key := uint64(id)
						if id%2 == 0 {
							cache.Put(key, data)
						} else {
							cache.Get(key)
						}
					}(w)
				}
				wg.Wait()
			}
		})
	}
}
