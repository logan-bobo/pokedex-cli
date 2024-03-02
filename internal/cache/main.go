package cahce

import (
	"fmt"
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

func (c *Cache) Get(key string) (bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	_, ok := c.Data[key]

	if ok {
		return true
	}

	return false

}

func (c *Cache) reaploop(interval time.Duration) {
	ticker := time.NewTicker(5 * time.Second)

	for {
		_ = <-ticker.C

		c.mu.Lock()

		for key, value := range c.Data {
			now := time.Now()

			if now.Unix() - value.createdAt.Unix() > int64(interval) {
				delete(c.Data, key)
			}
		}

		c.mu.Unlock()

	}
}

func NewCache(duration time.Duration) (*Cache) {
	c := Cache{
		Data: map[string]CacheEntry{},
		mu:   sync.Mutex{},
	}

	go c.reaploop(60)

	return &c
}

