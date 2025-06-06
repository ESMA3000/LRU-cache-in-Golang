package src

import "fmt"

type CacheManager struct {
	caches map[uint64]*LRUMap
}

func NewCacheManager() *CacheManager {
	return &CacheManager{
		caches: make(map[uint64]*LRUMap),
	}
}

func (cm *CacheManager) CreateCache(title string, key uint64, capacity uint8) {
	var cache *LRUMap = InitLRUMap(title, capacity)
	cm.caches[key] = cache
}

func (cm *CacheManager) GetCache(name uint64) *LRUMap {
	if cache, exists := cm.caches[name]; exists {
		return cache
	}
	return nil
}

func (cm *CacheManager) DestroyCache(name uint64) {
	if cache, exists := cm.caches[name]; exists {
		cache.Clear()
		delete(cm.caches, name)
	}
}

func (cm *CacheManager) ClearAllCaches() {
	for name := range cm.caches {
		cm.DestroyCache(name)
	}
}

func (cm *CacheManager) ListCaches() []string {
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
