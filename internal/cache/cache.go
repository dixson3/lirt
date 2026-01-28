package cache

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/dixson3/lirt/internal/config"
)

// Cache represents a file-based cache
type Cache struct {
	profile string
	ttl     time.Duration
}

// CachedData represents cached data with metadata
type CachedData struct {
	FetchedAt time.Time   `json:"fetchedAt"`
	Data      interface{} `json:"data"`
}

// New creates a new cache instance
func New(profile string, ttl time.Duration) *Cache {
	return &Cache{
		profile: profile,
		ttl:     ttl,
	}
}

// GetCacheDir returns the cache directory for this profile
func (c *Cache) GetCacheDir() string {
	return filepath.Join(config.GetConfigDir(), "cache", c.profile)
}

// ensureCacheDir creates the cache directory if it doesn't exist
func (c *Cache) ensureCacheDir() error {
	dir := c.GetCacheDir()
	return os.MkdirAll(dir, 0700)
}

// Get retrieves cached data if it exists and is not expired
func (c *Cache) Get(key string, target interface{}) (bool, error) {
	cachePath := filepath.Join(c.GetCacheDir(), key+".json")

	data, err := os.ReadFile(cachePath)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, fmt.Errorf("failed to read cache file: %w", err)
	}

	var cached CachedData
	if err := json.Unmarshal(data, &cached); err != nil {
		return false, fmt.Errorf("failed to unmarshal cache data: %w", err)
	}

	// Check if expired
	if time.Since(cached.FetchedAt) > c.ttl {
		return false, nil
	}

	// Unmarshal the actual data into target
	dataBytes, err := json.Marshal(cached.Data)
	if err != nil {
		return false, fmt.Errorf("failed to marshal cached data: %w", err)
	}

	if err := json.Unmarshal(dataBytes, target); err != nil {
		return false, fmt.Errorf("failed to unmarshal target data: %w", err)
	}

	return true, nil
}

// Set stores data in the cache
func (c *Cache) Set(key string, data interface{}) error {
	if err := c.ensureCacheDir(); err != nil {
		return err
	}

	cached := CachedData{
		FetchedAt: time.Now(),
		Data:      data,
	}

	jsonData, err := json.Marshal(cached)
	if err != nil {
		return fmt.Errorf("failed to marshal cache data: %w", err)
	}

	cachePath := filepath.Join(c.GetCacheDir(), key+".json")
	if err := os.WriteFile(cachePath, jsonData, 0600); err != nil {
		return fmt.Errorf("failed to write cache file: %w", err)
	}

	return nil
}

// Invalidate removes a cache entry
func (c *Cache) Invalidate(key string) error {
	cachePath := filepath.Join(c.GetCacheDir(), key+".json")
	if err := os.Remove(cachePath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to remove cache file: %w", err)
	}
	return nil
}

// Clear removes all cache entries for this profile
func (c *Cache) Clear() error {
	dir := c.GetCacheDir()
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return nil
	}
	return os.RemoveAll(dir)
}
