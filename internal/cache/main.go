package cache

import (
	"sync"
	"time"
)

type Cache struct {
	Data map[string]CacheEntry
	mu   sync.Mutex
}

type CacheEntry struct {
	createdAt time.Time
	Val       []byte
}

func (c *Cache) Add(key string, val []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.Data[key] = CacheEntry{
		createdAt: time.Now(),
		Val:       val,
	}
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	data, ok := c.Data[key]

	if ok {
		return data.Val, true
	}

	return data.Val, false

}

func (c *Cache) reaploop(interval time.Duration) {
	ticker := time.NewTicker(5 * time.Second)

	for {
		_ = <-ticker.C

		c.mu.Lock()

		for key, value := range c.Data {
			now := time.Now()

			if now.Unix()-value.createdAt.Unix() > int64(interval) {
				delete(c.Data, key)
			}
		}

		c.mu.Unlock()

	}
}

func NewCache(duration time.Duration) *Cache {
	c := Cache{
		Data: map[string]CacheEntry{},
		mu:   sync.Mutex{},
	}

	go c.reaploop(duration)

	return &c
}
