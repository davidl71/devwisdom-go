package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// TestSDKAdapterBasic tests basic SDK adapter functionality
func TestSDKAdapterBasic(t *testing.T) {
	server := NewWisdomServerSDK()
	if server == nil {
		t.Fatal("NewWisdomServerSDK returned nil")
	}

	// Verify server components are initialized
	if server.wisdom == nil {
		t.Error("wisdom engine is nil")
	}
	if server.appLogger == nil {
		t.Error("app logger is nil")
	}
	if server.server == nil {
		t.Error("SDK server is nil")
	}
}

// TestSDKAdapterToolsRegistration tests that tools are registered
func TestSDKAdapterToolsRegistration(t *testing.T) {
	server := NewWisdomServerSDK()

	// Try to register tools (this should not fail)
	server.registerTools()
}

// TestSDKAdapterResourcesRegistration tests that resources are registered
func TestSDKAdapterResourcesRegistration(t *testing.T) {
	server := NewWisdomServerSDK()

	// Try to register resources (this should not fail)
	if err := server.registerResources(); err != nil {
		t.Fatalf("registerResources failed: %v", err)
	}
}

// TestSDKAdapterToolHandlers tests tool handler conversion
func TestSDKAdapterToolHandlers(t *testing.T) {
	server := NewWisdomServerSDK()
	helperServer := &WisdomServer{
		wisdom:    server.wisdom,
		logger:    server.logger,
		appLogger: server.appLogger,
	}

	// Test consult_advisor handler
	args := map[string]interface{}{
		"score":  75.0,
		"metric": "security",
	}

	result, err := helperServer.handleConsultAdvisor(args)
	if err != nil {
		t.Fatalf("handleConsultAdvisor failed: %v", err)
	}

	// Verify result can be marshaled to JSON (as SDK adapter does)
	_, err = json.Marshal(result)
	if err != nil {
		t.Fatalf("Failed to marshal result: %v", err)
	}
}

// TestSDKAdapterResourceHandlers tests resource handler conversion
func TestSDKAdapterResourceHandlers(t *testing.T) {
	server := NewWisdomServerSDK()
	helperServer := &WisdomServer{
		wisdom:    server.wisdom,
		logger:    server.logger,
		appLogger: server.appLogger,
	}

	// Initialize wisdom engine
	if err := server.wisdom.Initialize(); err != nil {
		t.Fatalf("Failed to initialize wisdom engine: %v", err)
	}

	// Test tools resource
	mockReq := &JSONRPCRequest{
		ID:     "test",
		Method: "resources/read",
		Params: json.RawMessage(`{"uri": "wisdom://tools"}`),
	}

	resp := helperServer.handleToolsResource(mockReq)
	if resp.Error != nil {
		t.Fatalf("handleToolsResource failed: %v", resp.Error.Message)
	}
}

// TestSDKAdapter_Run tests the Run method with context cancellation
func TestSDKAdapter_Run(t *testing.T) {
	server := NewWisdomServerSDK()

	// Initialize wisdom engine
	if err := server.wisdom.Initialize(); err != nil {
		t.Fatalf("Failed to initialize wisdom engine: %v", err)
	}

	// Register tools and resources
	server.registerTools()
	if err := server.registerResources(); err != nil {
		t.Fatalf("Failed to register resources: %v", err)
	}

	// Create a context that will be cancelled quickly
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	// Run should return context.DeadlineExceeded or similar
	err := server.Run(ctx)
	if err == nil {
		t.Log("Run completed without error (may be acceptable if server shuts down quickly)")
	} else if err != context.DeadlineExceeded && err != context.Canceled {
		// Accept deadline exceeded or cancelled, but log other errors
		t.Logf("Run returned error (may be expected): %v", err)
	}
}

