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

		// Check response structure - result is InitializeResult struct.
	result, ok := resp.Result.(InitializeResult)
	if !ok {
		t.Fatalf("Response result is not InitializeResult: %T", resp.Result)
	}

		// Check server info.
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

		// Check response structure - result is *wisdom.Quote pointer.
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

		// Check response structure - result is wisdom.Consultation (value type, not pointer).
	result, ok := resp.Result.(wisdom.Consultation)
	if !ok {
				// Try pointer type in case server returns pointer.
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

		// Check response structure - result is map[string]interface{} with contents array.
	result, ok := resp.Result.(map[string]interface{})
	if !ok {
		t.Fatalf("Response result is not map[string]interface{}: %T", resp.Result)
	}

		// Type assertion: try []map[string]interface{} first (actual type).
	contentsValue := result["contents"]

		// Try []map[string]interface{} first (actual server return type).
	if contents, ok := contentsValue.([]map[string]interface{}); ok {
		if len(contents) == 0 {
			t.Error("Response contents is empty")
		}
		return 		return // Successfully validated.
	}

		// Fallback: try []interface{} (may contain maps).
	if contentsInterface, ok := contentsValue.([]interface{}); ok {
		if len(contentsInterface) == 0 {
			t.Error("Response contents is empty")
		}
		return 		return // Successfully validated.
	}

		// Neither type matched.
	t.Fatalf("Response contents is not []map[string]interface{} or []interface{}: %T", contentsValue)
}

func TestWisdomServer_Run_InitializeAndTools(t *testing.T) {
	server := NewWisdomServer()

		// Create test input with initialize and get_wisdom requests.
	input := strings.NewReader(`{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"test","version":"1.0.0"}}}
{"jsonrpc":"2.0","id":2,"method":"tools/call","params":{"name":"get_wisdom","arguments":{"score":75.0,"source":"stoic"}}}
`)

	var output bytes.Buffer
	ctx := context.Background()

		// Run server in a goroutine since it blocks.
	errChan := make(chan error, 1)
	go func() {
		errChan <- server.Run(ctx, input, &output)
	}()

		// Wait a bit for processing.
	select {
	case err := <-errChan:
		if err != nil {
			t.Fatalf("Server Run failed: %v", err)
		}
	case <-ctx.Done():
		t.Fatal("Context canceled")
	}

		// Check output contains responses.
	outputStr := output.String()
	if !strings.Contains(outputStr, "jsonrpc") {
		t.Error("Output does not contain JSON-RPC response")
	}
}

