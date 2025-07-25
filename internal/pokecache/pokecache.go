package pokecache

import (
	"sync"
	"time"
)

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

type Cache struct {
	table map[string]cacheEntry
	lock  sync.Mutex
}

func cleanCache(now time.Time, interval time.Duration, cache *Cache) {
	var dead []string
	cache.lock.Lock()
	defer cache.lock.Unlock()
	for key, val := range cache.table {
		if now.Sub(val.createdAt) > interval {
			dead = append(dead, key)
		}
	}
	for _, key := range dead {
		delete(cache.table, key)
	}
}

func reapLoop(interval time.Duration, cache *Cache) {
	ticker := time.NewTicker(interval)
	t := <-ticker.C
	for {
		cleanCache(t, interval, cache)
		t = <-ticker.C
	}
}

func NewCache(msInterval int) *Cache {
	cache := new(Cache)
	cache.table = make(map[string]cacheEntry)
	go reapLoop(time.Duration(msInterval)*time.Millisecond, cache)
	return cache
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.lock.Lock()
	defer c.lock.Unlock()
	val, ok := c.table[key]
	if ok {
		return val.val, true
	}
	return nil, false
}

func (c *Cache) Add(key string, val []byte) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.table[key] = cacheEntry{createdAt: time.Now(), val: val}

}
