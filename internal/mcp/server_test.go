package mcp

import (
	"bytes"
	"context"
	"encoding/json"
	"strings"
	"testing"
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

	// Check response structure
	result, ok := resp.Result.(map[string]interface{})
	if !ok {
		t.Fatalf("Response result is not a map: %T", resp.Result)
	}

	// Check server info
	serverInfo, ok := result["serverInfo"].(map[string]interface{})
	if !ok {
		t.Fatal("serverInfo is not a map")
	}
	if serverInfo["name"] != "devwisdom-go" {
		t.Errorf("serverInfo name = %q, want %q", serverInfo["name"], "devwisdom-go")
	}
	if serverInfo["version"] != Version {
		t.Errorf("serverInfo version = %q, want %q", serverInfo["version"], Version)
	}
}

func TestWisdomServer_HandleGetWisdom(t *testing.T) {
	server := NewWisdomServer()
	server.wisdom.Initialize()

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

	// Check response structure
	result, ok := resp.Result.(map[string]interface{})
	if !ok {
		t.Fatalf("Response result is not a map: %T", resp.Result)
	}

	if result["quote"] == nil {
		t.Error("Response missing quote field")
	}
	if result["source"] == nil {
		t.Error("Response missing source field")
	}
}

func TestWisdomServer_HandleConsultAdvisor(t *testing.T) {
	server := NewWisdomServer()
	server.wisdom.Initialize()

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

	// Check response structure
	result, ok := resp.Result.(map[string]interface{})
	if !ok {
		t.Fatalf("Response result is not a map: %T", resp.Result)
	}

	if result["advisor"] == nil {
		t.Error("Response missing advisor field")
	}
	if result["quote"] == nil {
		t.Error("Response missing quote field")
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
	server.wisdom.Initialize()

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
	server.wisdom.Initialize()

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

	// Check response structure
	result, ok := resp.Result.(map[string]interface{})
	if !ok {
		t.Fatalf("Response result is not a map: %T", resp.Result)
	}

	contents, ok := result["contents"].([]interface{})
	if !ok {
		t.Fatal("Response contents is not an array")
	}
	if len(contents) == 0 {
		t.Error("Response contents is empty")
	}
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

	// Notification (no ID) should not get a response
	req := &JSONRPCRequest{
		JSONRPC: "2.0",
		ID:      nil, // Notification
		Method:  "tools/call",
		Params: json.RawMessage(`{
			"name": "get_wisdom",
			"arguments": {"score": 75.0, "source": "stoic"}
		}`),
	}

	resp := server.handleRequest(req)
	if resp != nil {
		t.Error("Notifications should not receive responses")
	}
}

func TestWisdomServer_HandleGetDailyBriefing(t *testing.T) {
	server := NewWisdomServer()
	server.wisdom.Initialize()

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

	result, ok := resp.Result.(map[string]interface{})
	if !ok {
		t.Fatalf("Response result is not a map: %T", resp.Result)
	}

	if result["quote"] == nil {
		t.Error("Response missing quote field")
	}
	if result["source"] == nil {
		t.Error("Response missing source field")
	}
}

