package mcp

import (
	"bytes"
	"context"
	"encoding/json"
	"strings"
	"testing"
	"time"
)

// TestIntegration_AllTools tests all 5 MCP tools via stdio transport
func TestIntegration_AllTools(t *testing.T) {
	server := NewWisdomServer()

	// Create test input with initialize and all tool calls
	input := strings.NewReader(`{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"test","version":"1.0.0"}}}
{"jsonrpc":"2.0","id":2,"method":"tools/call","params":{"name":"get_wisdom","arguments":{"score":75.0,"source":"stoic"}}}
{"jsonrpc":"2.0","id":3,"method":"tools/call","params":{"name":"consult_advisor","arguments":{"metric":"security","score":40.0}}}
{"jsonrpc":"2.0","id":4,"method":"tools/call","params":{"name":"get_daily_briefing","arguments":{"score":70.0}}}
{"jsonrpc":"2.0","id":5,"method":"tools/call","params":{"name":"get_consultation_log","arguments":{"days":7}}}
`)

	var output bytes.Buffer
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Run server in a goroutine since it blocks
	errChan := make(chan error, 1)
	go func() {
		errChan <- server.Run(ctx, input, &output)
	}()

	// Wait for completion or timeout
	select {
	case err := <-errChan:
		if err != nil && err != context.DeadlineExceeded {
			t.Fatalf("Server Run failed: %v", err)
		}
	case <-ctx.Done():
		// Server should complete when input EOF is reached
	}

	// Parse and verify responses
	outputStr := output.String()
	lines := strings.Split(strings.TrimSpace(outputStr), "\n")

	if len(lines) < 5 {
		t.Fatalf("Expected at least 5 responses, got %d. Output: %s", len(lines), outputStr)
	}

	// Verify each response
	responses := make(map[int]*JSONRPCResponse)
	for i, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		var resp JSONRPCResponse
		if err := json.Unmarshal([]byte(line), &resp); err != nil {
			t.Fatalf("Failed to parse response %d: %v\nLine: %s", i+1, err, line)
		}

		// Extract ID (could be int or string)
		var id int
		switch v := resp.ID.(type) {
		case float64:
			id = int(v)
		case int:
			id = v
		default:
			t.Fatalf("Unexpected ID type: %T", resp.ID)
		}

		responses[id] = &resp

		// Verify JSON-RPC version
		if resp.JSONRPC != "2.0" {
			t.Errorf("Response %d: JSONRPC version = %q, want %q", id, resp.JSONRPC, "2.0")
		}

		// Verify no errors for successful requests
		if id <= 4 && resp.Error != nil {
			t.Errorf("Response %d: Unexpected error: %v", id, resp.Error)
		}
	}

	// Verify initialize response (id: 1)
	if resp, ok := responses[1]; ok {
		if resp.Error != nil {
			t.Errorf("Initialize response error: %v", resp.Error)
		}
		if resp.Result == nil {
			t.Error("Initialize response missing result")
		}
	} else {
		t.Error("Missing initialize response")
	}

	// Verify get_wisdom response (id: 2)
	if resp, ok := responses[2]; ok {
		if resp.Error != nil {
			t.Errorf("get_wisdom response error: %v", resp.Error)
		}
		if resp.Result == nil {
			t.Error("get_wisdom response missing result")
		}
	} else {
		t.Error("Missing get_wisdom response")
	}

	// Verify consult_advisor response (id: 3)
	if resp, ok := responses[3]; ok {
		if resp.Error != nil {
			t.Errorf("consult_advisor response error: %v", resp.Error)
		}
		if resp.Result == nil {
			t.Error("consult_advisor response missing result")
		}
	} else {
		t.Error("Missing consult_advisor response")
	}

	// Verify get_daily_briefing response (id: 4)
	if resp, ok := responses[4]; ok {
		if resp.Error != nil {
			t.Errorf("get_daily_briefing response error: %v", resp.Error)
		}
		if resp.Result == nil {
			t.Error("get_daily_briefing response missing result")
		}
	} else {
		t.Error("Missing get_daily_briefing response")
	}
}

