package core

import (
	"Excnahge-Cacher/api"
	"sync"
	"time"
)

type CacheItem struct {
	FromValue string
	ToValue   string
	Currency  float64
	ExpiresAt time.Time
}

type Cache struct {
	mu         sync.RWMutex
	data       map[string]CacheItem
	apiHandler *api.APIHandler
	stop       chan struct{}
}

func NewCache(handler *api.APIHandler) (*Cache, error) {
	cache := &Cache{
		data:       make(map[string]CacheItem),
		apiHandler: handler,
	}
	cache.StartCleanupWorker(30 * time.Second)
	err := cache.GenerateRates()
	if err != nil {
		return cache, err
	}

	return cache, nil
}

func (c *Cache) StartCleanupWorker(interval time.Duration) {
	ticker := time.NewTicker(interval)
	go func() {
		select {
		case <-ticker.C:
			c.cleanup()
		case <-c.stop:
			ticker.Stop()
			return
		}
	}()
}

func (c *Cache) cleanup() {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now()
	for key, item := range c.data {
		if now.After(item.ExpiresAt) {
			delete(c.data, key)
		}
	}
}

func (c *Cache) Stop() {
	close(c.stop)
}

func (c *Cache) GenerateRates() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	resp, err := c.apiHandler.GetAllCourses()
	if err != nil {
		return err
	}
	c.Drop()
	for key, value := range resp.Quotes {
		c.Set(key, value, time.Now().Add(5*time.Minute))
	}
	return nil
}

func InverseCourse(key string, value float64) (string, float64) {
	first := key[:3]
	second := key[2:]
	newKey := second + first
	newValue := 1. / value
	return newKey, newValue
}

func (c *Cache) Exists(key string) bool {
	_, exists := c.data[key]
	return exists
}

func (c *Cache) Set(key string, value float64, ttl time.Time) {
	c.mu.Lock()
	defer c.mu.Unlock()
	var currency, reversedCurrency CacheItem
	currency.ExpiresAt = ttl
	reversedCurrency.ExpiresAt = ttl
	reversedCurrency.FromValue = key[:3]
	currency.FromValue = key[2:]
	reversedCurrency.ToValue = key[2:]
	currency.ToValue = key[:3]
	currency.Currency = value
	newKey, newValue := InverseCourse(key, value)
	reversedCurrency.Currency = newValue
	if !c.Exists(key) {
		c.data[key] = currency
	}
	if !c.Exists(newKey) {
		c.data[newKey] = reversedCurrency
	}
}

func (c *Cache) Get(key string) CacheItem {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.data[key]
}

func (c *Cache) GetAll() []CacheItem {
	c.mu.RLock()
	defer c.mu.RUnlock()
	items := make([]CacheItem, 40)
	for _, value := range c.data {
		items = append(items, value)
	}
	return items
}

func (c *Cache) Drop() {
	c.mu.RLock()
	defer c.mu.RUnlock()
	for key, _ := range c.data {
		delete(c.data, key)
	}
}
