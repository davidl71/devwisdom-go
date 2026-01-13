package sefaria

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestClient_GetText(t *testing.T) {
	// Create mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/texts/Proverbs" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{
				"ref": "Proverbs 1",
				"heRef": "משלי א׳",
				"text": ["The proverbs of Solomon", "To know wisdom"],
				"he": ["מִשְׁלֵי שְׁלֹמֹה", "לָדַעַת חָכְמָה"]
			}`))
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	// Create client with test server URL
	client := &Client{
		httpClient: &http.Client{Timeout: 5 * time.Second},
		baseURL:    server.URL + "/api",
		cache:      NewCache(),
	}

	ctx := context.Background()
	resp, err := client.GetText(ctx, "Proverbs", 0, 0)
	if err != nil {
		t.Fatalf("GetText failed: %v", err)
	}

	if resp.Ref != "Proverbs 1" {
		t.Errorf("Expected ref 'Proverbs 1', got %q", resp.Ref)
	}

	if len(resp.Text) != 2 {
		t.Errorf("Expected 2 English verses, got %d", len(resp.Text))
	}

	if len(resp.He) != 2 {
		t.Errorf("Expected 2 Hebrew verses, got %d", len(resp.He))
	}
}

func TestClient_GetTextBySourceID(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/texts/Pirkei_Avot" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{
				"ref": "Pirkei Avot 1",
				"heRef": "פרקי אבות א׳",
				"text": ["Moses received the Torah"],
				"he": ["מֹשֶׁה קִבֵּל תּוֹרָה"]
			}`))
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	client := &Client{
		httpClient: &http.Client{Timeout: 5 * time.Second},
		baseURL:    server.URL + "/api",
		cache:      NewCache(),
	}

	ctx := context.Background()
	resp, err := client.GetTextBySourceID(ctx, "pirkei_avot", 0, 0)
	if err != nil {
		t.Fatalf("GetTextBySourceID failed: %v", err)
	}

	if resp.Ref != "Pirkei Avot 1" {
		t.Errorf("Expected ref 'Pirkei Avot 1', got %q", resp.Ref)
	}
}

func TestClient_GetTextBySourceID_UnknownSource(t *testing.T) {
	client := NewClient(nil)
	ctx := context.Background()

	_, err := client.GetTextBySourceID(ctx, "unknown_source", 0, 0)
	if err == nil {
		t.Error("Expected error for unknown source ID")
	}
}

func TestClient_Cache(t *testing.T) {
	callCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"ref": "Proverbs 1",
			"text": ["Test"],
			"he": ["בדיקה"]
		}`))
	}))
	defer server.Close()

	client := &Client{
		httpClient: &http.Client{Timeout: 5 * time.Second},
		baseURL:    server.URL + "/api",
		cache:      NewCache(),
	}

	ctx := context.Background()

	// First call - should hit API
	_, err := client.GetText(ctx, "Proverbs", 0, 0)
	if err != nil {
		t.Fatalf("First GetText failed: %v", err)
	}

	// Second call - should use cache
	_, err = client.GetText(ctx, "Proverbs", 0, 0)
	if err != nil {
		t.Fatalf("Second GetText failed: %v", err)
	}

	// Should only have called API once
	if callCount != 1 {
		t.Errorf("Expected 1 API call, got %d", callCount)
	}
}

func TestClient_CleanupCache(t *testing.T) {
	client := NewClient(nil)
	cache := NewCache()

	// Add an expired entry
	expiredEntry := &CacheEntry{
		Response:  &TextResponse{Ref: "Test"},
		Timestamp: time.Now().Add(-25 * time.Hour), // Expired (TTL is 24 hours)
		TTL:       24 * time.Hour,
	}
	cache.mu.Lock()
	cache.entries["test:key"] = expiredEntry
	cache.mu.Unlock()

	client.cache = cache

	// Cleanup should remove expired entry
	client.CleanupCache()

	cache.mu.RLock()
	_, exists := cache.entries["test:key"]
	cache.mu.RUnlock()

	if exists {
		t.Error("CleanupCache should have removed expired entry")
	}
}

func TestClient_buildEndpoint(t *testing.T) {
	client := NewClient(nil)

	tests := []struct {
		name     string
		book     string
		chapter  int
		verse    int
		expected string
	}{
		{"full book", "Proverbs", 0, 0, "texts/Proverbs"},
		{"full chapter", "Proverbs", 1, 0, "texts/Proverbs.1"},
		{"specific verse", "Proverbs", 1, 5, "texts/Proverbs.1.5"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := client.buildEndpoint(tt.book, tt.chapter, tt.verse)
			if result != tt.expected {
				t.Errorf("buildEndpoint(%q, %d, %d) = %q, want %q", tt.book, tt.chapter, tt.verse, result, tt.expected)
			}
		})
	}
}

func TestClient_GetText_ErrorHandling(t *testing.T) {
	tests := []struct {
		name           string
		statusCode     int
		responseBody   string
		expectedErrMsg string
	}{
		{
			name:           "404 not found",
			statusCode:     http.StatusNotFound,
			responseBody:   `{"error": "Not found"}`,
			expectedErrMsg: "Sefaria API returned status 404",
		},
		{
			name:           "500 server error",
			statusCode:     http.StatusInternalServerError,
			responseBody:   `{"error": "Internal error"}`,
			expectedErrMsg: "Sefaria API returned status 500",
		},
		{
			name:           "invalid JSON",
			statusCode:     http.StatusOK,
			responseBody:   `invalid json`,
			expectedErrMsg: "failed to parse Sefaria API response",
		},
		{
			name:           "empty response",
			statusCode:     http.StatusOK,
			responseBody:   `{"ref": "", "text": [], "he": []}`,
			expectedErrMsg: "Sefaria API response has no text content",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.statusCode)
				_, _ = w.Write([]byte(tt.responseBody))
			}))
			defer server.Close()

			client := &Client{
				httpClient: &http.Client{Timeout: 5 * time.Second},
				baseURL:    server.URL + "/api",
				cache:      NewCache(),
			}

			ctx := context.Background()
			_, err := client.GetText(ctx, "Proverbs", 0, 0)

			if err == nil {
				t.Error("Expected error but got nil")
			} else if !strings.Contains(err.Error(), tt.expectedErrMsg) {
				t.Errorf("Error message %q does not contain %q", err.Error(), tt.expectedErrMsg)
			}
		})
	}
}