// TestIntegration_AllResources tests all 4 MCP resources via stdio transport
func TestIntegration_AllResources(t *testing.T) {
	server := NewWisdomServer()

	// Create test input with initialize and all resource reads
	input := strings.NewReader(`{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"test","version":"1.0.0"}}}
{"jsonrpc":"2.0","id":2,"method":"resources/read","params":{"uri":"wisdom://sources"}}
{"jsonrpc":"2.0","id":3,"method":"resources/read","params":{"uri":"wisdom://advisors"}}
{"jsonrpc":"2.0","id":4,"method":"resources/read","params":{"uri":"wisdom://advisor/security"}}
{"jsonrpc":"2.0","id":5,"method":"resources/read","params":{"uri":"wisdom://consultations/7"}}
`)

	var output bytes.Buffer
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Run server in a goroutine since it blocks
	errChan := make(chan error, 1)
	go func() {
		errChan <- server.Run(ctx, input, &output)
	}()

	// Wait for completion or timeout
	select {
	case err := <-errChan:
		if err != nil && err != context.DeadlineExceeded {
			t.Fatalf("Server Run failed: %v", err)
		}
	case <-ctx.Done():
		// Server should complete when input EOF is reached
	}

	// Parse and verify responses
	outputStr := output.String()
	lines := strings.Split(strings.TrimSpace(outputStr), "\n")

	if len(lines) < 5 {
		t.Fatalf("Expected at least 5 responses, got %d. Output: %s", len(lines), outputStr)
	}

	// Verify each response
	responses := make(map[int]*JSONRPCResponse)
	for i, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		var resp JSONRPCResponse
		if err := json.Unmarshal([]byte(line), &resp); err != nil {
			t.Fatalf("Failed to parse response %d: %v\nLine: %s", i+1, err, line)
		}

		// Extract ID
		var id int
		switch v := resp.ID.(type) {
		case float64:
			id = int(v)
		case int:
			id = v
		default:
			t.Fatalf("Unexpected ID type: %T", resp.ID)
		}

		responses[id] = &resp

		// Verify JSON-RPC version
		if resp.JSONRPC != "2.0" {
			t.Errorf("Response %d: JSONRPC version = %q, want %q", id, resp.JSONRPC, "2.0")
		}

		// Verify no errors for successful requests
		if resp.Error != nil {
			t.Errorf("Response %d: Unexpected error: %v", id, resp.Error)
		}
	}

	// Verify all resource responses have results
	for id := 2; id <= 5; id++ {
		if resp, ok := responses[id]; ok {
			if resp.Result == nil {
				t.Errorf("Resource response %d missing result", id)
			}
		} else {
			t.Errorf("Missing resource response %d", id)
		}
	}
}

// TestIntegration_ErrorHandling tests error scenarios via stdio transport
func TestIntegration_ErrorHandling(t *testing.T) {
	server := NewWisdomServer()

	// Create test input with invalid requests
	input := strings.NewReader(`{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"test","version":"1.0.0"}}}
{"jsonrpc":"2.0","id":2,"method":"invalid/method","params":{}}
{"jsonrpc":"2.0","id":3,"method":"tools/call","params":{"name":"get_wisdom","arguments":{}}}
{"jsonrpc":"2.0","id":4,"method":"resources/read","params":{"uri":"wisdom://invalid"}}
`)

	var output bytes.Buffer
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Run server in a goroutine since it blocks
	errChan := make(chan error, 1)
	go func() {
		errChan <- server.Run(ctx, input, &output)
	}()

	// Wait for completion or timeout
	select {
	case err := <-errChan:
		if err != nil && err != context.DeadlineExceeded {
			t.Fatalf("Server Run failed: %v", err)
		}
	case <-ctx.Done():
		// Server should complete when input EOF is reached
	}

	// Parse and verify error responses
	outputStr := output.String()
	lines := strings.Split(strings.TrimSpace(outputStr), "\n")

	// Verify error responses
	for i, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		var resp JSONRPCResponse
		if err := json.Unmarshal([]byte(line), &resp); err != nil {
			t.Fatalf("Failed to parse response %d: %v\nLine: %s", i+1, err, line)
		}

		// Extract ID
		var id int
		switch v := resp.ID.(type) {
		case float64:
			id = int(v)
		case int:
			id = v
		default:
			continue // Skip non-numeric IDs
		}

		// ID 1 (initialize) should succeed
		if id == 1 {
			if resp.Error != nil {
				t.Errorf("Initialize should not error: %v", resp.Error)
			}
			continue
		}

		// IDs 2-4 should have errors
		if id >= 2 && id <= 4 {
			if resp.Error == nil {
				t.Errorf("Response %d should have error but doesn't", id)
			} else {
				// Verify error code is valid JSON-RPC error code
				if resp.Error.Code > 0 {
					t.Errorf("Response %d: Error code should be negative, got %d", id, resp.Error.Code)
				}
			}
		}
	}
}

