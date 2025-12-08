package wisdom

import (
	"os"
	"testing"
	"time"
)

func TestNewSourceCache(t *testing.T) {
	cache := NewSourceCache()
	if cache == nil {
		t.Fatal("NewSourceCache returned nil")
	}
	if cache.Size() != 0 {
		t.Errorf("New cache should be empty, got size %d", cache.Size())
	}
}

func TestSourceCache_SetAndGet(t *testing.T) {
	cache := NewSourceCache()
	config := &SourceConfig{
		ID:   "test",
		Name: "Test Source",
		Icon: "📜",
		Quotes: map[string][]Quote{
			"chaos": {{Quote: "Test", Source: "Test", Encouragement: "Test"}},
		},
	}

	cache.Set("test_key", config, "")

	retrieved, found := cache.Get("test_key")
	if !found {
		t.Fatal("Get returned false for just-set key")
	}
	if retrieved == nil {
		t.Fatal("Get returned nil config")
	}
	if retrieved.ID != "test" {
		t.Errorf("Get returned ID %q, want %q", retrieved.ID, "test")
	}
}

func TestSourceCache_Expiration(t *testing.T) {
	cache := NewSourceCache().WithTTL(100 * time.Millisecond)
	config := &SourceConfig{
		ID:   "test",
		Name: "Test",
		Icon: "📜",
		Quotes: map[string][]Quote{
			"chaos": {{Quote: "Test", Source: "Test", Encouragement: "Test"}},
		},
	}

	cache.Set("test_key", config, "")

	// Should be found immediately
	_, found := cache.Get("test_key")
	if !found {
		t.Fatal("Get returned false immediately after Set")
	}

	// Wait for expiration
	time.Sleep(150 * time.Millisecond)

	// Should not be found after expiration
	_, found = cache.Get("test_key")
	if found {
		t.Error("Get returned true for expired entry")
	}
}

func TestSourceCache_Invalidate(t *testing.T) {
	cache := NewSourceCache()
	config := &SourceConfig{
		ID:   "test",
		Name: "Test",
		Icon: "📜",
		Quotes: map[string][]Quote{
			"chaos": {{Quote: "Test", Source: "Test", Encouragement: "Test"}},
		},
	}

	cache.Set("test_key", config, "")
	cache.Invalidate("test_key")

	_, found := cache.Get("test_key")
	if found {
		t.Error("Get returned true for invalidated entry")
	}
}

func TestSourceCache_InvalidateAll(t *testing.T) {
	cache := NewSourceCache()
	config := &SourceConfig{
		ID:   "test",
		Name: "Test",
		Icon: "📜",
		Quotes: map[string][]Quote{
			"chaos": {{Quote: "Test", Source: "Test", Encouragement: "Test"}},
		},
	}

	cache.Set("key1", config, "")
	cache.Set("key2", config, "")
	cache.Set("key3", config, "")

	if cache.Size() != 3 {
		t.Errorf("Cache size = %d, want 3", cache.Size())
	}

	cache.InvalidateAll()

	if cache.Size() != 0 {
		t.Errorf("Cache size after InvalidateAll = %d, want 0", cache.Size())
	}
}

func TestSourceCache_FileModificationTracking(t *testing.T) {
	// Create a temporary file
	tmpFile, err := os.CreateTemp("", "test_source_*.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	tmpFile.Close()

	cache := NewSourceCache()
	config := &SourceConfig{
		ID:   "test",
		Name: "Test",
		Icon: "📜",
		Quotes: map[string][]Quote{
			"chaos": {{Quote: "Test", Source: "Test", Encouragement: "Test"}},
		},
	}

	cache.Set("test_key", config, tmpFile.Name())

	// Should be found
	_, found := cache.Get("test_key")
	if !found {
		t.Fatal("Get returned false for cached entry")
	}

	// Modify the file
	time.Sleep(10 * time.Millisecond) // Ensure different mod time
	if err := os.WriteFile(tmpFile.Name(), []byte("modified"), 0644); err != nil {
		t.Fatalf("Failed to modify file: %v", err)
	}

	// Should not be found after file modification
	_, found = cache.Get("test_key")
	if found {
		t.Error("Get returned true for entry with modified file")
	}
}

func TestSourceCache_ClearExpired(t *testing.T) {
	cache := NewSourceCache().WithTTL(50 * time.Millisecond)
	config := &SourceConfig{
		ID:   "test",
		Name: "Test",
		Icon: "📜",
		Quotes: map[string][]Quote{
			"chaos": {{Quote: "Test", Source: "Test", Encouragement: "Test"}},
		},
	}

	cache.Set("key1", config, "")
	cache.Set("key2", config, "")

	// Wait for expiration
	time.Sleep(100 * time.Millisecond)

	cleared := cache.ClearExpired()
	if cleared != 2 {
		t.Errorf("ClearExpired cleared %d entries, want 2", cleared)
	}
	if cache.Size() != 0 {
		t.Errorf("Cache size after ClearExpired = %d, want 0", cache.Size())
	}
}

func TestSourceCache_Disable(t *testing.T) {
	cache := NewSourceCache().Enable(false)
	config := &SourceConfig{
		ID:   "test",
		Name: "Test",
		Icon: "📜",
		Quotes: map[string][]Quote{
			"chaos": {{Quote: "Test", Source: "Test", Encouragement: "Test"}},
		},
	}

	cache.Set("test_key", config, "")

	_, found := cache.Get("test_key")
	if found {
		t.Error("Get returned true when cache is disabled")
	}
}