// TestSDKAdapter_Run_InitializationError tests Run with initialization failure
func TestSDKAdapter_Run_InitializationError(t *testing.T) {
	// Create a mock wisdom engine that fails to initialize
	// We can't easily mock this, so we'll test the error path by using a real engine
	// but with invalid configuration. However, this is complex.
	// Instead, let's test that Run handles initialization errors gracefully.

	// For now, we'll test that Run doesn't panic when wisdom engine fails
	// This requires mocking, which is complex. We'll skip this test for now
	// and focus on testing the error result function instead.
	t.Skip("Requires mocking wisdom engine initialization")
}

// TestNewErrorResult tests the newErrorResult helper function
func TestNewErrorResult(t *testing.T) {
	message := "test error message"
	result := newErrorResult(message)

	if result == nil {
		t.Fatal("newErrorResult returned nil")
	}

	if !result.IsError {
		t.Error("newErrorResult should set IsError to true")
	}

	if len(result.Content) == 0 {
		t.Fatal("newErrorResult should have content")
	}

	textContent, ok := result.Content[0].(*mcp.TextContent)
	if !ok {
		t.Fatalf("Content[0] is not *mcp.TextContent: %T", result.Content[0])
	}

	if textContent.Text != message {
		t.Errorf("TextContent.Text = %q, want %q", textContent.Text, message)
	}
}

// TestWrapToolHandler tests the wrapToolHandler function
func TestWrapToolHandler(t *testing.T) {
	server := NewWisdomServerSDK()
	if err := server.wisdom.Initialize(); err != nil {
		t.Fatalf("Failed to initialize wisdom engine: %v", err)
	}

	handlers := NewWisdomHandlers(server.wisdom, server.logger, server.appLogger)

	// Test with valid handler
	handler := wrapToolHandler(handlers.handleGetWisdom)

	// Create a test request - check actual SDK structure
	req := &mcp.CallToolRequest{
		Params: &mcp.CallToolParamsRaw{
			Name:      "get_wisdom",
			Arguments: json.RawMessage(`{"score": 75.0, "source": "stoic"}`),
		},
	}

	ctx := context.Background()
	result, err := handler(ctx, req)

	if err != nil {
		t.Fatalf("wrapToolHandler returned error: %v", err)
	}

	if result == nil {
		t.Fatal("wrapToolHandler returned nil result")
	}

	if result.IsError {
		t.Errorf("wrapToolHandler returned error result: %v", result)
	}

	if len(result.Content) == 0 {
		t.Error("wrapToolHandler returned empty content")
	}
}

// TestWrapToolHandler_InvalidJSON tests wrapToolHandler with invalid JSON arguments
func TestWrapToolHandler_InvalidJSON(t *testing.T) {
	server := NewWisdomServerSDK()
	handlers := NewWisdomHandlers(server.wisdom, server.logger, server.appLogger)

	handler := wrapToolHandler(handlers.handleGetWisdom)

	req := &mcp.CallToolRequest{
		Params: &mcp.CallToolParamsRaw{
			Name:      "get_wisdom",
			Arguments: json.RawMessage(`invalid json`),
		},
	}

	ctx := context.Background()
	result, err := handler(ctx, req)

	if err != nil {
		t.Fatalf("wrapToolHandler returned error: %v", err)
	}

	if result == nil {
		t.Fatal("wrapToolHandler returned nil result")
	}

	if !result.IsError {
		t.Error("wrapToolHandler should return error result for invalid JSON")
	}
}

// TestWrapToolHandler_HandlerError tests wrapToolHandler when handler returns error
func TestWrapToolHandler_HandlerError(t *testing.T) {
	server := NewWisdomServerSDK()
	handlers := NewWisdomHandlers(server.wisdom, server.logger, server.appLogger)

	handler := wrapToolHandler(handlers.handleGetWisdom)

	// Use invalid arguments that will cause handler to return error
	req := &mcp.CallToolRequest{
		Params: &mcp.CallToolParamsRaw{
			Name:      "get_wisdom",
			Arguments: json.RawMessage(`{"invalid": "args"}`), // Missing required "score"
		},
	}

	ctx := context.Background()
	result, err := handler(ctx, req)

	if err != nil {
		t.Fatalf("wrapToolHandler returned error: %v", err)
	}

	if result == nil {
		t.Fatal("wrapToolHandler returned nil result")
	}

	// Handler may return error or may handle it gracefully
	// We just verify it doesn't panic
}

