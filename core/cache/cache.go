package core

import (
	"sync"
	"time"
)

type CacheItem struct {
	FromValue     string
	ToValue       string
	Course        float64
	InverseCourse float64
	ttl           time.Duration
}

type Cache struct {
	mu   sync.RWMutex
	data map[string]CacheItem
}

func NewCache() *Cache {
	return &Cache{
		data: make(map[string]CacheItem),
	}
}

func (c *Cache) FindAllCourses(response)

func (c *Cache) Set(key string, value CacheItem, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	value.ttl = ttl
	c.data[key] = value
}

func (c *Cache) Get(key string) CacheItem {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.data[key]
}
