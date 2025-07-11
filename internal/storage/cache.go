package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type CacheEntry struct {
	Key        string      `json:"key"`
	Data       interface{} `json:"data"`
	ExpiresAt  time.Time   `json:"expiresAt"`
	CreatedAt  time.Time   `json:"createdAt"`
	AccessedAt time.Time   `json:"accessedAt"`
}

type Cache struct {
	directory string
	enabled   bool
}

var cache *Cache

func GetCache() *Cache {
	if cache == nil {
		cfg := GetConfig()
		cache = &Cache{
			directory: cfg.CacheDirectory,
			enabled:   cfg.CacheEnabled,
		}
		
		// Ensure cache directory exists
		if err := CreateDirectory(cache.directory); err != nil {
			fmt.Printf("Warning: Failed to create cache directory: %v\n", err)
			cache.enabled = false
		}
	}
	return cache
}

func (c *Cache) Set(key string, data interface{}, ttl time.Duration) error {
	if !c.enabled {
		return nil
	}

	entry := CacheEntry{
		Key:        key,
		Data:       data,
		ExpiresAt:  time.Now().Add(ttl),
		CreatedAt:  time.Now(),
		AccessedAt: time.Now(),
	}

	return c.writeEntry(key, entry)
}

func (c *Cache) Get(key string, result interface{}) (bool, error) {
	if !c.enabled {
		return false, nil
	}

	entry, exists, err := c.readEntry(key)
	if err != nil {
		return false, err
	}

	if !exists {
		return false, nil
	}

	// Check if expired
	if time.Now().After(entry.ExpiresAt) {
		c.Delete(key) // Clean up expired entry
		return false, nil
	}

	// Update access time
	entry.AccessedAt = time.Now()
	c.writeEntry(key, entry)

	// Marshal and unmarshal to convert interface{} to desired type
	dataBytes, err := json.Marshal(entry.Data)
	if err != nil {
		return false, fmt.Errorf("failed to marshal cached data: %w", err)
	}

	if err := json.Unmarshal(dataBytes, result); err != nil {
		return false, fmt.Errorf("failed to unmarshal cached data: %w", err)
	}

	return true, nil
}

func (c *Cache) Delete(key string) error {
	if !c.enabled {
		return nil
	}

	filePath := c.getFilePath(key)
	if err := os.Remove(filePath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete cache entry: %w", err)
	}

	return nil
}

func (c *Cache) Clear() error {
	if !c.enabled {
		return nil
	}

	entries, err := os.ReadDir(c.directory)
	if err != nil {
		return fmt.Errorf("failed to read cache directory: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		filePath := filepath.Join(c.directory, entry.Name())
		if err := os.Remove(filePath); err != nil {
			fmt.Printf("Warning: Failed to remove cache file %s: %v\n", filePath, err)
		}
	}

	return nil
}

func (c *Cache) CleanExpired() error {
	if !c.enabled {
		return nil
	}

	entries, err := os.ReadDir(c.directory)
	if err != nil {
		return fmt.Errorf("failed to read cache directory: %w", err)
	}

	now := time.Now()
	cleaned := 0

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		filePath := filepath.Join(c.directory, entry.Name())
		cacheEntry, exists, err := c.readEntryFromFile(filePath)
		if err != nil {
			fmt.Printf("Warning: Failed to read cache entry %s: %v\n", filePath, err)
			continue
		}

		if !exists || now.After(cacheEntry.ExpiresAt) {
			if err := os.Remove(filePath); err != nil {
				fmt.Printf("Warning: Failed to remove expired cache file %s: %v\n", filePath, err)
			} else {
				cleaned++
			}
		}
	}

	if cleaned > 0 {
		fmt.Printf("Cleaned %d expired cache entries\n", cleaned)
	}

	return nil
}

func (c *Cache) Stats() (map[string]interface{}, error) {
	if !c.enabled {
		return map[string]interface{}{
			"enabled": false,
		}, nil
	}

	entries, err := os.ReadDir(c.directory)
	if err != nil {
		return nil, fmt.Errorf("failed to read cache directory: %w", err)
	}

	totalEntries := 0
	expiredEntries := 0
	totalSize := int64(0)
	now := time.Now()

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		info, err := entry.Info()
		if err != nil {
			continue
		}

		totalSize += info.Size()
		totalEntries++

		filePath := filepath.Join(c.directory, entry.Name())
		cacheEntry, exists, err := c.readEntryFromFile(filePath)
		if err != nil {
			continue
		}

		if exists && now.After(cacheEntry.ExpiresAt) {
			expiredEntries++
		}
	}

	return map[string]interface{}{
		"enabled":        true,
		"directory":      c.directory,
		"totalEntries":   totalEntries,
		"expiredEntries": expiredEntries,
		"totalSize":      totalSize,
		"totalSizeMB":    float64(totalSize) / (1024 * 1024),
	}, nil
}

func (c *Cache) readEntry(key string) (CacheEntry, bool, error) {
	filePath := c.getFilePath(key)
	return c.readEntryFromFile(filePath)
}

func (c *Cache) readEntryFromFile(filePath string) (CacheEntry, bool, error) {
	var entry CacheEntry

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return entry, false, nil
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return entry, false, fmt.Errorf("failed to read cache file: %w", err)
	}

	if err := json.Unmarshal(data, &entry); err != nil {
		return entry, false, fmt.Errorf("failed to unmarshal cache entry: %w", err)
	}

	return entry, true, nil
}

func (c *Cache) writeEntry(key string, entry CacheEntry) error {
	filePath := c.getFilePath(key)

	data, err := json.MarshalIndent(entry, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal cache entry: %w", err)
	}

	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write cache file: %w", err)
	}

	return nil
}

func (c *Cache) getFilePath(key string) string {
	// Create a safe filename from the key
	filename := fmt.Sprintf("%s.json", sanitizeFilename(key))
	return filepath.Join(c.directory, filename)
}

func sanitizeFilename(filename string) string {
	// Replace invalid characters with underscores
	invalid := []string{"/", "\\", ":", "*", "?", "\"", "<", ">", "|", " "}
	result := filename
	
	for _, char := range invalid {
		result = filepath.Base(result) // Remove path separators
		if len(result) > 100 {
			result = result[:100] // Limit length
		}
	}
	
	return result
}

// Convenience functions for common cache operations
func CacheSet(key string, data interface{}, ttl time.Duration) error {
	return GetCache().Set(key, data, ttl)
}

func CacheGet(key string, result interface{}) (bool, error) {
	return GetCache().Get(key, result)
}

func CacheDelete(key string) error {
	return GetCache().Delete(key)
}

func CacheClear() error {
	return GetCache().Clear()
}

func CacheCleanExpired() error {
	return GetCache().CleanExpired()
}