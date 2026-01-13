package sefaria

import (
	"testing"
	"time"
)

func TestCache_Invalidate(t *testing.T) {
	cache := NewCache()
	response := &TextResponse{Ref: "Test", Text: []string{"test"}}

	// Set an entry
	cache.Set("key1", response)

	// Verify it exists
	if _, found := cache.Get("key1"); !found {
		t.Fatal("Entry should exist before invalidation")
	}

	// Invalidate it
	cache.Invalidate("key1")

	// Verify it's gone
	if _, found := cache.Get("key1"); found {
		t.Error("Entry should not exist after invalidation")
	}

	// Invalidate non-existent key (should not panic)
	cache.Invalidate("nonexistent")
}

func TestCache_InvalidateAll(t *testing.T) {
	cache := NewCache()
	response := &TextResponse{Ref: "Test", Text: []string{"test"}}

	// Set multiple entries
	cache.Set("key1", response)
	cache.Set("key2", response)
	cache.Set("key3", response)

	// Verify they exist
	if _, found := cache.Get("key1"); !found {
		t.Error("key1 should exist")
	}
	if _, found := cache.Get("key2"); !found {
		t.Error("key2 should exist")
	}
	if _, found := cache.Get("key3"); !found {
		t.Error("key3 should exist")
	}

	// Invalidate all
	cache.InvalidateAll()

	// Verify all are gone
	if _, found := cache.Get("key1"); found {
		t.Error("key1 should not exist after InvalidateAll")
	}
	if _, found := cache.Get("key2"); found {
		t.Error("key2 should not exist after InvalidateAll")
	}
	if _, found := cache.Get("key3"); found {
		t.Error("key3 should not exist after InvalidateAll")
	}
}

func TestCache_Cleanup(t *testing.T) {
	cache := NewCache()
	response := &TextResponse{Ref: "Test", Text: []string{"test"}}

	// Set default TTL to 1 hour for testing
	cache.SetTTL(1 * time.Hour)

	// Add a fresh entry
	cache.Set("fresh", response)

	// Add an expired entry (manually set timestamp)
	cache.mu.Lock()
	cache.entries["expired"] = &CacheEntry{
		Response:  response,
		Timestamp: time.Now().Add(-2 * time.Hour), // Expired
		TTL:       1 * time.Hour,
	}
	cache.mu.Unlock()

	// Verify both exist before cleanup
	if _, found := cache.Get("fresh"); !found {
		t.Error("fresh entry should exist")
	}
	if _, found := cache.Get("expired"); found {
		t.Error("expired entry should not be returned by Get()")
	}

	// Cleanup should remove expired entry
	cache.Cleanup()

	// Verify fresh entry still exists
	if _, found := cache.Get("fresh"); !found {
		t.Error("fresh entry should still exist after cleanup")
	}

	// Verify expired entry is removed
	cache.mu.RLock()
	_, exists := cache.entries["expired"]
	cache.mu.RUnlock()
	if exists {
		t.Error("expired entry should be removed by Cleanup")
	}
}

func TestCache_SetTTL(t *testing.T) {
	cache := NewCache()
	originalTTL := cache.ttl

	// Set new TTL
	newTTL := 2 * time.Hour
	cache.SetTTL(newTTL)

	if cache.ttl != newTTL {
		t.Errorf("TTL = %v, want %v", cache.ttl, newTTL)
	}

	// Verify new entries use new TTL
	response := &TextResponse{Ref: "Test", Text: []string{"test"}}
	cache.Set("test", response)

	cache.mu.RLock()
	entry := cache.entries["test"]
	cache.mu.RUnlock()

	if entry.TTL != newTTL {
		t.Errorf("Entry TTL = %v, want %v", entry.TTL, newTTL)
	}

	// Restore original TTL
	cache.SetTTL(originalTTL)
}

func TestCache_Expiration(t *testing.T) {
	cache := NewCache()
	// Set very short TTL for testing
	cache.SetTTL(100 * time.Millisecond)

	response := &TextResponse{Ref: "Test", Text: []string{"test"}}
	cache.Set("test", response)

	// Should exist immediately
	if _, found := cache.Get("test"); !found {
		t.Error("Entry should exist immediately after setting")
	}

	// Wait for expiration
	time.Sleep(150 * time.Millisecond)

	// Should not exist after expiration
	if _, found := cache.Get("test"); found {
		t.Error("Entry should not exist after expiration")
	}
}

func TestCache_ConcurrentAccess(t *testing.T) {
	cache := NewCache()
	response := &TextResponse{Ref: "Test", Text: []string{"test"}}

	// Test concurrent Set operations
	done := make(chan bool, 10)
	for i := 0; i < 10; i++ {
		go func(key string) {
			cache.Set(key, response)
			done <- true
		}("key" + string(rune(i)))
	}

	// Wait for all goroutines
	for i := 0; i < 10; i++ {
		<-done
	}

	// Verify all entries exist
	for i := 0; i < 10; i++ {
		key := "key" + string(rune(i))
		if _, found := cache.Get(key); !found {
			t.Errorf("Entry %q should exist", key)
		}
	}
}