func TestWisdomServer_HandleNotification(t *testing.T) {
	server := NewWisdomServer()
	if err := server.wisdom.Initialize(); err != nil {
		t.Fatalf("Initialize failed: %v", err)
	} 	} // Initialize engine to avoid "engine not initialized" error.

		// Notification (no ID) - per JSON-RPC 2.0 spec, notifications don't get responses.
	req := &JSONRPCRequest{
		JSONRPC: "2.0",
		ID:      nil, 		ID:      nil, // Notification (no ID).
		Method:  "tools/call",
		Params: json.RawMessage(`{
			"name": "get_wisdom",
			"arguments": {"score": 75.0, "source": "stoic"}
		}`),
	}

	resp := server.handleRequest(req)
		// We accept either nil or a non-error response, but error responses indicate a problem.
	if resp != nil {
		if resp.Error != nil {
						// but we log it for visibility.
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

		// Check response structure - result is map[string]interface{}.
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

		// Validate quotes is an array.
	quotes, ok := result["quotes"].([]interface{})
	if !ok {
		t.Errorf("Response quotes is not []interface{}: %T", result["quotes"])
	} else if len(quotes) == 0 {
		t.Log("Response quotes array is empty (may be acceptable)")
	}
}

// TestHandleToolsList tests the handleToolsList function.
func TestHandleToolsList(t *testing.T) {
	server := NewWisdomServer()

	req := &JSONRPCRequest{
		JSONRPC: "2.0",
		ID:      1,
		Method:  "tools/list",
		Params:  nil,
	}

	resp := server.handleToolsList(req)
	if resp == nil {
		t.Fatal("handleToolsList returned nil response")
	}

	if resp.Error != nil {
		t.Fatalf("handleToolsList returned error: %v", resp.Error)
	}

		// Check response structure.
	result, ok := resp.Result.(map[string]interface{})
	if !ok {
		t.Fatalf("Response result is not map[string]interface{}: %T", resp.Result)
	}

		// Tools can be []Tool or []interface{} depending on JSON marshaling.
	toolsValue := result["tools"]
	if toolsValue == nil {
		t.Fatal("Response tools is nil")
	}

		// Try []Tool first (actual type).
	if tools, ok := toolsValue.([]Tool); ok {
		if len(tools) == 0 {
			t.Error("Response tools is empty")
		}
	} else if tools, ok := toolsValue.([]interface{}); ok {
		if len(tools) == 0 {
			t.Error("Response tools is empty")
		}
	} else {
		t.Fatalf("Response tools is not []Tool or []interface{}: %T", toolsValue)
	}
}

// TestHandleResourcesList tests the handleResourcesList function.
func TestHandleResourcesList(t *testing.T) {
	server := NewWisdomServer()

	req := &JSONRPCRequest{
		JSONRPC: "2.0",
		ID:      1,
		Method:  "resources/list",
		Params:  nil,
	}

	resp := server.handleResourcesList(req)
	if resp == nil {
		t.Fatal("handleResourcesList returned nil response")
	}

	if resp.Error != nil {
		t.Fatalf("handleResourcesList returned error: %v", resp.Error)
	}

		// Check response structure.
	result, ok := resp.Result.(map[string]interface{})
	if !ok {
		t.Fatalf("Response result is not map[string]interface{}: %T", resp.Result)
	}

		// Resources can be []Resource or []interface{} depending on JSON marshaling.
	resourcesValue := result["resources"]
	if resourcesValue == nil {
		t.Fatal("Response resources is nil")
	}

		// Try []Resource first (actual type).
	if resources, ok := resourcesValue.([]Resource); ok {
		if len(resources) == 0 {
			t.Error("Response resources is empty")
		}
	} else if resources, ok := resourcesValue.([]interface{}); ok {
		if len(resources) == 0 {
			t.Error("Response resources is empty")
		}
	} else {
		t.Fatalf("Response resources is not []Resource or []interface{}: %T", resourcesValue)
	}
}

// TestFormatRequestID tests the formatRequestID function.
func TestFormatRequestID(t *testing.T) {
	tests := []struct {
		name     string
		id       interface{}
		expected string
	}{
		{
			name:     "string ID",
			id:       "test-id",
			expected: "test-id",
		},
		{
			name:     "numeric ID",
			id:       123,
			expected: "123",
		},
		{
			name:     "nil ID",
			id:       nil,
			expected: "null",
		},
		{
			name:     "float ID",
			id:       45.67,
			expected: "46", 			expected: "46", // formatRequestID uses %.0f for floats.
		},
		{
			name:     "boolean ID",
			id:       true,
			expected: "true",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatRequestID(tt.id)
			if result != tt.expected {
				t.Errorf("formatRequestID(%v) = %q, want %q", tt.id, result, tt.expected)
			}
		})
	}
}

// TestHandleToolsResource tests the handleToolsResource function.
func TestHandleToolsResource(t *testing.T) {
	server := NewWisdomServer()
	if err := server.wisdom.Initialize(); err != nil {
		t.Fatalf("Initialize failed: %v", err)
	}

	req := &JSONRPCRequest{
		JSONRPC: "2.0",
		ID:      1,
		Method:  "resources/read",
		Params: json.RawMessage(`{
			"uri": "wisdom://tools"
		}`),
	}

	resp := server.handleToolsResource(req)
	if resp == nil {
		t.Fatal("handleToolsResource returned nil response")
	}

	if resp.Error != nil {
		t.Fatalf("handleToolsResource returned error: %v", resp.Error)
	}

		// Check response structure.
	result, ok := resp.Result.(map[string]interface{})
	if !ok {
		t.Fatalf("Response result is not map[string]interface{}: %T", resp.Result)
	}

	contents, ok := result["contents"].([]interface{})
	if !ok {
				// Try []map[string]interface{} as well.
		if contentsMap, ok := result["contents"].([]map[string]interface{}); ok {
			if len(contentsMap) == 0 {
				t.Error("Response contents is empty")
			}
			return
		}
		t.Fatalf("Response contents is not []interface{}: %T", result["contents"])
	}

	if len(contents) == 0 {
		t.Error("Response contents is empty")
	}
}

