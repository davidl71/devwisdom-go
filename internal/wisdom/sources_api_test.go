package wisdom

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNewAPISourceLoader(t *testing.T) {
	loader := NewAPISourceLoader("https://api.example.com", 10*time.Second)
	if loader == nil {
		t.Fatal("NewAPISourceLoader returned nil")
	}
	if loader.baseURL != "https://api.example.com" {
		t.Errorf("baseURL = %q, want %q", loader.baseURL, "https://api.example.com")
	}
	if loader.timeout != 10*time.Second {
		t.Errorf("timeout = %v, want %v", loader.timeout, 10*time.Second)
	}
}

func TestAPISourceLoader_LoadSource_Success(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"id": "test",
			"name": "Test Source",
			"icon": "📜",
			"quotes": {
				"chaos": [{
					"quote": "Test quote",
					"source": "Test",
					"encouragement": "Test"
				}]
			}
		}`))
	}))
	defer server.Close()

	loader := NewAPISourceLoader(server.URL, 5*time.Second)
	ctx := context.Background()

	config, err := loader.LoadSource(ctx, "test")
	if err != nil {
		t.Fatalf("LoadSource failed: %v", err)
	}
	if config == nil {
		t.Fatal("LoadSource returned nil config")
	}
	if config.ID != "test" {
		t.Errorf("Config ID = %q, want %q", config.ID, "test")
	}
	if config.Name != "Test Source" {
		t.Errorf("Config Name = %q, want %q", config.Name, "Test Source")
	}
}

func TestAPISourceLoader_LoadSource_NotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	loader := NewAPISourceLoader(server.URL, 5*time.Second)
	ctx := context.Background()

	_, err := loader.LoadSource(ctx, "test")
	if err == nil {
		t.Error("LoadSource should fail for 404 response")
	}
}

func TestAPISourceLoader_LoadSource_Timeout(t *testing.T) {
	// Create a slow server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(2 * time.Second) // Longer than timeout
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	loader := NewAPISourceLoader(server.URL, 100*time.Millisecond)
	ctx := context.Background()

	_, err := loader.LoadSource(ctx, "test")
	if err == nil {
		t.Error("LoadSource should fail on timeout")
	}
}

func TestAPISourceLoader_LoadSourceWithRetry(t *testing.T) {
	attempts := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		if attempts < 2 {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"id": "test",
			"name": "Test",
			"icon": "📜",
			"quotes": {
				"chaos": [{"quote": "Test", "source": "Test", "encouragement": "Test"}]
			}
		}`))
	}))
	defer server.Close()

	loader := NewAPISourceLoader(server.URL, 5*time.Second)
	ctx := context.Background()

	config, err := loader.LoadSourceWithRetry(ctx, "test", 3)
	if err != nil {
		t.Fatalf("LoadSourceWithRetry failed: %v", err)
	}
	if config == nil {
		t.Fatal("LoadSourceWithRetry returned nil config")
	}
	if attempts < 2 {
		t.Errorf("Expected at least 2 attempts, got %d", attempts)
	}
}

func TestAPISourceLoader_LoadSourceWithTimeout(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"id": "test",
			"name": "Test",
			"icon": "📜",
			"quotes": {
				"chaos": [{"quote": "Test", "source": "Test", "encouragement": "Test"}]
			}
		}`))
	}))
	defer server.Close()

	loader := NewAPISourceLoader(server.URL, 10*time.Second)

	config, err := loader.LoadSourceWithTimeout("test", 5*time.Second)
	if err != nil {
		t.Fatalf("LoadSourceWithTimeout failed: %v", err)
	}
	if config == nil {
		t.Fatal("LoadSourceWithTimeout returned nil config")
	}
}