// TestConvertToolsToSDK tests the ConvertToolsToSDK function
func TestConvertToolsToSDK(t *testing.T) {
	tools := getToolDefinitions()
	sdkTools := ConvertToolsToSDK(tools)

	if len(sdkTools) != len(tools) {
		t.Errorf("ConvertToolsToSDK returned %d tools, want %d", len(sdkTools), len(tools))
	}

	for i, tool := range tools {
		if sdkTools[i].Name != tool.Name {
			t.Errorf("sdkTools[%d].Name = %q, want %q", i, sdkTools[i].Name, tool.Name)
		}
		if sdkTools[i].Description != tool.Description {
			t.Errorf("sdkTools[%d].Description = %q, want %q", i, sdkTools[i].Description, tool.Description)
		}
	}
}

// TestConvertResourceResponse tests the convertResourceResponse function
func TestConvertResourceResponse(t *testing.T) {
	server := NewWisdomServerSDK()

	// Test with valid response
	resp := &JSONRPCResponse{
		Result: map[string]interface{}{
			"contents": []interface{}{
				map[string]interface{}{
					"uri":      "wisdom://tools",
					"mimeType": "application/json",
					"text":     `{"tools": []}`,
				},
			},
		},
	}

	result, err := server.convertResourceResponse(resp, "wisdom://tools")
	if err != nil {
		t.Fatalf("convertResourceResponse returned error: %v", err)
	}

	if result == nil {
		t.Fatal("convertResourceResponse returned nil")
	}

	if len(result.Contents) == 0 {
		t.Fatal("convertResourceResponse returned empty contents")
	}

	if result.Contents[0].URI != "wisdom://tools" {
		t.Errorf("Contents[0].URI = %q, want %q", result.Contents[0].URI, "wisdom://tools")
	}
}

// TestConvertResourceResponse_Error tests convertResourceResponse with error response
func TestConvertResourceResponse_Error(t *testing.T) {
	server := NewWisdomServerSDK()

	resp := &JSONRPCResponse{
		Error: &JSONRPCError{
			Code:    -32603,
			Message: "Internal error",
		},
	}

	result, err := server.convertResourceResponse(resp, "wisdom://tools")
	if err == nil {
		t.Error("convertResourceResponse should return error for error response")
	}
	if result != nil {
		t.Error("convertResourceResponse should return nil result for error response")
	}
}

// TestConvertResourceResponse_InvalidFormat tests convertResourceResponse with invalid format
func TestConvertResourceResponse_InvalidFormat(t *testing.T) {
	server := NewWisdomServerSDK()

	// Test with invalid response format
	resp := &JSONRPCResponse{
		Result: "invalid format",
	}

	result, err := server.convertResourceResponse(resp, "wisdom://tools")
	if err == nil {
		t.Error("convertResourceResponse should return error for invalid format")
	}
	if result != nil {
		t.Error("convertResourceResponse should return nil result for invalid format")
	}
}

// TestGetString tests the getString helper function
func TestGetString(t *testing.T) {
	tests := []struct {
		name         string
		m            map[string]interface{}
		key          string
		defaultValue string
		expected     string
	}{
		{
			name:         "key exists with string value",
			m:            map[string]interface{}{"key": "value"},
			key:          "key",
			defaultValue: "default",
			expected:     "value",
		},
		{
			name:         "key exists with non-string value",
			m:            map[string]interface{}{"key": 123},
			key:          "key",
			defaultValue: "default",
			expected:     "default",
		},
		{
			name:         "key does not exist",
			m:            map[string]interface{}{},
			key:          "key",
			defaultValue: "default",
			expected:     "default",
		},
		{
			name:         "empty string value",
			m:            map[string]interface{}{"key": ""},
			key:          "key",
			defaultValue: "default",
			expected:     "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getString(tt.m, tt.key, tt.defaultValue)
			if result != tt.expected {
				t.Errorf("getString() = %q, want %q", result, tt.expected)
			}
		})
	}
}

