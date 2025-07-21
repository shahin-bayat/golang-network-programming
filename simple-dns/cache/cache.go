package cache

import (
	"sync"
	"time"

	"golang.org/x/net/dns/dnsmessage"
)

type CacheKey struct {
	Name dnsmessage.Name
	Type dnsmessage.Type
}

type CacheEntry struct {
	Resources []dnsmessage.Resource
	ExpiresAt time.Time
}

type DNSCache struct {
	mu      sync.RWMutex
	entries map[CacheKey]CacheEntry
}

func New() *DNSCache {
	return &DNSCache{
		entries: make(map[CacheKey]CacheEntry),
	}
}

func (c *DNSCache) Get(key CacheKey) (CacheEntry, bool) {
	c.mu.RLock()
	e, ok := c.entries[key]
	c.mu.RUnlock()

	if !ok {
		return e, false
	}

	if time.Now().After(e.ExpiresAt) {
		c.mu.Lock()
		defer c.mu.Unlock()
		delete(c.entries, key)
		return e, false
	}

	return e, true
}

func (c *DNSCache) Set(key CacheKey, entry CacheEntry) {
	c.mu.Lock()
	c.entries[key] = entry
	c.mu.Unlock()
}
