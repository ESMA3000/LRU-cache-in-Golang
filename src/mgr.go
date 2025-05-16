package src

type CacheManager struct {
	caches map[string]*LRUCache
}

func NewCacheManager() *CacheManager {
	return &CacheManager{
		caches: make(map[string]*LRUCache),
	}
}

func (cm *CacheManager) CreateCache(name string, capacity int) {
	cache := InitLRU(capacity)
	cm.caches[name] = &cache
}

func (cm *CacheManager) GetCache(name string) *LRUCache {
	if cache, exists := cm.caches[name]; exists {
		return cache
	}
	return nil
}

func (cm *CacheManager) DeleteCache(name string) {
	if cache, exists := cm.caches[name]; exists {
		cache.Clear()
		delete(cm.caches, name)
	}
}

func (cm *CacheManager) ListCaches() []string {
	names := make([]string, 0, len(cm.caches))
	for name := range cm.caches {
		names = append(names, name)
	}
	return names
}
