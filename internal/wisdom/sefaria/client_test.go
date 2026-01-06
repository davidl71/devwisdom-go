package sefaria

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestClient_GetText(t *testing.T) {
	// Create mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/texts/Proverbs" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{
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
			w.Write([]byte(`{
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
		w.Write([]byte(`{
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

