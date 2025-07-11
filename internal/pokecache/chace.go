package pokecache

import (
	"sync"
	"time"
)

type Cache struct {
	cache map[string]cacheEntry
	mu	 	*sync.Mutex
}
type cacheEntry struct {
	createdAt time.Time
	val 			[]byte
}

func NewCache(interval time.Duration) *Cache{
	c := &Cache{
		cache: make(map[string]cacheEntry),
		mu: 	 &sync.Mutex{}, 
	}
	go c.reapLoop(interval)
	return c
}

func (c *Cache) Add(key string, val []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cache[key] = cacheEntry{
		createdAt: time.Now(),
		val:       val,
	}
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	entry, ok := c.cache[key]
	if !ok {
		return nil, false
	}
	return entry.val, true
}

func (c *Cache) reapLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		c.mu.Lock()
		for key, entry := range c.cache {
			if time.Since(entry.createdAt) > interval {
				delete(c.cache, key)
			}
		}
		c.mu.Unlock()
	}
}