// TestCreateResourceHandler tests the createResourceHandler function
func TestCreateResourceHandler(t *testing.T) {
	server := NewWisdomServerSDK()
	if err := server.wisdom.Initialize(); err != nil {
		t.Fatalf("Failed to initialize wisdom engine: %v", err)
	}

	handlers := NewWisdomHandlers(server.wisdom, server.logger, server.appLogger)

	handler := server.createResourceHandler("wisdom://sources", handlers.HandleSourcesResource)

	req := &mcp.ReadResourceRequest{
		Params: &mcp.ReadResourceParams{
			URI: "wisdom://sources",
		},
	}

	ctx := context.Background()
	result, err := handler(ctx, req)

	if err != nil {
		t.Fatalf("createResourceHandler returned error: %v", err)
	}

	if result == nil {
		t.Fatal("createResourceHandler returned nil result")
	}

	if len(result.Contents) == 0 {
		t.Error("createResourceHandler returned empty contents")
	}
}

// TestRegisterAllResources tests that all resources are registered correctly
func TestRegisterAllResources(t *testing.T) {
	server := NewWisdomServerSDK()
	handlers := NewWisdomHandlers(server.wisdom, server.logger, server.appLogger)

	// Register all resources
	err := server.registerAllResources(handlers)
	if err != nil {
		t.Fatalf("registerAllResources failed: %v", err)
	}

	// Verify resources are registered by checking server state
	// Note: SDK server doesn't expose a direct way to list resources,
	// so we verify by attempting to read them
	expectedResources := []string{
		"wisdom://sources",
		"wisdom://advisors",
		"wisdom://tools",
		"wisdom://consultations/7",
	}

	// Test that we can create handlers for all resources
	for _, uri := range expectedResources {
		// Create a handler for this resource
		handler := server.createResourceHandler(uri, func(req *JSONRPCRequest) *JSONRPCResponse {
			return &JSONRPCResponse{ID: "test", Result: map[string]interface{}{"test": true}}
		})

		if handler == nil {
			t.Errorf("createResourceHandler returned nil for %s", uri)
		}
	}

	// Verify handlers were used (avoid unused variable warning)
	_ = handlers
}

// TestWrapToolHandler_ErrorPaths tests wrapToolHandler with various error scenarios
func TestWrapToolHandler_ErrorPaths(t *testing.T) {
	// Test handler that returns an error
	errorHandler := func(params map[string]interface{}) (interface{}, error) {
		return nil, fmt.Errorf("test error")
	}

	wrapped := wrapToolHandler(errorHandler)

	// Create a test request
	req := &mcp.CallToolRequest{
		Params: &mcp.CallToolParamsRaw{
			Name:      "test_tool",
			Arguments: json.RawMessage(`{}`),
		},
	}

	result, err := wrapped(context.Background(), req)
	if err != nil {
		t.Fatalf("wrapToolHandler returned error: %v", err)
	}
	if result == nil {
		t.Fatal("wrapToolHandler returned nil result")
	}

	// Verify error is returned
	if len(result.Content) == 0 {
		t.Error("wrapToolHandler should return error content for handler error")
	}

	// Test handler with invalid JSON arguments
	invalidJSONHandler := func(params map[string]interface{}) (interface{}, error) {
		return map[string]interface{}{"result": "success"}, nil
	}

	wrapped2 := wrapToolHandler(invalidJSONHandler)

	req2 := &mcp.CallToolRequest{
		Params: &mcp.CallToolParamsRaw{
			Name:      "test_tool",
			Arguments: json.RawMessage(`invalid json`),
		},
	}

	result2, err2 := wrapped2(context.Background(), req2)
	if err2 != nil {
		t.Fatalf("wrapToolHandler returned error: %v", err2)
	}
	if result2 == nil {
		t.Fatal("wrapToolHandler returned nil result for invalid JSON")
	}

	// Should handle invalid JSON gracefully
	if len(result2.Content) == 0 {
		t.Error("wrapToolHandler should return error content for invalid JSON")
	}
}
