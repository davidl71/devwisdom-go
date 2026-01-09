package mcp

import (
	"encoding/json"
	"testing"
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
	if err := server.registerTools(); err != nil {
		t.Fatalf("registerTools failed: %v", err)
	}
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
		"score": 75.0,
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