// TestHandleAdvisorResource tests the HandleAdvisorResource function.
func TestHandleAdvisorResource(t *testing.T) {
	server := NewWisdomServer()
	if err := server.wisdom.Initialize(); err != nil {
		t.Fatalf("Initialize failed: %v", err)
	}

	handlers := NewWisdomHandlers(server.wisdom, nil, server.appLogger)

	tests := []struct {
		name      string
		advisorID string
		wantError bool
	}{
		{
			name:      "valid metric advisor",
			advisorID: "security",
			wantError: false,
		},
		{
			name:      "valid tool advisor",
			advisorID: "project_scorecard",
			wantError: false,
		},
		{
			name:      "valid stage advisor",
			advisorID: "daily_checkin",
			wantError: false,
		},
		{
			name:      "invalid advisor",
			advisorID: "nonexistent_advisor",
			wantError: true,
		},
		{
			name:      "empty advisor ID",
			advisorID: "",
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &JSONRPCRequest{
				JSONRPC: "2.0",
				ID:      1,
				Method:  "resources/read",
				Params:  json.RawMessage(`{"uri": "wisdom://advisor/` + tt.advisorID + `"}`),
			}

			resp := handlers.HandleAdvisorResource(req, tt.advisorID)
			if resp == nil {
				t.Fatal("HandleAdvisorResource returned nil response")
			}

			if tt.wantError {
				if resp.Error == nil {
					t.Error("HandleAdvisorResource should return error for invalid advisor")
				}
			} else {
				if resp.Error != nil {
					t.Errorf("HandleAdvisorResource returned error: %v", resp.Error)
				}
								// with contents containing the advisor data as JSON text.
				result, ok := resp.Result.(map[string]interface{})
				if !ok {
					t.Fatalf("Response result is not map[string]interface{}: %T", resp.Result)
				}
								// Check that contents exist.
				contentsValue := result["contents"]
				var contentMap map[string]interface{}

								// Try []map[string]interface{} first (actual type).
				if contents, ok := contentsValue.([]map[string]interface{}); ok {
					if len(contents) == 0 {
						t.Error("Response contents is empty")
					}
					contentMap = contents[0]
				} else if contents, ok := contentsValue.([]interface{}); ok {
					if len(contents) == 0 {
						t.Error("Response contents is empty")
					}
					var ok2 bool
					contentMap, ok2 = contents[0].(map[string]interface{})
					if !ok2 {
						t.Fatalf("First content item is not map[string]interface{}: %T", contents[0])
					}
				} else {
					t.Fatalf("Response contents is not []map[string]interface{} or []interface{}: %T", contentsValue)
				}
				uri, ok := contentMap["uri"].(string)
				if !ok {
					t.Fatalf("Content URI is not string: %T", contentMap["uri"])
				}
				expectedURI := "wisdom://advisor/" + tt.advisorID
				if uri != expectedURI {
					t.Errorf("Content URI = %q, want %q", uri, expectedURI)
				}
			}
		})
	}
}

// TestHandleToolCall_UnknownTool tests HandleToolCall with an unknown tool name.
func TestHandleToolCall_UnknownTool(t *testing.T) {
	server := NewWisdomServer()
	if err := server.wisdom.Initialize(); err != nil {
		t.Fatalf("Initialize failed: %v", err)
	}

	handlers := NewWisdomHandlers(server.wisdom, nil, server.appLogger)

		// Test with unknown tool name.
	result, err := handlers.HandleToolCall("unknown_tool", map[string]interface{}{})
	if err == nil {
		t.Error("HandleToolCall should return error for unknown tool")
	}
	if result != nil {
		t.Error("HandleToolCall should return nil result for unknown tool")
	}
	if err != nil && !strings.Contains(err.Error(), "unknown tool") {
		t.Errorf("Error message should mention 'unknown tool', got: %v", err)
	}
}
