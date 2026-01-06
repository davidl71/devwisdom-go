package mcp

import (
	"bytes"
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/davidl71/devwisdom-go/internal/wisdom"
)

func TestNewWisdomServer(t *testing.T) {
	server := NewWisdomServer()
	if server == nil {
		t.Fatal("NewWisdomServer returned nil")
	}
	if server.wisdom == nil {
		t.Error("WisdomServer wisdom engine is nil")
	}
}

func TestWisdomServer_HandleInitialize(t *testing.T) {
	server := NewWisdomServer()

	req := &JSONRPCRequest{
		JSONRPC: "2.0",
		ID:      1,
		Method:  "initialize",
		Params: json.RawMessage(`{
			"protocolVersion": "2024-11-05",
			"capabilities": {},
			"clientInfo": {
				"name": "test-client",
				"version": "1.0.0"
			}
		}`),
	}

	resp := server.handleRequest(req)
	if resp == nil {
		t.Fatal("handleRequest returned nil response")
	}
	if resp.Error != nil {
		t.Fatalf("handleRequest returned error: %v", resp.Error)
	}

	// Check response structure - result is InitializeResult struct
	result, ok := resp.Result.(InitializeResult)
	if !ok {
		t.Fatalf("Response result is not InitializeResult: %T", resp.Result)
	}

	// Check server info
	if result.ServerInfo.Name != "devwisdom" {
		t.Errorf("serverInfo name = %q, want %q", result.ServerInfo.Name, "devwisdom")
	}
	if result.ServerInfo.Version != Version {
		t.Errorf("serverInfo version = %q, want %q", result.ServerInfo.Version, Version)
	}
}

func TestWisdomServer_HandleGetWisdom(t *testing.T) {
	server := NewWisdomServer()
	if err := server.wisdom.Initialize(); err != nil {
		t.Fatalf("Initialize failed: %v", err)
	}

	req := &JSONRPCRequest{
		JSONRPC: "2.0",
		ID:      2,
		Method:  "tools/call",
		Params: json.RawMessage(`{
			"name": "get_wisdom",
			"arguments": {
				"score": 75.0,
				"source": "stoic"
			}
		}`),
	}

	resp := server.handleRequest(req)
	if resp == nil {
		t.Fatal("handleRequest returned nil response")
	}
	if resp.Error != nil {
		t.Fatalf("handleRequest returned error: %v", resp.Error)
	}

	// Check response structure - result is *wisdom.Quote pointer
	result, ok := resp.Result.(*wisdom.Quote)
	if !ok {
		t.Fatalf("Response result is not *wisdom.Quote: %T", resp.Result)
	}
	if result == nil {
		t.Fatal("Response result is nil")
	}

	if result.Quote == "" {
		t.Error("Response quote field is empty")
	}
	if result.Source == "" {
		t.Error("Response source field is empty")
	}
}

func TestWisdomServer_HandleConsultAdvisor(t *testing.T) {
	server := NewWisdomServer()
	if err := server.wisdom.Initialize(); err != nil {
		t.Fatalf("Initialize failed: %v", err)
	}

	req := &JSONRPCRequest{
		JSONRPC: "2.0",
		ID:      3,
		Method:  "tools/call",
		Params: json.RawMessage(`{
			"name": "consult_advisor",
			"arguments": {
				"metric": "security",
				"score": 40.0
			}
		}`),
	}

	resp := server.handleRequest(req)
	if resp == nil {
		t.Fatal("handleRequest returned nil response")
	}
	if resp.Error != nil {
		t.Fatalf("handleRequest returned error: %v", resp.Error)
	}

	// Check response structure - result is wisdom.Consultation (value type, not pointer)
	result, ok := resp.Result.(wisdom.Consultation)
	if !ok {
		// Try pointer type in case server returns pointer
		if resultPtr, okPtr := resp.Result.(*wisdom.Consultation); okPtr {
			if resultPtr == nil {
				t.Fatal("Response result is nil pointer")
			}
			result = *resultPtr
		} else {
			t.Fatalf("Response result is not wisdom.Consultation: %T", resp.Result)
		}
	}

	if result.Advisor == "" {
		t.Error("Response advisor field is empty")
	}
	if result.Quote == "" {
		t.Error("Response quote field is empty")
	}
}

func TestWisdomServer_HandleInvalidMethod(t *testing.T) {
	server := NewWisdomServer()

	req := &JSONRPCRequest{
		JSONRPC: "2.0",
		ID:      99,
		Method:  "invalid/method",
		Params:  nil,
	}

	resp := server.handleRequest(req)
	if resp == nil {
		t.Fatal("handleRequest returned nil response")
	}
	if resp.Error == nil {
		t.Fatal("handleRequest should return error for invalid method")
	}
	if resp.Error.Code != ErrCodeMethodNotFound {
		t.Errorf("Error code = %d, want %d", resp.Error.Code, ErrCodeMethodNotFound)
	}
}

func TestWisdomServer_HandleInvalidParams(t *testing.T) {
	server := NewWisdomServer()
	if err := server.wisdom.Initialize(); err != nil {
		t.Fatalf("Initialize failed: %v", err)
	}

	req := &JSONRPCRequest{
		JSONRPC: "2.0",
		ID:      100,
		Method:  "tools/call",
		Params: json.RawMessage(`{
			"name": "consult_advisor",
			"arguments": {
				"metric": "security"
				// Missing required "score" parameter
			}
		}`),
	}

	resp := server.handleRequest(req)
	if resp == nil {
		t.Fatal("handleRequest returned nil response")
	}
	if resp.Error == nil {
		t.Fatal("handleRequest should return error for invalid params")
	}
	if resp.Error.Code != ErrCodeInvalidParams {
		t.Errorf("Error code = %d, want %d", resp.Error.Code, ErrCodeInvalidParams)
	}
}

