package src

import "fmt"

func NewCacheManager[U, K Uints, V any]() *CacheManager[U, K, V] {
	return &CacheManager[U, K, V]{
		caches: make(map[K]*LRUMap[U, K, V]),
	}
}

func (cm *CacheManager[U, K, V]) CreateCache(title string, key K, capacity U) {
	var cache *LRUMap[U, K, V] = InitLRUMap[U, K, V](title, capacity)
	cm.caches[key] = cache
}

func (cm *CacheManager[U, K, V]) GetCache(name K) *LRUMap[U, K, V] {
	if cache, exists := cm.caches[name]; exists {
		return cache
	}
	return nil
}

func (cm *CacheManager[U, K, V]) DestroyCache(name K) {
	if cache, exists := cm.caches[name]; exists {
		cache.Clear()
		delete(cm.caches, name)
	}
}

func (cm *CacheManager[U, K, V]) ClearAllCaches() {
	for name := range cm.caches {
		cm.DestroyCache(name)
	}
}

func (cm *CacheManager[U, K, V]) ListCaches() []string {
	var names []string = make([]string, 0, len(cm.caches))
	for m := range cm.caches {
		cache := cm.caches[m]
		names = append(names, cache.title)
		for _, node := range cache.Iterator(false) {
			names = append(names, fmt.Sprintf("Key: %d, Value: %v", node.key, node.value))
		}
	}
	return names
}