// TestIntegration_ProtocolCompliance tests JSON-RPC 2.0 protocol compliance
func TestIntegration_ProtocolCompliance(t *testing.T) {
	server := NewWisdomServer()

	// Create test input with various protocol scenarios
	input := strings.NewReader(`{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"test","version":"1.0.0"}}}
{"jsonrpc":"2.0","id":2,"method":"tools/call","params":{"name":"get_wisdom","arguments":{"score":75.0,"source":"stoic"}}}
{"jsonrpc":"2.0","id":null,"method":"tools/call","params":{"name":"get_wisdom","arguments":{"score":75.0,"source":"stoic"}}}
`)

	var output bytes.Buffer
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Run server in a goroutine since it blocks
	errChan := make(chan error, 1)
	go func() {
		errChan <- server.Run(ctx, input, &output)
	}()

	// Wait for completion or timeout
	select {
	case err := <-errChan:
		if err != nil && err != context.DeadlineExceeded {
			t.Fatalf("Server Run failed: %v", err)
		}
	case <-ctx.Done():
		// Server should complete when input EOF is reached
	}

	// Parse responses
	outputStr := output.String()
	lines := strings.Split(strings.TrimSpace(outputStr), "\n")

	responseCount := 0
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		var resp JSONRPCResponse
		if err := json.Unmarshal([]byte(line), &resp); err != nil {
			t.Fatalf("Failed to parse response: %v\nLine: %s", err, line)
		}

		responseCount++

		// Verify JSON-RPC version
		if resp.JSONRPC != "2.0" {
			t.Errorf("Response JSONRPC version = %q, want %q", resp.JSONRPC, "2.0")
		}

		// Verify ID is present (notifications should not get responses)
		if resp.ID == nil {
			t.Error("Response missing ID (notifications should not get responses)")
		}
	}

	// Should have 2 responses (initialize + get_wisdom, notification should not get response)
	if responseCount != 2 {
		t.Errorf("Expected 2 responses (notification should not get response), got %d", responseCount)
	}
}

// TestIntegration_SequentialRequests tests multiple sequential requests
func TestIntegration_SequentialRequests(t *testing.T) {
	server := NewWisdomServer()

	// Create test input with multiple sequential tool calls
	input := strings.NewReader(`{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"test","version":"1.0.0"}}}
{"jsonrpc":"2.0","id":2,"method":"tools/call","params":{"name":"get_wisdom","arguments":{"score":50.0,"source":"stoic"}}}
{"jsonrpc":"2.0","id":3,"method":"tools/call","params":{"name":"get_wisdom","arguments":{"score":75.0,"source":"tao"}}}
{"jsonrpc":"2.0","id":4,"method":"tools/call","params":{"name":"get_wisdom","arguments":{"score":90.0,"source":"pistis_sophia"}}}
`)

	var output bytes.Buffer
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Run server in a goroutine since it blocks
	errChan := make(chan error, 1)
	go func() {
		errChan <- server.Run(ctx, input, &output)
	}()

	// Wait for completion or timeout
	select {
	case err := <-errChan:
		if err != nil && err != context.DeadlineExceeded {
			t.Fatalf("Server Run failed: %v", err)
		}
	case <-ctx.Done():
		// Server should complete when input EOF is reached
	}

	// Parse and verify all responses
	outputStr := output.String()
	lines := strings.Split(strings.TrimSpace(outputStr), "\n")

	responses := make(map[int]*JSONRPCResponse)
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		var resp JSONRPCResponse
		if err := json.Unmarshal([]byte(line), &resp); err != nil {
			t.Fatalf("Failed to parse response: %v\nLine: %s", err, line)
		}

		// Extract ID
		var id int
		switch v := resp.ID.(type) {
		case float64:
			id = int(v)
		case int:
			id = v
		default:
			continue
		}

		responses[id] = &resp

		// Verify no errors
		if resp.Error != nil {
			t.Errorf("Response %d: Unexpected error: %v", id, resp.Error)
		}

		// Verify result exists
		if resp.Result == nil {
			t.Errorf("Response %d: Missing result", id)
		}
	}

	// Verify all 4 responses received
	if len(responses) != 4 {
		t.Errorf("Expected 4 responses, got %d", len(responses))
	}
}