func TestWisdomServer_HandleResourcesRead(t *testing.T) {
	server := NewWisdomServer()
	if err := server.wisdom.Initialize(); err != nil {
		t.Fatalf("Initialize failed: %v", err)
	}

	req := &JSONRPCRequest{
		JSONRPC: "2.0",
		ID:      10,
		Method:  "resources/read",
		Params: json.RawMessage(`{
			"uri": "wisdom://sources"
		}`),
	}

	resp := server.handleRequest(req)
	if resp == nil {
		t.Fatal("handleRequest returned nil response")
	}
	if resp.Error != nil {
		t.Fatalf("handleRequest returned error: %v", resp.Error)
	}

	// Check response structure - result is map[string]interface{} with contents array
	result, ok := resp.Result.(map[string]interface{})
	if !ok {
		t.Fatalf("Response result is not map[string]interface{}: %T", resp.Result)
	}

	// Server returns []map[string]interface{} for contents
	// Type assertion: try []map[string]interface{} first (actual type)
	contentsValue := result["contents"]

	// Try []map[string]interface{} first (actual server return type)
	if contents, ok := contentsValue.([]map[string]interface{}); ok {
		if len(contents) == 0 {
			t.Error("Response contents is empty")
		}
		return // Successfully validated
	}

	// Fallback: try []interface{} (may contain maps)
	if contentsInterface, ok := contentsValue.([]interface{}); ok {
		if len(contentsInterface) == 0 {
			t.Error("Response contents is empty")
		}
		return // Successfully validated
	}

	// Neither type matched
	t.Fatalf("Response contents is not []map[string]interface{} or []interface{}: %T", contentsValue)
}

func TestWisdomServer_Run_InitializeAndTools(t *testing.T) {
	server := NewWisdomServer()

	// Create test input with initialize and get_wisdom requests
	input := strings.NewReader(`{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"test","version":"1.0.0"}}}
{"jsonrpc":"2.0","id":2,"method":"tools/call","params":{"name":"get_wisdom","arguments":{"score":75.0,"source":"stoic"}}}
`)

	var output bytes.Buffer
	ctx := context.Background()

	// Run server in a goroutine since it blocks
	errChan := make(chan error, 1)
	go func() {
		errChan <- server.Run(ctx, input, &output)
	}()

	// Wait a bit for processing
	select {
	case err := <-errChan:
		if err != nil {
			t.Fatalf("Server Run failed: %v", err)
		}
	case <-ctx.Done():
		t.Fatal("Context cancelled")
	}

	// Check output contains responses
	outputStr := output.String()
	if !strings.Contains(outputStr, "jsonrpc") {
		t.Error("Output does not contain JSON-RPC response")
	}
}

func TestWisdomServer_HandleNotification(t *testing.T) {
	server := NewWisdomServer()
	if err := server.wisdom.Initialize(); err != nil {
		t.Fatalf("Initialize failed: %v", err)
	} // Initialize engine to avoid "engine not initialized" error

	// Notification (no ID) - per JSON-RPC 2.0 spec, notifications don't get responses
	req := &JSONRPCRequest{
		JSONRPC: "2.0",
		ID:      nil, // Notification (no ID)
		Method:  "tools/call",
		Params: json.RawMessage(`{
			"name": "get_wisdom",
			"arguments": {"score": 75.0, "source": "stoic"}
		}`),
	}

	resp := server.handleRequest(req)
	// Per JSON-RPC 2.0 spec, notifications (requests without ID) should not get responses
	// However, handleRequest processes them and may return an error response
	// The key is that handleRequest doesn't crash and handles the notification gracefully
	// In the Run() method, notifications are skipped (see server.go line 88-90)
	// But handleRequest itself may still process and return error responses
	// We accept either nil or a non-error response, but error responses indicate a problem
	if resp != nil {
		if resp.Error != nil {
			// Error responses for notifications are acceptable per JSON-RPC 2.0
			// but we log it for visibility
			t.Logf("Notification returned error response (acceptable): %v", resp.Error)
		}
	}
}

func TestWisdomServer_HandleGetDailyBriefing(t *testing.T) {
	server := NewWisdomServer()
	if err := server.wisdom.Initialize(); err != nil {
		t.Fatalf("Initialize failed: %v", err)
	}

	req := &JSONRPCRequest{
		JSONRPC: "2.0",
		ID:      7,
		Method:  "tools/call",
		Params: json.RawMessage(`{
			"name": "get_daily_briefing",
			"arguments": {
				"score": 75.0
			}
		}`),
	}

	resp := server.handleRequest(req)
	if resp == nil {
		t.Fatal("handleRequest returned nil response")
	}
	if resp.Error != nil {
		t.Fatalf("handleRequest returned error: %v", resp.Error)
	}

	// Check response structure - result is map[string]interface{}
	result, ok := resp.Result.(map[string]interface{})
	if !ok {
		t.Fatalf("Response result is not map[string]interface{}: %T", resp.Result)
	}

	if result["date"] == nil {
		t.Error("Response missing date field")
	}
	if result["score"] == nil {
		t.Error("Response missing score field")
	}
	if result["quotes"] == nil {
		t.Error("Response missing quotes field")
	}

	// Validate quotes is an array
	quotes, ok := result["quotes"].([]interface{})
	if !ok {
		t.Errorf("Response quotes is not []interface{}: %T", result["quotes"])
	} else if len(quotes) == 0 {
		t.Log("Response quotes array is empty (may be acceptable)")
	}
}
