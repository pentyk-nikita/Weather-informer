package cache

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type Cache struct {
	mu  sync.RWMutex
	dir string
	ttl time.Duration
}

type cacheItem struct {
	Data      []byte    `json:"data"`
	Timestamp time.Time `json:"timestamp"`
}

func New(dir string, ttl time.Duration) (*Cache, error) {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create cache directory: %w", err)
	}

	return &Cache{
		dir: dir,
		ttl: ttl,
	}, nil
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	filePath := filepath.Join(c.dir, key+".json")

	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, false
	}

	var item cacheItem
	if err := json.Unmarshal(data, &item); err != nil {
		return nil, false
	}

	if time.Since(item.Timestamp) > c.ttl {
		os.Remove(filePath)
		return nil, false
	}

	return item.Data, true
}

func (c *Cache) Set(key string, data []byte) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	item := cacheItem{
		Data:      data,
		Timestamp: time.Now(),
	}

	jsonData, err := json.Marshal(item)
	if err != nil {
		return fmt.Errorf("failed to marshal cache item: %w", err)
	}

	filePath := filepath.Join(c.dir, key+".json")
	return os.WriteFile(filePath, jsonData, 0644)
}